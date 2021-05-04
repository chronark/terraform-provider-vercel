package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/domain"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVercelDomain(t *testing.T) {

	domainName := "acceptancetestdomainone.com"
	updatedDomainName := "acceptancetestdomaintwo.com"
	var (

		// Holds the domain fetched from vercel when we create it at the beginning
		actualDomainAfterCreation domain.Domain

		// Renaming or changing a variable should not result in the recreation of the domain, so we expect to have the same id.
		actualDomainAfterUpdate domain.Domain
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckVercelDomainDestroy(domainName),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVercelDomainConfig(domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainStateHasValues(
						"vercel_domain.new", domain.Domain{Name: domainName},
					),
					testAccCheckVercelDomainExists("vercel_domain.new", &actualDomainAfterCreation),
					testAccCheckActualDomainHasValues(&actualDomainAfterCreation, &domain.Domain{Name: domainName}),
				),
			},
			{
				Config: testAccCheckVercelDomainConfig(updatedDomainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVercelDomainExists("vercel_domain.new", &actualDomainAfterUpdate),
					testAccCheckDomainStateHasValues(
						"vercel_domain.new", domain.Domain{Name: updatedDomainName},
					),
					testAccCheckActualDomainHasValues(&actualDomainAfterUpdate, &domain.Domain{Name: updatedDomainName}),
					testAccCheckDomainWasRecreated(&actualDomainAfterCreation, &actualDomainAfterUpdate),
				),
			},
		},
	})
}

// Combines multiple `resource.TestCheckResourceAttr` calls
func testAccCheckDomainStateHasValues(name string, want domain.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tests := []resource.TestCheckFunc{

			resource.TestCheckResourceAttr(
				name, "name", want.Name),
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

func testAccCheckDomainWasRecreated(s1, s2 *domain.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if s1.ID == s2.ID {
			return fmt.Errorf("Expected different IDs but they are the same.")
		}
		return nil
	}
}

func testAccCheckActualDomainHasValues(actual *domain.Domain, want *domain.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if actual.Name != want.Name {
			return fmt.Errorf("name does not match, expected: %s, got: %s", want.Name, actual.Name)
		}
		if actual.ID == "" {
			return fmt.Errorf("ID is empty")
		}

		return nil
	}
}

// Test whether the domain was destroyed properly and finishes the job if necessary
func testAccCheckVercelDomainDestroy(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := vercel.New(os.Getenv("VERCEL_TOKEN"))

		for _, rs := range s.RootModule().Resources {
			if rs.Type != name {
				continue
			}

			domain, err := client.Domain.Read(rs.Primary.ID, "")
			if err == nil {
				message := "Domain was not deleted from vercel during terraform destroy."
				deleteErr := client.Domain.Delete(domain.Name, "")
				if deleteErr != nil {
					return fmt.Errorf(message+" Automated removal did not succeed. Please manually remove @%s. Error: %w", domain.Name, err)
				}
				return fmt.Errorf(message + " It was removed now.")
			}

		}
		return nil
	}

}
func testAccCheckVercelDomainConfig(name string) string {
	return fmt.Sprintf(`
	resource "vercel_domain" "new" {
		name = "%s"
	}
	`, name)
}

func testAccCheckVercelDomainExists(n string, actual *domain.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s in %+v", n, s.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No domain set")
		}

		domain, err := vercel.New(os.Getenv("VERCEL_TOKEN")).Domain.Read(rs.Primary.Attributes["name"], "")
		if err != nil {
			return err
		}
		*actual = domain
		return nil
	}
}
