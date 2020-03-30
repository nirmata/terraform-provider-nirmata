package nirmata

import (
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
		},
	}
}

func resourceHostGroupDirectConnectCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(client.Client)

	name := d.Get("name").(string)
	labels := d.Get("labels").(map[string]interface{})

	cpID, err := apiClient.QueryByName(client.ServiceConfig, "CloudProvider", "Direct Connect")
	if err != nil {
		return err
	}

	hg := map[string]interface{}{
		"name":   name,
		"parent": cpID.UUID(),
		"labels": labels,
	}

	data, err := apiClient.PostFromJSON(client.ServiceConfig, "HostGroup", hg, nil)
	if err != nil {
		return err
	}

	hgID := data["id"].(string)
	d.SetId(hgID)
	return nil
}

func resourceHostGroupDirectConnectRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceHostGroupDirectConnectUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceHostGroupDirectConnectDelete(d *schema.ResourceData, m interface{}) error {
	return delete(d, m, client.ServiceConfig, "HostGroup")
}
