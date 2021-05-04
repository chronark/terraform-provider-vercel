package provider

import (
	"context"
	"fmt"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"token": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("VERCEL_TOKEN", nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"vercel_user": dataSourceUser(),
				"vercel_team": dataSourceTeam(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"vercel_env":     resourceEnv(),
				"vercel_project": resourceProject(),
				"vercel_secret":  resourceSecret(),
				"vercel_domain":  resourceDomain(),
				"vercel_dns":     resourceDNS(),
				"vercel_alias":   resourceAlias(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		token := d.Get("token").(string)

		if token == "" {
			return nil, diag.FromErr(fmt.Errorf("vercel token is not set, set manually or via `VERCEL_TOKEN` "))
		}

		client := vercel.New(token)

		return client, diag.Diagnostics{}
	}
}
