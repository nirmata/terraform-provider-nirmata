package nirmata

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resourcePromoteVersion() *schema.Resource {
	return &schema.Resource{

		Create: resourcePromoteVersionCreate,
		Read:   resourcePromoteVersionRead,
		Update: resourcePromoteVersionUpdate,
		Delete: resourcePromoteVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"rollout_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateName,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"channel": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application": {
				Type:     schema.TypeString,
				Required: true,
			},
			"catalog": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourcePromoteVersionCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	version := d.Get("version").(string)
	application := d.Get("application").(string)
	catalog := d.Get("catalog").(string)

	versionID, err := apiClient.QueryByName(client.ServiceCatalogs, "Version", version)
	if err != nil {
		fmt.Printf("Error - %v", err)
		return err
	}
	catalogID, appErr := apiClient.QueryByName(client.ServiceCatalogs, "Catalogs", catalog)
	if appErr != nil {
		log.Printf("Error  - %v", appErr)
		return err
	}
	fields := []string{"name", "id", "service", "modelIndex"}
	applicationList, appErr := apiClient.GetDescendants(catalogID, "Application", &client.GetOptions{Fields: fields})
	if appErr != nil {
		log.Printf("Error - %v", appErr)
		return err
	}

	var applicationID client.ID
	for _, apps := range applicationList {
		if application == apps["name"] {
			applicationID = client.NewID(client.ServiceCatalogs, "Application", apps["id"].(string))
		}
	}

	channelList, appErr := apiClient.GetDescendants(applicationID, "Channel", &client.GetOptions{Fields: fields})
	if appErr != nil {
		log.Printf("Error - %v", appErr)
		return err
	}
	parent := map[string]interface{}{
		"id":            applicationID.UUID(),
		"service":       "Catalog",
		"modelIndex":    "Application",
		"childRelation": "versions",
	}

	txnData := map[string]interface{}{
		"name":       name,
		"parent":     parent,
		"channel":    channelList,
		"modelIndex": "Version",
		"id":         versionID.UUID(),
	}

	rolloutId, txnErr := apiClient.PostFromJSON(client.ServiceCatalogs, "Rollout", txnData, nil)
	if txnErr != nil {
		log.Printf("[ERROR] - failed to promote version with data : %v", txnErr)
		return txnErr
	}

	changeID := rolloutId["changeId"].(string)
	d.SetId(changeID)
	log.Printf("[INFO] - version promoted %s %s", name, changeID)

	return nil
}

func resourcePromoteVersionRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePromoteVersionUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourcePromoteVersionDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
