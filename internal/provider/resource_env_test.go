package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/env"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/project"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVercelEnv_import(t *testing.T) {

	projectName, _ := uuid.GenerateUUID()
	var (
		actualProject project.Project
		actualEnv     env.Env
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckVercelProjectDestroy(projectName),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVercelEnvConfig(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectStateHasValues(
						"vercel_project.new", project.Project{Name: projectName},
					),
					testAccCheckVercelProjectExists("vercel_project.new", &actualProject),
					testAccCheckVercelEnvExists("vercel_env.new", &actualEnv),
					testAccVercelEnvImport(&actualProject, &actualEnv),
				),
			},
		},
	})
}

func testAccCheckVercelEnvConfig(name string) string {
	return fmt.Sprintf(`
	resource "vercel_project" "new" {
		name = "%s"
		git_repository {
			type = "github"
			repo = "chronark/terraform-provider-vercel"
		}
	}

	resource "vercel_env" "new" {
		key         = "FOO"
		value       = "BAR"

		project_id  = vercel_project.new.id
		target      = [ "production" ]
		type        = "plain"
	}
	`, name)
}

func testAccCheckVercelEnvExists(n string, actual *env.Env) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s in %+v", n, s.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No project set")
		}

		projectID := rs.Primary.Attributes["project_id"]
		teamID := rs.Primary.Attributes["team_id"]

		envs, err := vercel.New(os.Getenv("VERCEL_TOKEN")).Env.Read(projectID, teamID)
		if err != nil {
			return err
		}
		for _, env := range envs {
			if env.ID == rs.Primary.ID {
				*actual = env
				return nil
			}
		}

		return fmt.Errorf("Could not find env '%s' in project '%s'", rs.Primary.ID, projectID)
	}
}

func testAccVercelEnvImport(srcProject *project.Project, srcEnv *env.Env) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		data := resourceEnv().Data(nil)
		data.SetId(srcProject.Name + "/" + srcEnv.Key)
		ds, err := resourceEnv().Importer.StateContext(context.Background(), data, vercel.New(os.Getenv("VERCEL_TOKEN")))
		if err != nil {
			return err
		}
		if len(ds) != 1 {
			return fmt.Errorf("Expected 1 instance state from importer function. Got %d", len(ds))
		}

		if ds[0].Id() != srcEnv.ID {
			return fmt.Errorf("Imported env ID. Expected '%s'. Actual '%s'.", srcEnv.ID, ds[0].Id())
		}
		if ds[0].Get("project_id") != srcProject.ID {
			return fmt.Errorf("Imported env project ID. Expected '%s'. Actual '%s'.", srcProject.ID, ds[0].Get("project_id").(string))
		}
		return nil
	}
}
