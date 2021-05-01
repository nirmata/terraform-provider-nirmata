package nirmata

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resoureClusterAddon() *schema.Resource {
	return &schema.Resource{

		Create: resourceClusterAddonCreate,
		Read:   resourceClusterAddonRead,
		Update: resourceClusterAddonUpdate,
		Delete: resourceClusterAddonDelete,
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
				Required: true,
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
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_scheme": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
func resourceClusterAddonCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	clusterNameOrID := d.Get("cluster").(string)
	namespace := d.Get("namespace").(string)
	application := d.Get("application").(string)
	channel := d.Get("channel").(string)
	catalog := d.Get("catalog").(string)
	environment := d.Get("environment").(string)
	labels := d.Get("labels")

	if namespace == "" {
		namespace = application
	}
	if environment == "" {
		environment = application + "-" + clusterNameOrID
	}
	clusterID, err := apiClient.QueryByName(client.ServiceClusters, "KubernetesCluster", clusterNameOrID)
	if err != nil {
		log.Printf("[ERROR] - failed to resolve cluster %s %v", clusterNameOrID, err)
		return err
	}
	clusterID1, err := apiClient.GetRelationID(clusterID, "ClusterAddOns")
	if err != nil {
		log.Printf("[ERROR] - failed to resolve cluster %s %v", clusterID1, err)
		return err
	}
	data := map[string]interface{}{
		"name":        name,
		"parent":      clusterID1.UUID(),
		"namespace":   namespace,
		"application": application,
		"channel":     channel,
		"environment": environment,
		"labels":      labels,
		"catalog":     catalog,
		"modelIndex":  "ClusterAddOn",
	}

	log.Printf("[DEBUG] - creating cluster addon %s with %+v", name, data)
	addonData, err := apiClient.PostFromJSON(client.ServiceClusters, "ClusterAddOn", data, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create cluster addon %s with data %v: %v", name, data, err)
		return err
	}

	addonUUID := addonData["id"].(string)
	d.SetId(addonUUID)
	log.Printf("[INFO] - created cluster addon %s %s", name, addonUUID)

	return nil
}

func resourceClusterAddonRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

var addonMap = map[string]string{
	"service_name":   "serviceName",
	"service_port":   "servicePort",
	"service_scheme": "serviceScheme",
	"labels":         "labels",
}

func resourceClusterAddonUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	id := client.NewID(client.ServiceClusters, "clusteraddon", d.Id())
	addonChanges := buildChanges(d, addonMap, "service_name", "service_port", "service_scheme", "labels")
	if len(addonChanges) > 0 {
		addonData, err := apiClient.Get(id, &client.GetOptions{})
		if err != nil {
			log.Printf("[ERROR] - failed to retrieve %s from %v: %v", "cluster addon", id.Map(), err)
			return err
		}
		d, plainErr := client.NewObject(addonData)
		if plainErr != nil {
			log.Printf("[ERROR] - failed to decode %s %v: %v", "cluster addon", d, err)
			return err
		}
		_, plainErr = apiClient.PutWithIDFromJSON(d.ID(), addonData)
		if plainErr != nil {
			log.Printf("[ERROR] - failed to update %s %v: %v", "cluster", d.ID().Map(), err)
			return err
		}
		log.Printf("[DEBUG] updated %v %v", d.ID().Map(), addonData)
	}
	return nil
}

func resourceClusterAddonDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceClusters, "ClusterAddOn", d.Id())

	if err := apiClient.Delete(id, nil); err != nil {
		return err
	}

	log.Printf("[INFO] - deleted addon %s %s", name, id.UUID())
	return nil
}
