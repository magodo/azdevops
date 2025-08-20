package azdevops_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/magodo/azdevops"
	"github.com/stretchr/testify/require"
)

func TestClient_Process(t *testing.T) {
	ep, ok := os.LookupEnv("AZDEVOPS_ADO_ENDPOINT")
	if !ok {
		t.Skip(`"AZDEVOPS_ADO_ENDPOINT" not specified`)
	}
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)
	conn, err := azdevops.NewConnection(t.Context(), ep, &azdevops.Credential{AADTokenCredential: cred}, nil)
	require.NoError(t, err)

	c, err := conn.NewDefaultClient(t.Context())
	require.NoError(t, err)

	// Get the default template
	{
		resp, err := c.Do(t.Context(), http.MethodGet, "_apis/process/processes", nil, nil, map[string]string{"api-version": "7.1-preview.1"})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var templates []map[string]any
		require.NoError(t, azdevops.UnmarshalBody(resp, &templates, azdevops.UnmarshalCollection))

		var templateId string
		for _, tpl := range templates {
			if tpl["isDefault"].(bool) {
				templateId = tpl["id"].(string)
				break
			}
		}
		require.NotEqual(t, "", templateId)
	}
}
