package provider

import (
	"context"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/dns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNS() *schema.Resource {
	return &schema.Resource{
		Description: "https://vercel.com/docs/api#endpoints/dns\nCurrently this will only fetch the last 1000 records. Please create an issue if you require more.",

		CreateContext: resourceDNSCreate,
		ReadContext:   resourceDNSRead,
		DeleteContext: resourceDNSDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Description: "The domain for this DNS record",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"id": {
				Description: "The unique identifier of the dns record.",
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
			"type": {
				Description: "The type of record, it could be any valid DNS record.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "A subdomain name or an empty string for the root domain.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"value": {
				Description: "The record value.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"ttl": {
				Description: "The TTL value. Must be a number between 60 and 2147483647. Default value is 60.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				ForceNew:    true,
			},
			"creator": {
				Description: "The ID of the user who created the record or system if the record is an automatic record.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created": {
				Description: "The date when the record was created.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"updated": {
				Description: "The date when the record was updated.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"created_at": {
				Description: "The date when the record was created in milliseconds since the UNIX epoch.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"updated_at": {
				Description: "The date when the record was updated in milliseconds since the UNIX epoch.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourceDNSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)

	payload := dns.CreateRecord{
		Name:  d.Get("name").(string),
		Type:  d.Get("type").(string),
		Value: d.Get("value").(string),
		TTL:   d.Get("ttl").(int),
	}
	domain := d.Get("domain").(string)
	teamId := d.Get("team_id").(string)
	dnsId, err := client.DNS.Create(domain, payload, teamId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dnsId)

	return resourceDNSRead(ctx, d, meta)
}

func resourceDNSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	record, err := client.DNS.Read(d.Get("domain").(string), d.Id(), d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", record.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("type", record.Type)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("value", record.Value)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("creator", record.Creator)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("created", record.Created)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("updated", record.Updated)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("created_at", record.CreatedAt)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("updated_at", record.UpdatedAt)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceDNSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)

	err := client.DNS.Delete(d.Get("domain").(string), d.Id(), d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diag.Diagnostics{}
}
