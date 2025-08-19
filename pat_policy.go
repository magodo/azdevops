package azdevops

import (
	"encoding/base64"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type PATOptions struct {
}

type PATPolicy struct {
	pat string
}

func NewPATPolicy(pat string, opts *PATOptions) *PATPolicy {
	if opts == nil {
		opts = &PATOptions{}
	}
	return &PATPolicy{
		pat: pat,
	}
}

func (p *PATPolicy) Do(req *policy.Request) (*http.Response, error) {
	if p.pat == "" {
		return req.Next()
	}

	req.Raw().Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("_:"+p.pat)))
	res, err := req.Next()
	if err != nil {
		return nil, err
	}
	return res, err
}
