package nirmata

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resoureClusterAddOn() *schema.Resource {
	return &schema.Resource{

		Create: resourceClusterAddOnCreate,
		Read:   resourceClusterAddOnRead,
		Update: resourceClusterAddOnUpdate,
		Delete: resourceClusterAddOnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateName,
			},
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"catalog": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"channel": {
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
func resourceClusterAddOnCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	clusterNameOrID := d.Get("cluster").(string)
	namespace := d.Get("namespace").(string)
	application := d.Get("application").(string)
	channel := d.Get("channel").(string)
	catalog := d.Get("catalog").(string)
	environment := d.Get("environment").(string)
	labels := d.Get("labels")
	isAppPresent := false

	clusterID, err := apiClient.QueryByName(client.ServiceClusters, "KubernetesCluster", clusterNameOrID)
	if err != nil {
		log.Printf("[ERROR] - failed to resolve cluster %s %v", clusterNameOrID, err)
		return err
	}

	catalogID, catalogErr := fetchID(apiClient, client.ServiceCatalogs, "Catalogs", catalog)
	if catalogErr != nil {
		log.Printf("[ERROR] - failed to get catalog %s %v", catalog, catalogErr)
		return catalogErr
	}

	fields := []string{"name", "id"}
	applicationList, appErr := apiClient.GetDescendants(catalogID, "Application", &client.GetOptions{Fields: fields})
	if appErr != nil {
		log.Printf("[ERROR] - failed to get application - %v", appErr)
		return err
	}
	for _, apps := range applicationList {
		if application == apps["name"] {
			isAppPresent = true
		}
	}

	if !isAppPresent {
		log.Printf("[ERROR] - no application found with name %v in catalog - %v", application, catalog)
		return fmt.Errorf("[ERROR] - no application found with name %v in catalog - %v", application, catalog)
	}

	addOnId, err := apiClient.GetRelationID(clusterID, "ClusterAddOns")
	if err != nil {
		log.Printf("[ERROR] - failed to resolve cluster %s %v", addOnId, err)
		return err
	}

	if namespace == "" {
		namespace = application
	}
	if environment == "" {
		if isUUID(clusterNameOrID) {
			clusterData, err := apiClient.Get(clusterID, nil)
			if err != nil {
				log.Printf("[ERROR] - failed to get cluster details %s %v", clusterID, err)
				return err
			}
			clusterNameOrID = application + "-" + clusterData["name"].(string)
		}
		environment = application + "-" + clusterNameOrID
	}

	data := map[string]interface{}{
		"name":        name,
		"parent":      addOnId.UUID(),
		"namespace":   namespace,
		"application": application,
		"channel":     channel,
		"environment": environment,
		"labels":      labels,
		"catalog":     catalog,
		"modelIndex":  "ClusterAddOn",
	}

	log.Printf("[DEBUG] - creating cluster add-on %s with %+v", name, data)
	addonData, err := apiClient.PostFromJSON(client.ServiceClusters, "ClusterAddOn", data, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create cluster add-on %s with data %v: %v", name, data, err)
		return err
	}

	addOnUUID := addonData["id"].(string)

	addonID := client.NewID(client.ServiceClusters, "ClusterAddOn", addOnUUID)
	state, waitErr := waitForAddonState(apiClient, d.Timeout(schema.TimeoutCreate), addonID)
	if waitErr != nil {
		log.Printf("[ERROR] - failed to check add-on status. Error - %v", waitErr)
		return waitErr
	}

	if strings.EqualFold("failed", state) {
		return fmt.Errorf("cluster add-on creation failed")
	}

	d.SetId(addOnUUID)
	log.Printf("[INFO] - created cluster addon %s %s", name, addOnUUID)

	return nil
}

func resourceClusterAddOnRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	cluster := d.Get("cluster").(string)
	_, err := fetchID(apiClient, client.ServiceClusters, "KubernetesCluster", cluster)
	if err != nil {
		log.Printf("[INFO] cluster  %+v not found", err)
		return err
	}
	d.SetId("")

	return nil
}

var addOnMap = map[string]string{
	"labels": "labels",
}

func resourceClusterAddOnUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	id := client.NewID(client.ServiceClusters, "clusteraddon", d.Id())
	addOnChanges := buildChanges(d, addOnMap, "labels")
	if len(addOnChanges) > 0 {
		addOnData, err := apiClient.Get(id, &client.GetOptions{})
		if err != nil {
			log.Printf("[ERROR] - failed to retrieve %s from %v: %v", "cluster add-on", id.Map(), err)
			return err
		}
		d, plainErr := client.NewObject(addOnData)
		if plainErr != nil {
			log.Printf("[ERROR] - failed to decode %s %v: %v", "cluster ad-don", d, plainErr)
			return err
		}
		_, plainErr = apiClient.PutWithIDFromJSON(d.ID(), addOnChanges)
		if plainErr != nil {
			log.Printf("[ERROR] - failed to update %s %v: %v", "cluster", d.ID().Map(), plainErr)
			return err
		}
		log.Printf("[DEBUG] updated %v %v", d.ID().Map(), addOnChanges)
	}
	return nil
}

func resourceClusterAddOnDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceClusters, "ClusterAddOn", d.Id())

	if err := apiClient.Delete(id, nil); err != nil {
		return err
	}

	log.Printf("[INFO] - deleted add-on %s %s", name, id.UUID())
	return nil
}
