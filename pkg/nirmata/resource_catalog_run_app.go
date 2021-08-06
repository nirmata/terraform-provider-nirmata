package nirmata

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resourceRunApplication() *schema.Resource {
	return &schema.Resource{

		Create: resourceRunApplicationCreate,
		Read:   resourceRunApplicationRead,
		Update: resourceRunApplicationUpdate,
		Delete: resourceRunApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
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
			"environments": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"clusters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceRunApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	version := d.Get("version").(string)
	channel := d.Get("channel").(string)
	application := d.Get("application").(string)
	catalog := d.Get("catalog").(string)
	environments := d.Get("environments").([]interface{})
	clusters := d.Get("clusters").([]interface{})

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
	fields := []string{"name", "id"}
	applicationList, appErr := apiClient.GetDescendants(catalogID, "Application", &client.GetOptions{Fields: fields})
	if appErr != nil {
		log.Printf("Error - %v", appErr)
		return err
	}

	var applicationID client.ID
	var channelID string
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

	for _, c := range channelList {
		if channel == c["name"] {
			channelID = c["id"].(string)
		}
	}
	channelMap := map[string]interface{}{
		"id":         channelID,
		"service":    "Catalog",
		"modelIndex": "Channel",
	}
	var catalogEnvArr = make([]interface{}, 0)
	var catalogClusterArr = make([]interface{}, 0)
	if len(environments) != 0 {
		for _, env := range environments {
			envId, _ := apiClient.QueryByName(client.ServiceEnvironments, "Environment", env.(string))
			obj := map[string]interface{}{
				"id":         envId.UUID(),
				"service":    "Environments",
				"modelIndex": "Environment",
			}
			catalogEnvArr = append(catalogEnvArr, obj)
		}
	}

	if len(clusters) != 0 {
		for _, env := range environments {
			cluId, _ := apiClient.QueryByName(client.ServiceClusters, "KubernetesCluster", env.(string))
			obj := map[string]interface{}{
				"id":         cluId.UUID(),
				"service":    "Cluster",
				"modelIndex": "KubernetesCluster",
			}
			catalogClusterArr = append(catalogClusterArr, obj)
		}
	}

	txnData := map[string]interface{}{
		"runName":    name,
		"name":       name + "-" + fmt.Sprint(rand.Int()),
		"parent":     versionID.UUID(),
		"envIds":     catalogEnvArr,
		"clusterIds": catalogClusterArr,
		"channel":    channelMap,
		"modelIndex": "Rollout",
	}

	rolloutId, txnErr := apiClient.PostFromJSON(client.ServiceCatalogs, "Rollout", txnData, nil)
	if txnErr != nil {
		log.Printf("[ERROR] - failed to create rollout with data : %v", txnErr)
		return txnErr
	}

	changeID := rolloutId["id"].(string)
	d.SetId(changeID)
	log.Printf("[INFO] - created  rollout %s %s", name, changeID)

	return nil
}

func resourceRunApplicationRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRunApplicationUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceRunApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceCatalogs, "Rollout", d.Id())

	if err := apiClient.Delete(id, nil); err != nil {
		return err
	}

	log.Printf("[INFO] - deleted rollout %s %s", name, id.UUID())
	return nil
}
