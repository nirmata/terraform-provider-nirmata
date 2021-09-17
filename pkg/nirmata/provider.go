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
				DefaultFunc: schema.EnvDefaultFunc("NIRMATA_URL", "https://nirmata.io"),
				Description: "Nirmata URL (HTTPS) address",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"nirmata_cluster":                   resourceManagedCluster(),
			"nirmata_cluster_imported":          resourceClusterImported(),
			"nirmata_cluster_type_aks":          resourceAksClusterType(),
			"nirmata_cluster_type_oke":          resourceOkeClusterType(),
			"nirmata_cluster_type_gke":          resourceGkeClusterType(),
			"nirmata_cluster_type_eks":          resourceEksClusterType(),
			"nirmata_cluster_type_registered":   resourceRegisteredClusterType(),
			"nirmata_cluster_registered":        resourceClusterRegistered(),
			"nirmata_cluster_direct_connect":    resourceClusterDirectConnect(),
			"nirmata_host_group_direct_connect": resourceHostGroupDirectConnect(),
			"nirmata_environment":               resourceEnvironment(),
			"nirmata_environment_type":          resourceEnvironmentType(),
			"nirmata_cluster_addons":            resoureClusterAddOn(),
			"nirmata_aws_role_credentials":      resoureAwsRoleCredentials(),
			"nirmata_catalog":                   resourceCatalog(),
			"nirmata_git_application":           resourceGitApplication(),
			"nirmata_run_application":           resourceRunApplication(),
			"nirmata_promote_version":           resourcePromoteVersion(),
			"nirmata_helm_application":          resourceHelmApplication(),
			"nirmata_catalog_application":       resourceCatalogApplication(),
		},

		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("url").(string)
	token := d.Get("token").(string)

	return client.NewClient(url, token, nil, false), nil
}
