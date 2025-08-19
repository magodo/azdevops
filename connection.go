package azdevops

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/google/uuid"
)

type ConnectionOption struct {
	ModuleName    string
	ModuleVersion string
	policy.ClientOptions
	runtime.PipelineOptions
}

type Connection struct {
	ep            string
	pl            runtime.Pipeline
	tr            tracing.Tracer
	resourceAreas map[uuid.UUID]ResourceArea

	// clientCache caches the client for different location url (key)
	clientCache     map[string]*Client
	clientCacheLock sync.Mutex
}

func NewConnection(ctx context.Context, endpoint string, cred *Credential, options *ConnectionOption) (*Connection, error) {
	if options == nil {
		options = &ConnectionOption{}
	}

	pl, err := NewPipeline(options.ModuleName, options.ModuleVersion, cred, &options.PipelineOptions, &options.ClientOptions)
	if err != nil {
		return nil, err
	}
	tr := options.TracingProvider.NewTracer(options.ModuleName, options.ModuleVersion)

	conn := &Connection{
		ep:            endpoint,
		pl:            pl,
		tr:            tr,
		resourceAreas: map[uuid.UUID]ResourceArea{},
		clientCache:   map[string]*Client{},
	}

	{
		// Construct a temporary annonymous client to get the resource areas of this ADO instance
		pl, err := NewPipeline(options.ModuleName, options.ModuleVersion, nil, &options.PipelineOptions, &options.ClientOptions)
		if err != nil {
			return nil, err
		}
		c := &Client{ep: endpoint, pl: pl, tr: tr}
		resp, err := c.Do(ctx, http.MethodGet, "_apis/ResourceAreas", nil, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to query resource areas information: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, runtime.NewResponseError(resp)
		}

		var areas []ResourceArea
		if err := UnmarshalBody(resp, &areas, UnmarshalCollection); err != nil {
			return nil, err
		}
		for _, area := range areas {
			if area.Id != nil {
				conn.resourceAreas[*area.Id] = area
			}
		}
	}

	return conn, nil
}

func (conn *Connection) NewDefaultClient(ctx context.Context) (*Client, error) {
	return conn.newAreaClient(ctx, conn.ep)
}

func (conn *Connection) NewAreaClient(ctx context.Context, areaId uuid.UUID) (*Client, error) {
	ep := conn.ep
	area, ok := conn.resourceAreas[areaId]
	if !ok {
		return nil, fmt.Errorf("unknown area %s", areaId)
	}
	if v := area.LocationUrl; v != nil {
		ep = *v
	}

	return conn.newAreaClient(ctx, ep)
}

func (conn *Connection) newAreaClient(ctx context.Context, ep string) (*Client, error) {
	conn.clientCacheLock.Lock()
	defer conn.clientCacheLock.Unlock()
	c, ok := conn.clientCache[ep]
	if ok {
		return c, nil
	}

	c = &Client{
		ep:   ep,
		pl:   conn.pl,
		tr:   conn.tr,
		apis: map[uuid.UUID]Api{},
	}

	{
		// Populate the apis for this client
		resp, err := c.Do(ctx, http.MethodOptions, "_apis", nil, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to query apis information: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, runtime.NewResponseError(resp)
		}

		var apis []Api
		if err := UnmarshalBody(resp, &apis, UnmarshalCollection); err != nil {
			return nil, err
		}
		for _, api := range apis {
			if api.Id != nil {
				c.apis[*api.Id] = api
			}
		}
	}

	conn.clientCache[ep] = c

	return c, nil
}
