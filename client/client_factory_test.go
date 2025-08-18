package client_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/magodo/azdevops/client"
)

func TestClientFactory(t *testing.T) {
	ep, ok := os.LookupEnv("AZDEVOPS_ADO_ENDPOINT")
	if !ok {
		t.Skip(`"AZDEVOPS_ADO_ENDPOINT" not specified`)
	}
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.NewClientFactory(t.Context(), ep, cred, nil)
	if err != nil {
		t.Fatal(err)
	}
}
