package provider

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"vercel": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("Unable to create new provider: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	require.NotEmpty(t, os.Getenv("VERCEL_TOKEN"))
}
