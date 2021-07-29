package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/chronark/terraform-provider-vercel/pkg/util"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	projectApi "github.com/chronark/terraform-provider-vercel/pkg/vercel/project"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Description: "https://vercel.com/docs/api#endpoints/projects",

		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceProjectImportState,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Internal id of this project",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the project.",
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
			"git_repository": {
				Description: "The git repository that will be connected to the project. Any pushes to the specified connected git repository will be automatically deployed.",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "The git provider of the repository. Must be either `github`, `gitlab`, or `bitbucket`.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"repo": {
							Description: "The name of the git repository. For example: `chronark/terraform-provider-vercel`",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"domain": {
				Description: "Add a domain to the project by passing the project.",
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the production domain.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"redirect": {
							Description: "Target destination domain for redirect.",
							Type:        schema.TypeString,
							Optional:    true,
						}, "redirect_status_code": {
							Description: "The redirect status code (301, 302, 307, 308).",
							Type:        schema.TypeInt,
							Optional:    true,
						}, "git_branch": {
							Description: "it branch for the domain to be auto assigned to. The Project's production branch is the default (null).",
							Type:        schema.TypeString,
							Optional:    true,
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
				Default:     nil,
			},
			"serverless_function_region": {
				Description: "The region to deploy Serverless Functions in this project.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"root_directory": {
				Description: "The name of a directory or relative path to the source code of your project. When null is used it will default to the project root.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
			},
			"node_version": {
				Description: "The Node.js Version for this project.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"alias": {
				Description: "A list of production domains for the project.",
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceProjectImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// TODO: need to escape the slash inside project and team names?
	parts := strings.Split(d.Id(), "/")
	if len(parts) > 1 {
		teamID, projectID := parts[0], parts[1]
		d.Set("team_id", teamID)
		d.SetId(projectID)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)
	// Terraform does not have nested objects with different types yet, so I am using a `TypeList`
	// Here we have to typecast to list first and then take the first item and cast again.
	repo := d.Get("git_repository").([]interface{})[0].(map[string]interface{})
	project := projectApi.CreateProject{
		Name: d.Get("name").(string),
		GitRepository: struct {
			Type string `json:"type"`
			Repo string `json:"repo"`
		}{
			Type: repo["type"].(string),
			Repo: repo["repo"].(string),
		},
	}

	framework, frameworkSet := d.GetOk("framework")
	if frameworkSet {
		project.Framework = framework.(string)
	}
	publicSource, publicSourceSet := d.GetOk("public_source")
	if publicSourceSet {
		project.PublicSource = publicSource.(bool)
	}
	installCommand, installCommandSet := d.GetOk("install_command")
	if installCommandSet {
		project.InstallCommand = installCommand.(string)
	}
	buildCommand, buildCommandSet := d.GetOk("build_command")
	if buildCommandSet {
		project.BuildCommand = buildCommand.(string)
	}
	devCommand, devCommandSet := d.GetOk("dev_command")
	if devCommandSet {
		project.DevCommand = devCommand.(string)
	}
	outputDirectory, outputDirectorySet := d.GetOk("output_directory")
	if outputDirectorySet {
		project.OutputDirectory = outputDirectory.(string)
	}

	serverlessFunctionRegion, serverlessFunctionRegionSet := d.GetOk("serverless_function_region")
	if serverlessFunctionRegionSet {
		project.ServerlessFunctionRegion = serverlessFunctionRegion.(string)
	}
	rootDirectory, rootDirectorySet := d.GetOk("root_directory")
	if rootDirectorySet {
		project.RootDirectory = rootDirectory.(string)
	}
	nodeVersion, nodeVersionSet := d.GetOk("node_version")
	if nodeVersionSet {
		project.NodeVersion = nodeVersion.(string)

	}

	id, err := client.Project.Create(project, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	// repo := d.Get("git_repository").([]interface{})[0].(map[string]interface{})

	_, domainSet := d.GetOk("domain")
	if domainSet {

		rawDomain := d.Get("domain").([]interface{})[0].(map[string]interface{})
		domain := projectApi.Domain{
			Name: rawDomain["name"].(string),
			// Redirect:           rawDomain["redirect"].(string),
			// RedirectStatusCode: rawDomain["redirect_status_code"].(int),
		}

		err = client.Project.AddDomain(id, domain, d.Get("team_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(id)

	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	id := d.Id()

	project, err := client.Project.Read(id, d.Get("team_id").(string))
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

	aliases := make([]string, 0)
	for i := 0; i < len(project.Alias); i++ {
		aliases = append(aliases, project.Alias[i].Domain)
	}
	err = d.Set("alias", aliases)
	if err != nil {
		return diag.FromErr(err)
	}

	gitRepository := make([]map[string]interface{}, 1)
	gitRepository[0] = map[string]interface{}{
		"type": project.Link.Type,
	}
	switch project.Link.Type {
	case "gitlab":
		gitRepository[0]["repo"] = fmt.Sprintf("%s/%s", project.Link.ProjectNamespace, project.Link.ProjectName)
	case "github":
		gitRepository[0]["repo"] = fmt.Sprintf("%s/%s", project.Link.Org, project.Link.Repo)
	case "bitbucket":
		gitRepository[0]["repo"] = fmt.Sprintf("%s/%s", project.Link.Owner, project.Link.Slug)
	default:
		return diag.Errorf("Can't recognize '%s' repository type", project.Link.Type)
	}

	err = d.Set("git_repository", gitRepository)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)
	var update projectApi.UpdateProject

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

	err := client.Project.Update(d.Id(), update, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}


	if d.HasChange("domain") {
  	rawOldDomains, rawNewDomains := d.GetChange("domain")

	  var (
		  oldDomains []string
		  newDomains []string
	  )

	  for _, d := range rawNewDomains.([]interface{}) {
		  newDomains = append(newDomains, d.(map[string]interface{})["name"].(string))
	  }

	  for _, d := range rawOldDomains.([]interface{}) {
		  oldDomains = append(oldDomains, d.(map[string]interface{})["name"].(string))
	  }

	  toAdd := util.Difference(newDomains, oldDomains)
	  toRemove := util.Difference(oldDomains, newDomains)

	  for _, dom := range toAdd {
	  	err := client.Project.AddDomain(d.Id(), projectApi.Domain{Name: dom}, d.Get("team_id").(string))
			if err != nil {
				return diag.FromErr(err)
			}
	  }

	  for _, dom := range toRemove {
	  	err := client.Project.RemoveDomain(d.Id(), dom, d.Get("team_id").(string))
			if err != nil {
				return diag.FromErr(err)
			}
	  }
	}

	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)
	err := client.Project.Delete(d.Id(), d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
