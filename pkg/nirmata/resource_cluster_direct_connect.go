package nirmata

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/nirmata/go-client/pkg/client"
)

func resourceClusterDirectConnect() *schema.Resource {

	return &schema.Resource{
		Create: resourceClusterDirectConnectCreate,
		Read:   resourceClusterDirectConnectRead,
		Update: resourceClusterDirectConnectUpdate,
		Delete: resourceClusterDirectConnectDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"host_group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceClusterDirectConnectCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(client.Client)

	name := d.Get("name").(string)
	policy := d.Get("policy").(string)
	hostGroup := d.Get("host_group").(string)

	clusterData := map[string]interface{}{
		"name":           name,
		"mode":           "managed",
		"policySelector": policy,
		"hostGroupSelector": map[string]interface{}{
			"matchLabels": map[string]interface{}{
				"name": hostGroup,
			},
		},
	}

	resultData, err := apiClient.PostFromJSON(client.ServiceClusters, "HostCluster", clusterData, nil)
	if err != nil {
		return err
	}

	clusterID := resultData["id"].(string)
	d.SetId(clusterID)
	return nil
}

func resourceClusterDirectConnectRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceClusterDirectConnectUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceClusterDirectConnectDelete(d *schema.ResourceData, m interface{}) error {
	return delete(d, m, client.ServiceClusters, "HostCluster")
}
