package nirmata

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

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

	catalogID, catErr := apiClient.QueryByName(client.ServiceCatalogs, "Catalogs", catalog)
	if catErr != nil {
		log.Printf("Error: catalog not found - %v", catErr)
		return catErr
	}
	fields := []string{"name", "id"}
	applicationList, appErr := apiClient.GetDescendants(catalogID, "Application", &client.GetOptions{Fields: fields})
	if appErr != nil {
		log.Printf("Error application not found - %v", appErr)
		return appErr
	}

	var applicationID client.ID
	var channelID string
	var versionID string
	for _, apps := range applicationList {
		if application == apps["name"] {
			applicationID = client.NewID(client.ServiceCatalogs, "Application", apps["id"].(string))
		}
	}

	channelList, channelErr := apiClient.GetDescendants(applicationID, "Channel", &client.GetOptions{Fields: fields})
	if appErr != nil {
		log.Printf("Error channel not found - %v", appErr)
		return channelErr
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

	versionList, versionErr := apiClient.GetDescendants(applicationID, "Version", &client.GetOptions{Fields: fields})
	if appErr != nil {
		log.Printf("Error version not found - %v", appErr)
		return versionErr
	}

	for _, c := range versionList {
		if version == c["name"] {
			versionID = c["id"].(string)
		}
	}

	var catalogEnvArr = make([]interface{}, 0)
	var catalogClusterArr = make([]interface{}, 0)
	if len(environments) != 0 {
		for _, env := range environments {
			envId, envErr := apiClient.QueryByName(client.ServiceEnvironments, "Environment", env.(string))
			if envErr != nil {
				log.Printf("Error environment not found - %v", envErr)
				return envErr
			}
			obj := map[string]interface{}{
				"id":         envId.UUID(),
				"service":    "Environments",
				"modelIndex": "Environment",
			}
			catalogEnvArr = append(catalogEnvArr, obj)
		}
	}

	if len(clusters) != 0 {
		for _, cluster := range clusters {
			cluId, cErr := apiClient.QueryByName(client.ServiceClusters, "KubernetesCluster", cluster.(string))
			if cErr != nil {
				log.Printf("Error cluster not found - %v", cErr)
				return cErr
			}
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
		"parent":     versionID,
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

	rolloutID := client.NewID(client.ServiceCatalogs, "Rollout", changeID)
	state, waitErr := waitForRollutState(apiClient, d.Timeout(schema.TimeoutCreate), rolloutID)
	if waitErr != nil {
		log.Printf("[ERROR] - failed to check rollout status. Error - %v", waitErr)
		return waitErr
	}

	if strings.EqualFold("failed", state) {
		status, err := getRolloutStatus(apiClient, rolloutID)
		if err != nil {
			log.Printf("[ERROR] - failed to retrieve rollout failure details: %v", err)
			return fmt.Errorf("rollout creation failed")
		}
		return fmt.Errorf("rollout creation failed: %s", status)
	}

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
	environments := d.Get("environments").([]interface{})
	application := d.Get("application").(string)
	fields := []string{"name", "id"}
	if len(environments) != 0 {
		for _, env := range environments {
			envId, envErr := apiClient.QueryByName(client.ServiceEnvironments, "Environment", env.(string))
			if envErr != nil {
				log.Printf("Error environment not found - %v", envErr)
				return envErr
			}
			applicationList, appErr := apiClient.GetDescendants(envId, "Application", &client.GetOptions{Fields: fields})
			if appErr != nil {

				log.Printf("Error applications not found - %v", appErr)
				return appErr
			}
			for _, apps := range applicationList {
				if application == apps["name"] {
					applicationID := client.NewID(client.ServiceEnvironments, "Application", apps["id"].(string))
					id := client.NewID(client.ServiceEnvironments, "Application", apps["id"].(string))
					if err := apiClient.Delete(id, nil); err != nil {
						return err
					}
					waitForDeletedState(apiClient, d.Timeout(schema.TimeoutCreate), applicationID)
				}
			}
		}
	}

	log.Printf("[INFO] - deleted rollout %s", name)
	return nil
}
