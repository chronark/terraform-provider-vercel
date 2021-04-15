package provider

import (
	"context"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/env"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnv() *schema.Resource {
	return &schema.Resource{
		Description: "https://vercel.com/docs/api#endpoints/projects/get-project-environment-variables",

		CreateContext: resourceEnvCreate,
		ReadContext:   resourceEnvRead,
		UpdateContext: resourceEnvUpdate,
		DeleteContext: resourceEnvDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "The unique project identifier.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"team_id": {
				Description: "By default, you can access resources contained within your own user account. To access resources owned by a team, you can pass in the team ID",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
			},
			"type": {
				Description: "The type can be `plain`, `secret`, or `system`.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"id": {
				Description: "Unique id for this variable.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"key": {
				Description: "The name of the environment variable.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"value": {
				Description: "If the type is `plain`, a string representing the value of the environment variable. If the type is `secret`, the secret ID of the secret attached to the environment variable. If the type is `system`, the name of the System Environment Variable.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"target": {
				Description: "The target can be a list of `development`, `preview`, or `production`.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    3,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Description: "A number containing the date when the variable was created in milliseconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"updated_at": {
				Description: "A number containing the date when the variable was updated in milliseconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func toCreateOrUpdateEnv(d *schema.ResourceData) env.CreateOrUpdateEnv {

	// Casting each target because go does not allow typecasting from interface{} to []string
	targetList := d.Get("target").([]interface{})
	target := make([]string, len(targetList))
	for i := 0; i < len(target); i++ {
		target[i] = targetList[i].(string)
	}

	return env.CreateOrUpdateEnv{
		Type:   d.Get("type").(string),
		Key:    d.Get("key").(string),
		Value:  d.Get("value").(string),
		Target: target,
	}
}

func resourceEnvCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)

	payload := toCreateOrUpdateEnv(d)

	envID, err := client.Env.Create(d.Get("project_id").(string), payload, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(envID)

	return resourceEnvRead(ctx, d, meta)
}

func resourceEnvRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	id := d.Id()
	allEnvVariables, err := client.Env.Read(d.Get("project_id").(string), d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// Filter the current variable out of all existing ones
	var currentVar env.Env
	for _, envVar := range allEnvVariables {
		if envVar.ID == id {
			currentVar = envVar
			break
		}
	}

	err = d.Set("type", currentVar.Type)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("key", currentVar.Key)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("value", currentVar.Value)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("target", currentVar.Target)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("updated_at", currentVar.UpdatedAt)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("created_at", currentVar.CreatedAt)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceEnvUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)

	// Vercel expects an object with all 4 keys, so there's not point in checking for individual changes.
	if d.HasChanges("type", "key", "value", "taget") {

		projectID := d.Get("project_id").(string)
		envID := d.Id()
		payload := toCreateOrUpdateEnv(d)

		err := client.Env.Update(projectID, envID, payload, d.Get("team_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceEnvRead(ctx, d, meta)
}

func resourceEnvDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)

	projectID := d.Get("project_id").(string)
	envKey := d.Get("key").(string)

	err := client.Env.Delete(projectID, envKey, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diag.Diagnostics{}
}
