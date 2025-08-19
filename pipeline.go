package azdevops

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// Credential specifies the token credential to use.
// One of the fields must be specified.
type Credential struct {
	AADTokenCredential azcore.TokenCredential
	PAT                *string
}

// NewPipeline creates a pipeline from connection options. Policies from ClientOptions are
// placed after policies from PipelineOptions.
func NewPipeline(module, version string, cred *Credential, plOpts *runtime.PipelineOptions, cliOpts *policy.ClientOptions) (runtime.Pipeline, error) {
	if plOpts == nil {
		plOpts = &runtime.PipelineOptions{}
	}
	if cliOpts == nil {
		cliOpts = &policy.ClientOptions{}
	}
	if cred != nil {
		var authPolicy policy.Policy
		switch {
		case cred.AADTokenCredential != nil:
			authPolicy = runtime.NewBearerTokenPolicy(cred.AADTokenCredential, []string{"499b84ac-1321-427f-aa17-267ca6975798"}, &policy.BearerTokenOptions{
				InsecureAllowCredentialWithHTTP: cliOpts.InsecureAllowCredentialWithHTTP,
			})
		case cred.PAT != nil:
			authPolicy = NewPATPolicy(*cred.PAT, nil)
		default:
			return runtime.Pipeline{}, fmt.Errorf("Either the AAD Token Credential or Personal Access Token shall be specified")
		}
		// we don't want to modify the underlying array in plOpts.PerRetry
		perRetry := make([]policy.Policy, len(plOpts.PerRetry), len(plOpts.PerRetry)+1)
		copy(perRetry, plOpts.PerRetry)
		perRetry = append(perRetry, authPolicy)
		plOpts.PerRetry = perRetry
	}
	if plOpts.APIVersion.Name == "" {
		plOpts.APIVersion.Name = "api-version"
	}
	return runtime.NewPipeline(module, version, *plOpts, cliOpts), nil
}
