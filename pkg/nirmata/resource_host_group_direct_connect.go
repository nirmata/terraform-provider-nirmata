package nirmata

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resourceHostGroupDirectConnect() *schema.Resource {

	return &schema.Resource{
		Create: resourceHostGroupDirectConnectCreate,
		Read:   resourceHostGroupDirectConnectRead,
		Update: resourceHostGroupDirectConnectUpdate,
		Delete: resourceHostGroupDirectConnectDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"curl_script": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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

func resourceHostGroupDirectConnectCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(client.Client)
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		name := d.Get("name").(string)
		labels := d.Get("labels").(map[string]interface{})

		cpID, err := apiClient.QueryByName(client.ServiceConfig, "CloudProvider", "Direct Connect")
		if err != nil {
			return resource.RetryableError(err.OrigErr())
		}

		hg := map[string]interface{}{
			"name":   name,
			"parent": cpID.UUID(),
			"labels": labels,
		}

		data, err := apiClient.PostFromJSON(client.ServiceConfig, "HostGroup", hg, nil)
		if err != nil {
			return resource.RetryableError(err.OrigErr())
		}

		hgID := data["id"].(string)
		d.SetId(hgID)

		hgID := data["id"].(string)
		d.SetId(hgID)
	
		scriptPath := "nirmata-host-agent/setup-nirmata-agent.sh"
		curlScript := fmt.Sprintf("sudo curl -sSL %s/%s | sudo sh -s -- --cloud other --hostgroup %s",
			apiClient.Address(), scriptPath, hgID)
		d.Set("curl_script", curlScript)
	
		updateHostGroupData(d, data)

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func updateHostGroupData(d *schema.ResourceData, data map[string]interface{}) {
	updateStateAndStatus(d, data)
}

func updateStateAndStatus(d *schema.ResourceData, data map[string]interface{}) {
	if data["state"] != nil {
		clusterState := data["state"].(string)
		d.Set("state", clusterState)
	}

	if data["status"] != nil {
		clusterStatus := data["status"].([]interface{})
		d.Set("status", clusterStatus)
	}
}

func resourceHostGroupDirectConnectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(client.Client)
	id := clientID(d, client.ServiceConfig, "HostGroup")
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

func resourceHostGroupDirectConnectUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceHostGroupDirectConnectDelete(d *schema.ResourceData, m interface{}) error {
	return delete(d, m, client.ServiceConfig, "HostGroup", nil)
}
