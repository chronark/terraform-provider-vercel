package provider

import (
	"context"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/secret"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		Description: "https://vercel.com/docs/api#endpoints/secrets",

		CreateContext: resourceSecretCreate,
		ReadContext:   resourceSecretRead,
		DeleteContext: resourceSecretDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The unique identifier of the secret.",
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
			"name": {
				Description: "The name of the secret.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"value": {
				Description: "The value of the new secret.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
			},

			"user_id": {
				Description: "The unique identifier of the user who created the secret.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "A number containing the date when the variable was created in milliseconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourceSecretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)

	payload := secret.CreateSecret{
		Name:  d.Get("name").(string),
		Value: d.Get("value").(string),
	}

	secretID, err := client.Secret.Create(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(secretID)

	return resourceSecretRead(ctx, d, meta)
}

func resourceSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	id := d.Id()

	secret, err := client.Secret.Read(id, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", secret.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("team_id", secret.TeamID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("user_id", secret.UserID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("created_at", secret.CreatedAt)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceSecretDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)

	err := client.Secret.Delete(d.Get("name").(string), d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diag.Diagnostics{}
}
