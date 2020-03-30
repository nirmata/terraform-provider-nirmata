package nirmata

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func delete(d *schema.ResourceData, m interface{}, s client.Service, model string) error {
	apiClient := m.(client.Client)

	uuid := d.Id()
	id := client.NewID(s, model, uuid)
	if err := apiClient.Delete(id, nil); err != nil {
		if !strings.Contains(err.Error(), "404") {
			return err
		}
	}

	d.SetId("")
	return nil
}
