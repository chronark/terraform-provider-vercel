package provider

import (
	"fmt"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/project"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

func TestAccVercelProject(t *testing.T) {

	projectName, _ := uuid.GenerateUUID()
	updatedProjectName, _ := uuid.GenerateUUID()
	var (

		// Holds the project fetched from vercel when we create it at the beginning
		actualProjectAfterCreation project.Project

		// Renaming or changing a variable should not result in the recreation of the project, so we expect to have the same id.
		actualProjectAfterUpdate project.Project

		// Used everywhere else
		actualProject project.Project
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckVercelProjectDestroy(projectName),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVercelProjectConfig(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectStateHasValues(
						"vercel_project.new", project.Project{Name: projectName},
					),
					testAccCheckVercelProjectExists("vercel_project.new", &actualProjectAfterCreation),
					testAccCheckActualProjectHasValues(&actualProjectAfterCreation, &project.Project{Name: projectName}),
				),
			},
			{
				Config: testAccCheckVercelProjectConfig(updatedProjectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVercelProjectExists("vercel_project.new", &actualProjectAfterUpdate),
					testAccCheckProjectStateHasValues(
						"vercel_project.new", project.Project{Name: updatedProjectName},
					),
					testAccCheckActualProjectHasValues(&actualProjectAfterUpdate, &project.Project{Name: updatedProjectName}),
					testAccCheckProjectWasNotRecreated(&actualProjectAfterCreation, &actualProjectAfterUpdate),
				),
			},
			{
				Config: testAccCheckVercelProjectConfigWithOverridenCommands(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVercelProjectExists("vercel_project.new", &actualProject),
					testAccCheckProjectStateHasValues(
						"vercel_project.new", project.Project{
							Name:            projectName,
							InstallCommand:  "echo install",
							BuildCommand:    "echo build",
							DevCommand:      "echo dev",
							OutputDirectory: "out",
						},
					),
					testAccCheckActualProjectHasValues(&actualProject, &project.Project{
						Name:            projectName,
						InstallCommand:  "echo install",
						BuildCommand:    "echo build",
						DevCommand:      "echo dev",
						OutputDirectory: "out",
					},
					),
				),
			},
		},
	})
}

// Combines multiple `resource.TestCheckResourceAttr` calls
func testAccCheckProjectStateHasValues(name string, want project.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tests := []resource.TestCheckFunc{

			resource.TestCheckResourceAttr(
				name, "install_command", want.InstallCommand),
			resource.TestCheckResourceAttr(
				name, "build_command", want.BuildCommand),
			resource.TestCheckResourceAttr(
				name, "dev_command", want.DevCommand),
			resource.TestCheckResourceAttr(
				name, "output_directory", want.OutputDirectory),
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

// Chaning the name or value of a project should not result in a recreation meaning the UID assigned by vercel
// should not have changed.
func testAccCheckProjectWasNotRecreated(s1, s2 *project.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if s1.ID != s2.ID {
			return fmt.Errorf("Expected same IDs but they are not the same.")
		}
		return nil
	}
}

func testAccCheckActualProjectHasValues(actual *project.Project, want *project.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if actual.Name != want.Name {
			return fmt.Errorf("name does not match, expected: %s, got: %s", want.Name, actual.Name)
		}
		if actual.ID == "" {
			return fmt.Errorf("ID is empty")
		}

		if actual.InstallCommand != want.InstallCommand {
			return fmt.Errorf("install_command does not match: expected: %s, got: %s", want.InstallCommand, actual.InstallCommand)
		}
		if actual.BuildCommand != want.BuildCommand {
			return fmt.Errorf("build_command does not match: expected: %s, got: %s", want.BuildCommand, actual.BuildCommand)
		}
		if actual.DevCommand != want.DevCommand {
			return fmt.Errorf("dev_command does not match: expected: %s, got: %s", want.DevCommand, actual.DevCommand)
		}
		if actual.OutputDirectory != want.OutputDirectory {
			return fmt.Errorf("output_directory does not match: expected: %s, got: %s", want.OutputDirectory, actual.OutputDirectory)
		}

		return nil
	}
}

// Test whether the project was destroyed properly and finishes the job if necessary
func testAccCheckVercelProjectDestroy(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := vercel.New(os.Getenv("VERCEL_TOKEN"))

		for _, rs := range s.RootModule().Resources {
			if rs.Type != name {
				continue
			}

			project, err := client.Project.Read(rs.Primary.ID, "")
			if err == nil {
				message := "Project was not deleted from vercel during terraform destroy."
				deleteErr := client.Project.Delete(project.Name, "")
				if deleteErr != nil {
					return fmt.Errorf(message+" Automated removal did not succeed. Please manually remove @%s. Error: %w", project.Name, err)
				}
				return fmt.Errorf(message + " It was removed now.")
			}

		}
		return nil
	}

}
func testAccCheckVercelProjectConfig(name string) string {
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

func testAccCheckVercelProjectConfigWithOverridenCommands(name string) string {
	return fmt.Sprintf(`
	resource "vercel_project" "new" {
		name = "%s"
		git_repository {
			type = "github"
			repo = "chronark/mercury"
		}
		install_command  = "echo install"
		build_command 	 = "echo build"
		dev_command 	 = "echo dev"
		output_directory = "out"
	}
	`, name)
}

func testAccCheckVercelProjectExists(n string, actual *project.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s in %+v", n, s.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No project set")
		}

		project, err := vercel.New(os.Getenv("VERCEL_TOKEN")).Project.Read(rs.Primary.ID, "")
		if err != nil {
			return err
		}
		*actual = project
		return nil
	}
}
