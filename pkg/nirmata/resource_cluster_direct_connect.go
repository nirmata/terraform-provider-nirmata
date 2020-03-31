package nirmata

import (
	"strings"

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
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
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
	updateClusterData(d, resultData)

	return nil
}

func updateClusterData(d *schema.ResourceData, data map[string]interface{}) {
	updateStateAndStatus(d, data)
}

func resourceClusterDirectConnectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(client.Client)
	id := clientID(d, client.ServiceClusters, "HostCluster")
	data, err := apiClient.Get(id, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}

		return err
	}

	updateClusterData(d, data)
	return nil
}

func resourceClusterDirectConnectUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceClusterDirectConnectDelete(d *schema.ResourceData, m interface{}) error {
	params := map[string]string{"action": "delete"}
	return delete(d, m, client.ServiceClusters, "HostCluster", params)
}
