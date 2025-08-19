package azdevops_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	"github.com/magodo/azdevops"
)

func TestNewConnection_Anonymous(t *testing.T) {
	ep, ok := os.LookupEnv("AZDEVOPS_ADO_ENDPOINT")
	if !ok {
		t.Skip(`"AZDEVOPS_ADO_ENDPOINT" not specified`)
	}
	_, err := azdevops.NewConnection(t.Context(), ep, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestConnection_NewClient(t *testing.T) {
	ep, ok := os.LookupEnv("AZDEVOPS_ADO_ENDPOINT")
	if !ok {
		t.Skip(`"AZDEVOPS_ADO_ENDPOINT" not specified`)
	}
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := azdevops.NewConnection(t.Context(), ep, &azdevops.Credential{AADTokenCredential: cred}, nil)
	if err != nil {
		t.Fatal(err)
	}

	// New default client
	if _, err := conn.NewDefaultClient(t.Context()); err != nil {
		t.Fatal(err)
	}

	// New a release are client
	if _, err := conn.NewAreaClient(t.Context(), uuid.MustParse("efc2f575-36ef-48e9-b672-0c6fb4a48ac5")); err != nil {
		t.Fatal(err)
	}
}
