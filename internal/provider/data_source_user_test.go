package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceUser(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vercel_user.me", "username", "chronark"),
				),
			},
		},
	})
}

const testAccDataSourceUser = `data "vercel_user" "me" {}`
