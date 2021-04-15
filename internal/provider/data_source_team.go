package provider

import (
	"context"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTeam() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves information related to an existing team. https://vercel.com/docs/api#endpoints/teams",
		ReadContext: dataSourceTeamRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"slug": {
				Required: true,
				Type:     schema.TypeString,
			},
			"name": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"creator_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"avatar": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"created": {
				Computed: true,
				Type:     schema.TypeInt,
			},
		},
	}
}

func dataSourceTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)

	slug := d.Get("slug").(string)

	team, err := client.Team.Read(slug)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("name", team.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("creator_id", team.CreatorId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("avatar", team.Avatar)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("created", team.Created.Unix())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(team.Id)

	return diag.Diagnostics{}
}
