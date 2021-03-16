package provider

import (
	"context"
	"encoding/json"

	"github.com/chronark/terraform-provider-vercel/internal/vercel"
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

type response struct {
	User struct {
		UID             string `json:"uid"`
		Email           string `json:"email"`
		Name            string `json:"name"`
		Username        string `json:"username"`
		Avatar          string `json:"avatar"`
		PlatformVersion int    `json:"platformVersion"`
		Billing         struct {
			Plan        string      `json:"plan"`
			Period      interface{} `json:"period"`
			Trial       interface{} `json:"trial"`
			Cancelation interface{} `json:"cancelation"`
			Addons      interface{} `json:"addons"`
		} `json:"billing"`
		Bio      string `json:"bio"`
		Website  string `json:"website"`
		Profiles []struct {
			Service string `json:"service"`
			Link    string `json:"link"`
		} `json:"profiles"`
	} `json:"user"`
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*vercel.Client)

	res, err := client.Call("GET", "/www/user")
	if err != nil {
		return diag.FromErr(err)
	}
	defer res.Body.Close()

	var decodedResponse response
	err = json.NewDecoder(res.Body).Decode(&decodedResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("email", decodedResponse.User.Email)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", decodedResponse.User.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("username", decodedResponse.User.Username)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("avatar", decodedResponse.User.Avatar)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("platformversion", decodedResponse.User.PlatformVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("bio", decodedResponse.User.Bio)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("website", decodedResponse.User.Website)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("profiles", decodedResponse.User.Profiles)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(decodedResponse.User.UID)

	return diags
}
