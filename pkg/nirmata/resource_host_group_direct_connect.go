package nirmata

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resourceHostGroupDirectConnect() *schema.Resource {

	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

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

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
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

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(client.Client)

	uuid := d.Id()
	id := client.NewID(client.ServiceConfig, "HostGroup", uuid)
	if err := apiClient.Delete(id, nil); err != nil {
		if !strings.Contains(err.Error(), "404") {
			return err
		}
	}

	d.SetId("")
	return nil
}
