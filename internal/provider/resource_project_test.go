package provider

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccVercelProjectBasic(t *testing.T) {

	projectName, _ := uuid.GenerateUUID()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckVercelProjectConfigBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVercelProjectExists("vercel_project.new"),
					resource.TestCheckResourceAttr(
						"vercel_project.new", "name", projectName),
				),
			},
		},
	})
}

func testAccCheckVercelProjectConfigBasic(name string) string {
	return fmt.Sprintf(`
	resource "vercel_project" "new" {
		name = "%s"
		git_repository {
			type = "github"
			repo = "chronark/mercury"
		}
	}
	`, name)
}

func testAccCheckVercelProjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s in %+v", n, s.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No project set")
		}

		return nil
	}
}
