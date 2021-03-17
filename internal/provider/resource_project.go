package provider

import (
	"context"
	"log"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/project"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the project.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"git_repository": {
				Description: "The Git Repository that will be connected to the project. Any pushes to the specified connected Git Repository will be automatically deployed.",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "The Git Provider of the repository. Must be either `github`, `gitlab`, or `bitbucket`.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"repo": {
							Description: "The name of the Git Repository.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"account_id": {
				Description: "The unique ID of the user or team the project belongs to.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "A number containing the date when the project was created in milliseconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"updated_at": {
				Description: "A number containing the date when the project was updated in milliseconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"env": {
				Description: "A list of environment variables configured for the project.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "`secret` or `public`",
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
						},
						"id": {
							Description: "Unique id for this variable.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"key": {
							Description: "The name of this variable",
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
						},
						"value": {
							Description: "The value of this variable.",
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
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
				},
			},
			"framework": {
				Description: "The framework that is being used for this project. When null is used no framework is selected.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"public_source": {
				Description: " Specifies whether the source code and logs of the deployments for this project should be public or not.",
				Type:        schema.TypeBool,
				Optional:    true,
			},

			"install_command": {
				Description: "The install command for this project. When null is used this value will be automatically detected.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"build_command": {
				Description: "The build command for this project. When null is used this value will be automatically detected.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"dev_command": {
				Description: "The dev command for this project. When null is used this value will be automatically detected.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"output_directory": {
				Description: "The output directory of the project. When null is used this value will be automatically detected.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"serverless_function_region": {
				Description: "The region to deploy Serverless Functions in this project.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"root_directory": {
				Description: "The name of a directory or relative path to the source code of your project. When null is used it will default to the project root.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"node_version": {
				Description: "The Node.js Version for this project.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)
	// Terraform does not have nested objects with different types yet, so I am using a `TypeList`
	// Here we have to typecast to list first and then take the first item and cast again.
	repo := d.Get("git_repository").([]interface{})[0].(map[string]interface{})
	project := project.CreateProject{
		Name: d.Get("name").(string),
		GitRepository: struct {
			Type string `json:"type"`
			Repo string `json:"repo"`
		}{
			Type: repo["type"].(string),
			Repo: repo["repo"].(string),
		},
	}

	framework, isSet := d.GetOk("framework")
	if isSet {
		project.Framework = framework.(string)
	}

	id, err := client.Project.Create(project)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	id := d.Id()

	project, err := client.Project.Read(id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", project.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("account_id", project.AccountID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("created_at", project.CreatedAt)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("updated_at", project.UpdatedAt)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("env", project.Env)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("framework", project.Framework)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("public_source", project.PublicSource)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("install_command", project.InstallCommand)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("build_command", project.BuildCommand)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("dev_command", project.DevCommand)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("output_directory", project.OutputDirectory)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("serverless_function_region", project.ServerlessFunctionRegion)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("root_directory", project.RootDirectory)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("node_version", project.NodeVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)
	var update project.UpdateProject
	if d.HasChange("name") {
		update.Name = d.Get("name").(string)
	}
	if d.HasChange("framework") {
		update.Framework = d.Get("framework").(string)
	}
	if d.HasChange("public_source") {
		update.PublicSource = d.Get("public_source").(bool)
	}
	if d.HasChange("install_command") {
		update.InstallCommand = d.Get("install_command").(string)
	}
	if d.HasChange("build_command") {
		update.BuildCommand = d.Get("build_command").(string)
	}
	if d.HasChange("dev_command") {
		update.DevCommand = d.Get("dev_command").(string)
	}
	if d.HasChange("output_directory") {
		update.OutputDirectory = d.Get("output_directory").(string)
	}
	if d.HasChange("serverless_function_region") {
		update.ServerlessFunctionRegion = d.Get("serverless_function_region").(string)
	}
	if d.HasChange("root_directory") {
		update.RootDirectory = d.Get("root_directory").(string)
	}
	if d.HasChange("node_version") {
		update.NodeVersion = d.Get("node_version").(string)
	}
	err := client.Project.Update(d.Id(), update)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)
	err := client.Project.Delete(d.Id())
	if err != nil {
		log.Println("THE ERROR WAS HERE1")
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
