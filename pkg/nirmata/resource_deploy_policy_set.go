package nirmata

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resourceDeployPolicySet() *schema.Resource {
	return &schema.Resource{

		Create: resourceDeployPolicySetCreate,
		Read:   resourceDeployPolicySetRead,
		Update: resourceDeployPolicySetUpdate,
		Delete: resourceDeployPolicySetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"policy_grop_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateName,
			},
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delete_from_cluster": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceDeployPolicySetCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	nameOrID := d.Get("policy_grop_name").(string)
	clusterNameOrID := d.Get("cluster").(string)

	clusterID, err := fetchID(apiClient, client.ServiceClusters, "KubernetesCluster", clusterNameOrID)
	if err != nil {
		log.Printf("[ERROR] - failed to resolve cluster %s %v", clusterNameOrID, err)
		return err
	}

	groupID, err := fetchID(apiClient, client.ServiceClusters, "PolicyGroup", nameOrID)
	if err != nil {
		log.Printf("[ERROR] - failed to resolve policy group %s %v", clusterNameOrID, err)
		return err
	}

	clusterRef := map[string]interface{}{
		"id":         clusterID.UUID(),
		"service":    "Cluster",
		"modelIndex": "KubernetesCluster",
	}

	data := map[string]interface{}{
		"name":       clusterNameOrID,
		"parent":     groupID.UUID(),
		"clusterRef": clusterRef,
	}

	log.Printf("[DEBUG] - deploying policy group %s with %+v", nameOrID, data)
	deployData, err := apiClient.PostFromJSON(client.ServiceClusters, "PolicyGroupCluster", data, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to deploy policy group %s with data %v: %v", nameOrID, data, err)
		return err
	}

	deployDataUUID := deployData["id"].(string)
	d.SetId(deployDataUUID)

	policyGroupClusterID := client.NewID(client.ServiceClusters, "PolicyGroupCluster", deployData["id"].(string))

	state, waitErr := waitForPolicyDeploySetSyncStatus(apiClient, d.Timeout(schema.TimeoutCreate), policyGroupClusterID)
	if waitErr != nil {
		log.Printf("[ERROR] - failed to deploy policy group sync status. Error - %v", waitErr)
		return waitErr
	}

	if strings.EqualFold("failed", state) {
		status, err := getPolicyDeployGroupStatus(apiClient, policyGroupClusterID)
		if err != nil {
			log.Printf("[ERROR] - failed to retrieve policy group sync details: %v", err)
			return fmt.Errorf(" [ERROR] - policy group sync failed")
		}
		return fmt.Errorf(" [ERROR] - policy group sync failed: %s", status)
	}

	log.Printf("[INFO] - policy group deploy successfully %s %s", nameOrID, deployDataUUID)

	return nil
}

func resourceDeployPolicySetRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("cluster").(string)
	id := client.NewID(client.ServiceClusters, "PolicyGroupCluster", d.Id())

	_, err := apiClient.Get(id, &client.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] deployed policy group %+v not found", id.Map())
			d.SetId("")
			return nil
		}

		log.Printf("[ERROR] failed to retrieve deployed policy group details %s (%s): %v", name, id, err)
		return err
	}

	log.Printf("[INFO] - retrieved policy group %s %s", name, id.UUID())
	return nil
}

func resourceDeployPolicySetUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceDeployPolicySetDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("cluster").(string)
	deleteFromCluster := d.Get("delete_from_cluster").(bool)
	id := client.NewID(client.ServiceClusters, "PolicyGroupCluster", d.Id())

	params := map[string]string{
		"deleteFromCluster": strconv.FormatBool(deleteFromCluster),
	}

	if err := apiClient.Delete(id, params); err != nil {
		return err
	}

	log.Printf("[INFO] - deleted policy group %s %s", name, id.UUID())
	return nil
}
