package provider

import (
	"context"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Retrieves information related to the authenticated user. https://vercel.com/docs/api#endpoints/user/get-the-authenticated-user",
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"email": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"name": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"username": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"avatar": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"platformversion": {
				Computed: true,
				Type:     schema.TypeInt,
			},

			"bio": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"website": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"profiles": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"link": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)

	user, err := client.User.Read()
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("email", user.Email)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", user.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("username", user.Username)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("avatar", user.Avatar)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("platformversion", user.PlatformVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("bio", user.Bio)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("website", user.Website)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("profiles", user.Profiles)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.UID)

	return diag.Diagnostics{}
}
