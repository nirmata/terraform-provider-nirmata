package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/terraform-provider-nirmata/pkg/client"
)

// Provider returns a Nirmata terraform.Provider
func Provider() *schema.Provider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{

			"token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				DefaultFunc: schema.EnvDefaultFunc("NIRMATA_TOKEN", nil),
				Description: "Nirmata API Access Token",
			},

			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NIRMATA_URL", "https://nirmata.io"),
				Description: "Nirmata URL (HTTPS) address",
			},
		},

		ResourcesMap: map[string]*schema.Resource{},

		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("url").(string)
	token := d.Get("token").(string)

	return client.NewClient(url, token, nil, false), nil
}
