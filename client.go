package azdevops

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/google/uuid"
)

type Client struct {
	ep   string
	pl   runtime.Pipeline
	tr   tracing.Tracer
	apis map[uuid.UUID]Api
}

func (c *Client) Do(ctx context.Context, method string, path string, body []byte, additionalHeaders map[string]string, additionalQueries map[string]string) (*http.Response, error) {
	endpoint := strings.TrimRight(c.ep, "/") + "/" + strings.TrimLeft(path, "/")
	req, err := runtime.NewRequest(ctx, method, endpoint)
	if err != nil {
		return nil, err
	}
	if body != nil {
		// TODO: Support other media types
		if err := runtime.MarshalAsJSON(req, body); err != nil {
			return nil, err
		}
	}
	for k, v := range additionalHeaders {
		req.Raw().Header.Add(k, v)
	}

	if len(additionalQueries) != 0 {
		q := req.Raw().URL.Query()
		for k, v := range additionalQueries {
			q.Add(k, v)
		}
		req.Raw().URL.RawQuery = q.Encode()
	}

	return c.pl.Do(req)
}
