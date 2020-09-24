package nirmata

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/nirmata/go-client/pkg/client"
)

func resourceProviderManagedCluster() *schema.Resource {
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
			"type_selector": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceClusterCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	name := d.Get("name").(string)
	nodeCount := d.Get("node_count").(int)
	typeSelector := d.Get("type_selector").(string)

	clusterTypeID, err := apiClient.QueryByName(client.ServiceClusters, "ClusterType", typeSelector)
	if err != nil {
		fmt.Printf("Error - %v", err)
		return err
	}
	cspec, err := apiClient.GetRelation(clusterTypeID, "clusterSpecs")
	if err != nil {
		fmt.Printf("Error - %v", err)
		return err
	}

	pmc := map[string]interface{}{
		"name":         name,
		"mode":         "providerManaged",
		"typeSelector": typeSelector,
		"config": map[string]interface{}{
			"modelIndex":    "ClusterConfig",
			"version":       cspec["version"],
			"nodeCount":     nodeCount,
			"cloudProvider": cspec["cloud"],
		},
	}

	data, err := apiClient.PostFromJSON(client.ServiceClusters, "kubernetesCluster", pmc, nil)
	if err != nil {
		return err
	}

	pmcID := data["id"].(string)
	d.SetId(pmcID)
	return nil
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	name := d.Get("name").(string)

	clusterID, err := apiClient.QueryByName(client.ServiceClusters, "clustertypes", name)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	data, err := apiClient.Get(clusterID, &client.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		fmt.Printf("failed to retrieve cluster details %s (%s): %v", name, clusterID, err)
		return err
	}

	nodePools := data["nodePools"].([]interface{})
	if len(nodePools) == 0 {
		return fmt.Errorf("failed to find nodepool for cluster %s (%s)", name, clusterID)
	}

	if len(nodePools) > 1 {
		fmt.Printf("found %d nodepools for cluster %s (%s)", len(nodePools), name, clusterID)
	}

	nodePool := nodePools[0]
	np := nodePool.(map[string]interface{})
	d.Set("nodeCount",np["nodeCount"])
	return nil
}

func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	var nodeCount int
	name := d.Get("name").(string)

	if d.HasChanges("node_count") {
		_,NewNodeCount := d.GetChange("node_count")
		nodeCount = NewNodeCount.(int)
	}
	clusterID, err := apiClient.QueryByName(client.ServiceClusters, "KubernetesCluster", name)
	if err != nil {
		fmt.Printf("failed to find cluster %s: %v", name, err)
		return err
	}

	data, err := apiClient.Get(clusterID, &client.GetOptions{})
	if err != nil {
		fmt.Printf("failed to retrieve cluster details %s (%s): %v", name, clusterID, err)
		return err
	}

	nodePools := data["nodePools"].([]interface{})
	if len(nodePools) == 0 {
		return fmt.Errorf("failed to find nodepool for cluster %s (%s)", name, clusterID)
	}

	if len(nodePools) > 1 {
		fmt.Printf("found %d nodepools for cluster %s (%s)", len(nodePools), name, clusterID)
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

	fmt.Printf("Updated node count to %d for nodepool %s in cluster %s", nodeCount, np["name"], name)
	return nil
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	name := d.Get("name").(string)

	id, err := apiClient.QueryByName(client.ServiceClusters, "kubernetesCluster", name)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	params := map[string]string{
		"action": "delete",
	}

	if err := apiClient.Delete(id, params); err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("Deleted cluster %s", name)

	return nil
}
