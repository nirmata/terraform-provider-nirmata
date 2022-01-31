package nirmata

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/nirmata/go-client/pkg/client"
)

func resourceManagedCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterUpdate,
		Delete: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateName,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"override_credentials": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_field_override": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"nodepool_field_override": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"delete_action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "delete",
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"nodepools": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: clusterNodePoolSchema,
				},
			},
			"creation_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

var clusterNodePoolSchema = map[string]*schema.Schema{
	"node_count": {
		Type:         schema.TypeInt,
		Required:     true,
		ValidateFunc: validateNodeCount,
	},
	"min_count": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		ValidateFunc: validateNodeCount,
	},
	"max_count": {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validateNodeCount,
	},
	"enable_auto_scaling": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
}

func resourceClusterCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	labels := d.Get("labels")
	nodepools := d.Get("nodepools").([]interface{})
	typeSelector := d.Get("cluster_type").(string)
	credentials := d.Get("override_credentials").(string)
	timeout := d.Get("creation_timeout").(int)
	systemMetadata := d.Get("system_metadata")
	clusterFieldOverride := d.Get("cluster_field_override")
	nodepoolFieldOverride := d.Get("nodepool_field_override")

	spec, _, nodepool, err := getClusterTypeSpec(apiClient, typeSelector, nodepools)
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return err
	}

	mode := spec["clusterMode"]
	fieldsToOverride := map[string]interface{}{
		"cluster":  clusterFieldOverride,
		"nodePool": nodepoolFieldOverride,
	}

	data := map[string]interface{}{
		"name":                name,
		"mode":                mode,
		"typeSelector":        typeSelector,
		"credentialsSelector": credentials,
		"labels":              labels,
	}

	data["config"] = map[string]interface{}{
		"modelIndex":     "ClusterConfig",
		"version":        spec["version"],
		"nodeCount":      nil,
		"cloudProvider":  spec["cloud"],
		"systemMetadata": systemMetadata,
		"overrideValues": fieldsToOverride,
	}

	data["nodePools"] = nodepool

	clusterData, err := apiClient.PostFromJSON(client.ServiceClusters, "kubernetesCluster", data, nil)
	if err != nil {
		return err
	}

	clusterUUID := clusterData["id"].(string)
	d.SetId(clusterUUID)

	clusterID := client.NewID(client.ServiceClusters, "KubernetesCluster", clusterUUID)
	var cluster_timeout = d.Timeout(schema.TimeoutCreate)
	if timeout != 0 {
		cluster_timeout = time.Duration(timeout) * time.Minute
	}
	state, waitErr := waitForClusterState(apiClient, cluster_timeout, clusterID)
	if waitErr != nil {
		log.Printf("[ERROR] - failed to check cluster status. Error - %v", waitErr)
		return waitErr
	}

	if strings.EqualFold("failed", state) {
		status, err := getClusterStatus(apiClient, clusterID)
		if err != nil {
			log.Printf("[ERROR] - failed to retrieve cluster failure details: %v", err)
			return fmt.Errorf("cluster creation failed")
		}
		return fmt.Errorf("cluster creation failed: %s", status)
	}

	log.Printf("[INFO] created cluster %s with ID %s", name, clusterID)
	return nil
}

