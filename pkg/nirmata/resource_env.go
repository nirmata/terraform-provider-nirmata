package nirmata

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{

		Create: resourceEnvironmentCreate,
		Read:   resourceEnvironmentRead,
		Update: resourceEnvironmentUpdate,
		Delete: resourceEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateName,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment_update_action": {
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

func resourceEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	envType := d.Get("type").(string)
	clusterNameOrID := d.Get("cluster").(string)
	namespace := d.Get("namespace").(string)
	labels := d.Get("labels")
	policy := d.Get("environment_update_action").(string)

	if policy == "" || policy != "update" {
		policy = "notify"
	}

	clusterID, err := fetchID(apiClient, client.ServiceClusters, "KubernetesCluster", clusterNameOrID)
	if err != nil {
		log.Printf("[ERROR] - failed to resolve cluster %s %v", envType, err)
		return err
	}

	data := map[string]interface{}{
		"name":         name,
		"resourceType": envType,
		"hostCluster":  clusterID.Map(),
		"namespace":    namespace,
		"labels":       labels,
		"updatePolicy": map[string]interface{}{
			"configUpdateAction": policy,
		},
	}

	log.Printf("[DEBUG] - creating environment %s with %+v", name, data)
	envData, err := apiClient.PostFromJSON(client.ServiceEnvironments, "Environment", data, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create environment %s with data %v: %v", name, data, err)
		return err
	}

	envUUID := envData["id"].(string)
	d.SetId(envUUID)
	log.Printf("[INFO] - created environment %s %s", name, envUUID)

	return nil
}

func resourceEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceEnvironments, "Environment", d.Id())

	_, err := apiClient.Get(id, &client.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] cluster %+v not found", id.Map())
			d.SetId("")
			return nil
		}

		log.Printf("[ERROR] failed to retrieve environment details %s (%s): %v", name, id, err)
		return err
	}

	log.Printf("[INFO] - retrieved environment %s %s", name, id.UUID())
	return nil
}

func resourceEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceEnvironments, "Environment", d.Id())

	if err := apiClient.Delete(id, nil); err != nil {
		return err
	}

	log.Printf("[INFO] - deleted environment %s %s", name, id.UUID())
	return nil
}
