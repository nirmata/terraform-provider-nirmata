package nirmata

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resourceHelmApplication() *schema.Resource {
	return &schema.Resource{

		Create: resourceHelmApplicationCreate,
		Read:   resourceHelmApplicationRead,
		Update: resourceHelmApplicationUpdate,
		Delete: resourceHelmApplicationDelete,
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
				Optional: true,
			},
			"application": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"chart_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceHelmApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	catalog := d.Get("catalog").(string)
	application := d.Get("application").(string)
	repository := d.Get("repository").(string)
	appVersion := d.Get("app_version").(string)
	chartVersion := d.Get("chart_version").(string)
	location := d.Get("location").(string)

	repositoryData, cerr := apiClient.QueryByName(client.ServiceCatalogs, "ChartRepository", repository)
	if cerr != nil {
		log.Printf("[ERROR] - Failed to fetch controller YAML %s: %v \n", name, cerr)
		return cerr
	}
	fieldsToOverride := map[string]interface{}{
		"label": "chartVersion: " + chartVersion + ", appVersion: " + appVersion,
		"value": "chartVersion: " + chartVersion + ", appVersion: " + appVersion,
	}
	txn := map[string]interface{}{
		"location": location,
		"name":     application,
		"version":  fieldsToOverride,
	}

	data, err := apiClient.PostFromJSON(client.ServiceCatalogs, "ChartRepository/"+repositoryData.UUID()+"/getChart", txn, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create cluster type  with data : %v", err)
		return err
	}
	catID, cerr := apiClient.QueryByName(client.ServiceCatalogs, "Catalogs", catalog)
	if cerr != nil {
		log.Printf("[ERROR] - failed to find catalog with name : %v", name)
		return cerr
	}

	appData := map[string]interface{}{
		"name":         name,
		"parent":       catID.Map(),
		"upstreamType": "helm",
		"helmConfig": map[string]interface{}{
			"chartVersion":     chartVersion,
			"appVersion":       appVersion,
			"valueFileContent": data["values.yaml"],
			"chartName":        application,
			"chartRepo": map[string]interface{}{
				"id": repositoryData.UUID(),
			},
		},
	}

	helmData, marshalErr := json.Marshal(appData)
	if marshalErr != nil {
		fmt.Printf("Error - %v", err)
		return err
	}

	log.Printf("[DEBUG] - creating  application %s with %+v", name, appData)
	appId, err := apiClient.PostWithID(catID, "applications", helmData, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create  application %s with data %v: %v", name, appData, err)
		return err
	}
	catalogAppUUID := appId["id"].(string)
	d.SetId(catalogAppUUID)
	log.Printf("[INFO] - created application %s %s", name, catalogAppUUID)
	return nil
}

func resourceHelmApplicationRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceHelmApplicationUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceHelmApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceCatalogs, "Application", d.Id())

	if err := apiClient.Delete(id, nil); err != nil {
		return err
	}

	log.Printf("[INFO] - deleted application %s %s", name, id.UUID())
	return nil
}
