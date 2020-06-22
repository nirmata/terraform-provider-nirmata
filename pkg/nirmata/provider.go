package nirmata

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

// Provider returns the Nirmata Terraform Provider
func Provider() *schema.Provider {

	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NIRMATA_TOKEN", nil),
				Description: "Nirmata API Access Token",
				Sensitive:   true,
			},

			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NIRMATA_URL", "https://devtest2.nirmata.io/"),
				Description: "Nirmata URL (HTTPS) address",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"nirmata_host_group_direct_connect": resourceHostGroupDirectConnect(),
			"nirmata_cluster_direct_connect":    resourceClusterDirectConnect(),
			"nirmata_ProviderManaged_cluster":   resourceProviderManagedCluster(),
			"nirmata_aks_clusterType":           resourceAksClusterType(),
			"nirmata_oke_clusterType":           resourceOkeClusterType(),
			"nirmata_gke_clusterType":           resourceGkeClusterType(),
			"nirmata_eks_clusterType":           resourceEksClusterType(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("url").(string)
	token := d.Get("token").(string)

	return client.NewClient(url, token, nil, false), nil
}