func getClusterTypeSpec(api client.Client, typeSelector string, nodepools []interface{}) (map[string]interface{}, []map[string]interface{}, []interface{}, error) {
	typeID, err := api.QueryByName(client.ServiceClusters, "ClusterType", typeSelector)
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return nil, nil, nil, err
	}
	var nodePoolObjArr = make([]interface{}, 0)
	nodepoolTypes, _ := api.GetDescendants(typeID, "NodePoolType", nil)
	cloudConfigSpec, _ := api.GetDescendants(typeID, "CloudConfigSpec", nil)

	for key, nodepool := range nodepools {
		element, ok := nodepool.(map[string]interface{})
		if ok {
			maxCount := element["node_count"].(int)
			minCount := 1
			if element["enable_auto_scaling"] == true {
				minCount = element["min_count"].(int)
				maxCount = element["max_count"].(int)
			}

			nodePoolObj := map[string]interface{}{
				"modelIndex":        "NodePool",
				"name":              "node-pool-" + strconv.Itoa(key),
				"minCount":          minCount,
				"maxCount":          maxCount,
				"nodeCount":         element["node_count"],
				"enableAutoScaling": element["enable_auto_scaling"],
				"typeSelector":      nodepoolTypes[key]["name"],
			}
			nodePoolObjArr = append(nodePoolObjArr, nodePoolObj)
		}
	}

	spec, err := api.GetRelation(typeID, "clusterSpecs")
	if err != nil {
		fmt.Println(err)
		return nil, nil, nil, err
	}
	return spec, cloudConfigSpec, nodePoolObjArr, nil
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	clusterID := client.NewID(client.ServiceClusters, "KubernetesCluster", d.Id())

	data, err := apiClient.Get(clusterID, &client.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] cluster %+v not found", clusterID.Map())
			d.SetId("")
			return nil
		}

		log.Printf("[ERROR] failed to retrieve cluster details %s (%s): %v", name, clusterID, err)
		return err
	}

	d.Set("state", data["state"])

	nodePools := data["nodePools"].([]interface{})
	if len(nodePools) == 0 {
		log.Printf("[INFO] failed to find nodepool for cluster %s (%s)", name, clusterID)
	} else {
		setErr := d.Set("nodepools", nodePools)
		if setErr != nil {
			log.Printf("[ERROR] failed to set nodepools: %v", nodePools)
		}
	}

	return nil
}

func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	if err := updateNodeCount(d, apiClient); err != nil {
		return err
	}
	labelsChanges := buildChanges(d, clustervaultMap, "labels")
	if len(labelsChanges) > 0 {
		if err := updateClusterLabels(d, apiClient); err != nil {
			log.Printf("[ERROR] - failed to update labels with data : %v", err)
			return err
		}
	}
	return nil
}

func updateNodeCount(d *schema.ResourceData, apiClient client.Client) error {
	var nodeCount int
	name := d.Get("name").(string)

	if d.HasChanges("node_count") {
		_, newNodeCount := d.GetChange("node_count")
		nodeCount = newNodeCount.(int)
	}

	if nodeCount == 0 {
		return nil
	}

	clusterID := client.NewID(client.ServiceClusters, "KubernetesCluster", d.Id())
	data, err := apiClient.Get(clusterID, &client.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] failed to retrieve cluster details %s (%s): %v", name, clusterID, err)
		return err
	}

	nodePools := data["nodePools"].([]interface{})
	if len(nodePools) == 0 {
		return fmt.Errorf("failed to find nodepool for cluster %s (%s)", name, clusterID)
	}

	if len(nodePools) > 1 {
		log.Printf("[INFO] found %d nodepools for cluster %s (%s)", len(nodePools), name, clusterID)
	}

	nodePool := nodePools[0]
	np := nodePool.(map[string]interface{})
	jsonObj := map[string]int{
		"nodeCount": nodeCount,
	}

	jsonString, jsonErr := json.Marshal(jsonObj)
	if jsonErr != nil {
		return fmt.Errorf("failed to marshall %v to JSON: %v", jsonObj, err)
	}

	restRequest := &client.RESTRequest{
		Service:     client.ServiceClusters,
		ContentType: "application/json",
		Path:        fmt.Sprintf("/NodePool/%s", np["id"].(string)),
		Data:        jsonString,
	}

	if _, err := apiClient.Put(restRequest); err != nil {
		return fmt.Errorf("failed to marshall %v to JSON: %v", jsonObj, err)
	}

	log.Printf("[INFO] Updated node count to %d for nodepool %s in cluster %s", nodeCount, np["name"], name)
	return nil
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	name := d.Get("name").(string)
	clusterID := client.NewID(client.ServiceClusters, "KubernetesCluster", d.Id())

	action := "delete"
	if d.Get("delete_action") != nil {
		action = d.Get("delete_action").(string)
	}

	params := map[string]string{
		"action": action,
	}

	if err := apiClient.Delete(clusterID, params); err != nil {
		return err
	}

	waitForDeletedState(apiClient, d.Timeout(schema.TimeoutCreate), clusterID)

	log.Printf("[INFO] Deleted cluster %s", name)
	return nil
}
