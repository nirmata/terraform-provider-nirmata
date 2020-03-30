package nirmata

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceHostGroupDirectConnectRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceHostGroupDirectConnectUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

<<<<<<< HEAD
func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(client.Client)
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {

		uuid := d.Id()
		id := client.NewID(client.ServiceConfig, "HostGroup", uuid)
		if err := apiClient.Delete(id, nil); err != nil {
			if !strings.Contains(err.Error(), "404") {
				return resource.RetryableError(err)
			}
		}

		d.SetId("")
		return nil
	})
	if err != nil {
		return err
	}
	return nil
=======
func resourceHostGroupDirectConnectDelete(d *schema.ResourceData, m interface{}) error {
	return delete(d, m, client.ServiceConfig, "HostGroup")
>>>>>>> f32a3ccac2ad49ea846bbddd3d538f9fff732d56
}
