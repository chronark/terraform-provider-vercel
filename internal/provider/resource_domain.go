package provider

import (
	"context"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Description: "https://vercel.com/docs/api#endpoints/projects/get-project-environment-variables",

		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		DeleteContext: resourceDomainDelete,

		Schema: map[string]*schema.Schema{
			"team_id": {
				Description: "By default, you can access resources contained within your own user account. To access resources owned by a team, you can pass in the team ID",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
			},
			"name": {
				Description: "The name of the production domain.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"id": {
				Description: "Unique id for this variable.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"service_type": {
				Description: "The type of service the domain is handled by. external if the DNS is externally handled, or zeit.world if handled with Vercel.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"ns_verified_at": {
				Description: "The date at which the domain's nameservers were verified based on the intended set.",
				Type:        schema.TypeInt,
				Computed:    true,
			},

			"txt_verified_at": {
				Description: "The date at which the domain's TXT DNS record was verified.",
				Type:        schema.TypeInt,
				Computed:    true,
			},

			"cdn_enabled": {
				Description: "Whether the domain has the Vercel Edge Network enabled or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},

			"created_at": {
				Description: "The date when the domain was created in the registry.",
				Type:        schema.TypeInt,
				Computed:    true,
			},

			"expires_at": {
				Description: "The date at which the domain is set to expire. null if not bought with Vercel.",
				Type:        schema.TypeInt,
				Computed:    true,
			},

			"bought_at": {
				Description: "If it was purchased through Vercel, the date when it was purchased.",
				Type:        schema.TypeInt,
				Computed:    true,
			},

			"transfer_started_at": {
				Description: "If transferred into Vercel, The date when the domain transfer was initiated",
				Type:        schema.TypeInt,
				Computed:    true,
			},

			"transferred_at": {
				Description: "The date at which the domain was successfully transferred into Vercel. null if the transfer is still processing or was never transferred in.",
				Type:        schema.TypeInt,
				Computed:    true,
			},

			"verification_record": {
				Description: "The ID of the verification record in the registry.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"verified": {
				Description: "If the domain has the ownership verified.",
				Type:        schema.TypeBool,
				Computed:    true,
			},

			"nameservers": {
				Description: "A list of the current nameservers of the domain.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"intended_nameservers": {
				Description: "A list of the intended nameservers for the domain to point to Vercel DNS.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)

	id, err := client.Domain.Create(d.Get("name").(string), d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return resourceDomainRead(ctx, d, meta)
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	domain, err := client.Domain.Read(d.Get("name").(string), d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(domain.ID)
	err = d.Set("service_type", domain.ServiceType)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("ns_verified_at", domain.NsVerifiedAt)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("txt_verified_at", domain.TxtVerifiedAt)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("cdn_enabled", domain.CdnEnabled)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("created_at", domain.CreatedAt)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("expires_at", domain.ExpiresAt)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("bought_at", domain.BoughtAt)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("transfer_started_at", domain.TransferredAt)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("transferred_at", domain.TransferredAt)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("verification_record", domain.VerificationRecord)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("verified", domain.Verified)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("nameservers", domain.Nameservers)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("intended_nameservers", domain.IntendedNameservers)
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)

	domainName := d.Get("name").(string)

	err := client.Domain.Delete(domainName, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diag.Diagnostics{}
}
