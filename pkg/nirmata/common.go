package nirmata

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func delete(d *schema.ResourceData, m interface{}, s client.Service, model string, params map[string]string) error {
	apiClient := m.(client.Client)
	id := clientID(d, s, model)
	if err := apiClient.Delete(id, params); err != nil {
		if !strings.Contains(err.Error(), "404") {
			return err
		}
	}

	d.SetId("")
	return nil
}

func clientID(d *schema.ResourceData, s client.Service, model string) client.ID {
	uuid := d.Id()
	return client.NewID(s, model, uuid)
}
