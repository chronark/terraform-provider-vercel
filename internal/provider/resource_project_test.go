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

func TestAccVercelProjectBasic(t *testing.T) {

	projectName, _ := uuid.GenerateUUID()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckVercelProjectDestroy(projectName),
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

func testAccCheckVercelProjectDestroy(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := vercel.New(os.Getenv("VERCEL_TOKEN"))

		for _, rs := range s.RootModule().Resources {
			if rs.Type != name {
				continue
			}

			project, err := client.Project.Read(rs.Primary.ID)
			if err == nil {
				_ = client.Project.Delete(project.ID)
				return fmt.Errorf("Project still exists on vercel")
			}

		}
		return nil
	}

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
