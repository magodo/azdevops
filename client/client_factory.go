package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/google/uuid"
	"github.com/magodo/azdevops/internal/pipeline"
	"github.com/magodo/azdevops/model"
)

type FactoryOption struct {
	ModuleName    string
	ModuleVersion string
	policy.ClientOptions
}

type ClientFactory struct {
	ep                string
	pl                runtime.Pipeline
	tr                tracing.Tracer
	coreApiCache      map[uuid.UUID]model.Api
	resourceAreaCache map[uuid.UUID]model.ResourceArea
}

func NewClientFactory(ctx context.Context, endpoint string, cred azcore.TokenCredential, options *FactoryOption) (*ClientFactory, error) {
	if options == nil {
		options = &FactoryOption{}
	}

	pl, err := pipeline.NewPipeline(options.ModuleName, options.ModuleVersion, cred, runtime.PipelineOptions{}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}
	tr := options.TracingProvider.NewTracer(options.ModuleName, options.ModuleVersion)
	f := &ClientFactory{ep: endpoint, pl: pl, tr: tr}

	c := Client{ep: f.ep, pl: f.pl, tr: f.tr}

	normalize := func(p *string) {
		if p != nil {
			*p = strings.ToLower(*p)
		}
	}

	{
		// Each location URL has the `_apis` endpoint available. The information returned from subdomains might not be contained in the one returned from the main domain (e.g. `dev.azure.com/<org>`).
		// Here we are retrieving the main domain's apis information, only to get the path for the (Location, ResourceAreas).
		resp, err := c.Do(ctx, http.MethodOptions, "_apis", nil, nil, nil)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			errmsg, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("API error for reading _apis (%d): %v", resp.StatusCode, errmsg)
		}

		var apis []model.Api
		if err := UnmarshalBody(resp, &apis, UnmarshalCollection); err != nil {
			return nil, err
		}

		m := map[uuid.UUID]model.Api{}
		for _, api := range apis {
			normalize(api.Area)
			normalize(api.ResourceName)
			m[*api.Id] = api
		}
		f.coreApiCache = m
	}

	{
		// TODO: Instead of hardcode the path here, we shall use the apis info to lookup the path for (area, resource) of (Location, ResourceAreas).
		resp, err := c.Do(ctx, http.MethodGet, "_apis/ResourceAreas", nil, nil, nil)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			errmsg, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("API error for reading _apis/ResourceAreas (%d): %v", resp.StatusCode, errmsg)
		}

		var areas []model.ResourceArea
		if err := UnmarshalBody(resp, &areas, UnmarshalCollection); err != nil {
			return nil, err
		}

		m := map[uuid.UUID]model.ResourceArea{}
		for _, area := range areas {
			normalize(area.Name)
			m[*area.Id] = area
		}
		f.resourceAreaCache = m
	}

	return f, nil
}
