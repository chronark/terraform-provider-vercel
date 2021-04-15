package provider

import (
	"fmt"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/secret"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

func TestAccVercelSecret(t *testing.T) {

	secretName, _ := uuid.GenerateUUID()
	secretValue, _ := uuid.GenerateUUID()
	updatedSecretValue, _ := uuid.GenerateUUID()
	var (

		// Holds the secret fetched from vercel when we create it at the beginning
		actualSecretAfterCreation secret.Secret

		// Renaming or chaning a variable results in the recreation of the secret, so we expect this value to have a different id.
		actualSecretAfterUpdate secret.Secret
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckVercelSecretDestroy(secretName),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVercelSecretConfig(secretName, secretValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretStateHasValues(
						"vercel_secret.new", secret.CreateSecret{Name: secretName, Value: secretValue},
					),
					testAccCheckVercelSecretExists("vercel_secret.new", &actualSecretAfterCreation),
					testAccCheckActualSecretHasValues(&actualSecretAfterCreation, &secret.Secret{Name: secretName}),
				),
			},
			{
				Config: testAccCheckVercelSecretConfig(secretName, updatedSecretValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVercelSecretExists("vercel_secret.new", &actualSecretAfterUpdate),
					testAccCheckSecretStateHasValues(
						"vercel_secret.new", secret.CreateSecret{Name: secretName, Value: updatedSecretValue},
					),
					testAccCheckActualSecretHasValues(&actualSecretAfterUpdate, &secret.Secret{Name: secretName}),
					testAccCheckSecretWasRecreated(&actualSecretAfterCreation, &actualSecretAfterUpdate),
				),
			},
		},
	})
}

// Combines multiple `resource.TestCheckResourceAttr` calls
func testAccCheckSecretStateHasValues(name string, want secret.CreateSecret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tests := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(
				name, "name", want.Name),
			resource.TestCheckResourceAttr(
				name, "value", want.Value),
		}

		for _, test := range tests {
			err := test(s)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// Chaning the name or value of a secret results in a recreation meaning the UID assigned by vercel
// should have changed.
func testAccCheckSecretWasRecreated(s1, s2 *secret.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if s1.UID == s2.UID {
			return fmt.Errorf("Expected different UIDs but they are the same.")
		}
		return nil
	}
}

func testAccCheckActualSecretHasValues(actual *secret.Secret, want *secret.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if actual.Name != want.Name {
			return fmt.Errorf("name is not correct, expected: %s, got: %s", want.Name, actual.Name)
		}
		if actual.UID == "" {
			return fmt.Errorf("UID is empty")
		}

		return nil
	}
}

// Test whether the secret was destroyed properly and finishes the job if necessary
func testAccCheckVercelSecretDestroy(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := vercel.New(os.Getenv("VERCEL_TOKEN"))

		for _, rs := range s.RootModule().Resources {
			if rs.Type != name {
				continue
			}

			secret, err := client.Secret.Read(rs.Primary.ID, "")
			if err == nil {
				message := "Secret was not deleted from vercel during terraform destroy."
				deleteErr := client.Secret.Delete(secret.Name, "")
				if deleteErr != nil {
					return fmt.Errorf(message+" Automated removal did not succeed. Please manually remove @%s. Error: %w", secret.Name, err)
				}
				return fmt.Errorf(message + " It was removed now.")
			}

		}
		return nil
	}

}
func testAccCheckVercelSecretConfig(name string, value string) string {
	return fmt.Sprintf(`
	resource "vercel_secret" "new" {
		name = "%s"
		value = "%s"
	}
	`, name, value)
}

func testAccCheckVercelSecretExists(n string, actual *secret.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s in %+v", n, s.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No secret set")
		}

		secret, err := vercel.New(os.Getenv("VERCEL_TOKEN")).Secret.Read(rs.Primary.ID, "")
		if err != nil {
			return err
		}
		*actual = secret
		return nil
	}
}
