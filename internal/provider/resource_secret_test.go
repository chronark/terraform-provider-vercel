package provider

import (
	"fmt"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

func TestAccVercelSecret(t *testing.T) {

	secretName, _ := uuid.GenerateUUID()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckVercelSecretDestroy(secretName),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVercelSecretConfig(secretName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVercelSecretExists("vercel_secret.new"),
					resource.TestCheckResourceAttr(
						"vercel_secret.new", "name", secretName),
				),
			},
		},
	})
}
// Test whether the secret was destroyed properly and finishes the job if necessary
func testAccCheckVercelSecretDestroy(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := vercel.New(os.Getenv("VERCEL_TOKEN"))

		for _, rs := range s.RootModule().Resources {
			if rs.Type != name {
				continue
			}

			secret, err := client.Secret.Read(rs.Primary.ID)
			if err == nil {
				message := "Secret was not deleted from vercel during terraform destroy."
				deleteErr := client.Secret.Delete(secret.Name)
				if deleteErr != nil {
					return fmt.Errorf(message + " Automated removal did not succeed. Please manually remove @%s. Error: %w", secret.Name, err)
				}
				return fmt.Errorf(message + " It was removed now.")
			}

		}
		return nil
	}

}
func testAccCheckVercelSecretConfig(name string) string {
	return fmt.Sprintf(`
	resource "vercel_secret" "new" {
		name = "%s"
		value = "secret value"
	}
	`, name)
}

func testAccCheckVercelSecretExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s in %+v", n, s.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No secret set")
		}

		return nil
	}
}
