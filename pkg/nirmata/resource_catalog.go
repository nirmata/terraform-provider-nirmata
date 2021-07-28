package nirmata

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resourceCatalog() *schema.Resource {
	return &schema.Resource{

		Create: resourceCatalogCreate,
		Read:   resourceCatalogRead,
		Update: resourceCatalogUpdate,
		Delete: resourceCatalogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateName,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceCatalogCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	labels := d.Get("labels")

	data := map[string]interface{}{
		"name":        name,
		"description": description,
		"labels":      labels,
	}

	log.Printf("[DEBUG] - creating catalog %s with %+v", name, data)
	catalogData, err := apiClient.PostFromJSON(client.ServiceCatalogs, "Catalogs", data, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create catalog %s with data %v: %v", name, data, err)
		return err
	}
	catalogDataUUID := catalogData["id"].(string)
	d.SetId(catalogDataUUID)
	log.Printf("[INFO] - created catalog %s %s", name, catalogDataUUID)

	return nil
}

func resourceCatalogRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceCatalogs, "Catalogs", d.Id())

	_, err := apiClient.Get(id, &client.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] failed to retrieve catalog details %s (%s): %v", name, id, err)
		return err
	}

	log.Printf("[INFO] - retrieved catalog %s %s", name, id.UUID())
	return nil
}

func resourceCatalogUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCatalogDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceCatalogs, "Catalogs", d.Id())

	if err := apiClient.Delete(id, nil); err != nil {
		return err
	}

	log.Printf("[INFO] - deleted catalog %s %s", name, id.UUID())
	return nil
}
