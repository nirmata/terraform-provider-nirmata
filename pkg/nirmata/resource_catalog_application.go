package nirmata

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resourceCatalogApplication() *schema.Resource {
	return &schema.Resource{

		Create: resourceCatalogApplicationCreate,
		Read:   resourceCatalogApplicationRead,
		Update: resourceCatalogApplicationUpdate,
		Delete: resourceCatalogApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateName,
			},
			"catalog": {
				Type:     schema.TypeString,
				Required: true,
			},
			"yamls": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"release_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCatalogApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	catalog := d.Get("catalog").(string)
	yamls := d.Get("yamls").(string)
	catID, cerr := apiClient.QueryByName(client.ServiceCatalogs, "Catalogs", catalog)
	if cerr != nil {
		log.Printf("[ERROR] - failed to find catalog with name : %v", name)
		return cerr
	}
	uuid := ""
	path := "createApplication"
	params := map[string]string{"name": name}
	results, perr := apiClient.PostWithID(catID, path, []byte(yamls), params)
	if results != nil {
		log.Printf("Received response %v", results["changes"])
	}

	if perr != nil {
		log.Printf("Error - failed to create application: %s", perr)
		return perr
	}
	fields := []string{"version", "name"}
	changes := results["changes"].(map[string]interface{})
	modifiedIds := changes["modifiedIds"].([]interface{})
	for _, product := range modifiedIds {
		id := product.(map[string]interface{})
		uuid = id["uuid"].(string)
		log.Printf("[INFO] -application %s %s", "id", id["uuid"])
	}

	d.SetId(uuid)
	log.Printf("[INFO] - created application %s %s", name, uuid)
	appID := client.NewID(client.ServiceCatalogs, "Applications", uuid)
	version, versionErr := apiClient.GetDescendant(appID, "Version", &client.GetOptions{Fields: fields})
	if versionErr != nil {
		log.Printf("Error version not found - %v", versionErr)
		return versionErr
	}
	d.Set("version", version["version"].(string))
	d.Set("release_name", version["name"].(string))

	return nil
}

func resourceCatalogApplicationRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceCatalogs, "applications", d.Id())

	_, err := apiClient.Get(id, &client.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] failed to retrieve application detail %s (%s): %v", name, id, err)
		return err
	}

	log.Printf("[INFO] - retrieved application %s %s", name, id.UUID())
	return nil
}

func resourceCatalogApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCatalogApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceCatalogs, "Application", d.Id())

	if err := apiClient.Delete(id, nil); err != nil {
		return err
	}

	log.Printf("[INFO] - deleted application %s %s", name, id.UUID())
	return nil
}
