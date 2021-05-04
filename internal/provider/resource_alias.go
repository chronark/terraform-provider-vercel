package provider

import (
	"context"
	"fmt"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/alias"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlias() *schema.Resource {
	return &schema.Resource{
		Description: "https://vercel.com/docs/api#endpoints/projects/add-a-domain-to-a-project",

		CreateContext: resourceAliasCreate,
		ReadContext:   resourceAliasRead,
		UpdateContext: resourceAliasUpdate,
		DeleteContext: resourceAliasDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "The unique Project identifier.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"id": {
				Description: "The unique identifier of the alias.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"team_id": {
				Description: "By default, you can access resources contained within your own user account. To access resources owned by a team, you can pass in the team ID",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
			},
			"domain": {
				Description: "The name of the production domain.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"redirect": {
				Description: "Target destination domain for redirect",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceAliasCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)

	projectId := d.Get("project_id").(string)
	domain := d.Get("domain").(string)
	teamId := d.Get("team_id").(string)

	payload := alias.CreateOrUpdateAlias{
		Domain: domain,
	}
	redirect, set := d.GetOk("redirect")
	if set {
		payload.Redirect = redirect.(string)
	}
	err := client.Alias.Create(projectId, payload, teamId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s-%s", projectId, domain))

	return resourceAliasRead(ctx, d, meta)
}
func resourceAliasRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)
	projectId := d.Get("project_id").(string)
	domain := d.Get("domain").(string)
	teamId := d.Get("team_id").(string)
	alias, err := client.Alias.Read(
		projectId, domain, teamId,
	)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("domain", alias.Domain)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("redirect", alias.Redirect)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}

}

func resourceAliasUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	payload := alias.CreateOrUpdateAlias{}

	if d.HasChange("domain") {
		payload.Domain = d.Get("domain").(string)
	}
	if d.HasChange("redirect") {
		payload.Redirect = d.Get("redirect").(string)
	}

	err := client.Alias.Update(d.Get("project_id").(string), payload, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceAliasRead(ctx, d, meta)

}

func resourceAliasDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)

	err := client.Alias.Delete(d.Get("project_id").(string), d.Get("domain").(string), d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diag.Diagnostics{}
}
