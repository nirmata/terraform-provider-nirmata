package nirmata

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"strconv"
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
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) > 64 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 64 characters", k))
					}
					if !regexp.MustCompile(`^[\w+=,.@-]*$`).MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must match [\\w+=,.@-]", k))
					}
					return
				},
			},
			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if v.(int) > 999 {
						errors = append(errors, fmt.Errorf(
							"%q The node count must be between 1 and 1000", k))
					}
					return
				},
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
				Optional:  true,
				Elem: &schema.Schema{
				  Type: schema.TypeString,
				},
			  },
			"cluster_field_override": {
				Type:     schema.TypeMap,
				Optional:  true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceClusterCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	nodeCount := d.Get("node_count").(int)
	typeSelector := d.Get("cluster_type").(string)
	credentials := d.Get("override_credentials").(string)
	systemMetadata := d.Get("system_metadata")
	clusterFieldOverride := d.Get("cluster_field_override")

	spec,_,nodepool, err := getClusterTypeSpec(apiClient, typeSelector)
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return err
	}
	mode := spec["clusterMode"]
	fieldsToOverride :=  map[string]interface{} {
		"cluster" : clusterFieldOverride,
	}
	data := map[string]interface{}{
		"name":         name,
		"mode":         mode,
		"typeSelector": typeSelector,
		"credentialsSelector": credentials,
	}
	data["config"] = map[string]interface{}{
			"modelIndex":    "ClusterConfig",
			"version":       spec["version"],
			"nodeCount":     nodeCount,
			"cloudProvider": spec["cloud"],
			"systemMetadata": systemMetadata,
			"overrideValues": fieldsToOverride,
		}
	data["nodePools"] = nodepool;

	clusterData, err := apiClient.PostFromJSON(client.ServiceClusters, "kubernetesCluster", data, nil)
	if err != nil {
		return err
	}

	clusterUUID := clusterData["id"].(string)
	d.SetId(clusterUUID)

	clusterID := client.NewID(client.ServiceClusters, "KubernetesCluster", clusterUUID)

	state, waitErr := waitForClusterState(apiClient, d.Timeout(schema.TimeoutCreate), clusterID)
	if waitErr != nil {
		log.Printf("[ERROR] - failed to check cluster status. Error - %v", waitErr)
		return nil
	}

	if strings.EqualFold("failed", state) {
		status, err := getClusterStatus(apiClient, clusterID)
		if err != nil {
			log.Printf("[ERROR] - failed to retrieve cluster failure details: %v", err)
			return fmt.Errorf("cluster creation failed")
		}

		return fmt.Errorf("cluster creation failed: %s", status)
	}

	log.Printf("created cluster %s with ID %s", name, clusterID)
	return nil
}
func getClusterTypeSpec(api client.Client, typeSelector string ) (map[string]interface{}, []map[string]interface {}, []interface{},error) {
	typeID, err := api.QueryByName(client.ServiceClusters, "ClusterType", typeSelector)
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return nil,nil,nil, err;
	}
	var nodePoolObjArr = make([]interface{}, 0);
	nodepoolTypes, _ := api.GetDescendants(typeID, "NodePoolType", nil);
	cloudConfigSpec, _ := api.GetDescendants(typeID, "CloudConfigSpec", nil);
	for key, nodepoolType := range nodepoolTypes {
		nodePoolObj := map[string]interface{}{
			"modelIndex":      "NodePool",
			"name":            "node-pool-"+ strconv.Itoa(key),
			"minCount": 		1,
          	"maxCount": 		1,
			"nodeCount": 		1,
			"typeSelector":     nodepoolType["name"],
		}
		nodePoolObjArr = append(nodePoolObjArr, nodePoolObj);
	}
	spec, err := api.GetRelation(typeID, "clusterSpecs")
	if err != nil {
		fmt.Println(err)
		return nil,nil, nil,err;
	}
	return spec,cloudConfigSpec,nodePoolObjArr, nil;
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)

	clusterID, err := apiClient.QueryByName(client.ServiceClusters, "KubernetesCluster", name)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] cluster does not exist %s (%s): %v", name, err)
			d.SetId("")
			return nil
		}

		return err
	}

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
	d.Set("nodeCount", np["nodeCount"])

	return nil
}

func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	var nodeCount int
	name := d.Get("name").(string)

	if d.HasChanges("node_count") {
		_, NewNodeCount := d.GetChange("node_count")
		nodeCount = NewNodeCount.(int)
	}
	clusterID, err := apiClient.QueryByName(client.ServiceClusters, "KubernetesCluster", name)
	if err != nil {
		log.Printf("[ERROR] failed to find cluster %s: %v", name, err)
		return err
	}

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

	id, err := apiClient.QueryByName(client.ServiceClusters, "kubernetesCluster", name)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] cluster does not exist %s (%s): %v", name, err)
			d.SetId("")
			return nil
		}

		log.Printf("[ERROR] - %v", err)
		return err
	}

	params := map[string]string{
		"action": "delete",
	}

	if err := apiClient.Delete(id, params); err != nil {
		return err
	}

	log.Printf("[INFO] Deleted cluster %s", name)
	return nil
}
