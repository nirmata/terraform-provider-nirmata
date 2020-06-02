package nirmata

import (
	"fmt"
	"regexp"

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

	hg := map[string]interface{}{
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

	data, err := apiClient.PostFromJSON(client.ServiceClusters, "kubernetesCluster", hg, nil)
	if err != nil {
		return err
	}

	hgID := data["id"].(string)
	d.SetId(hgID)
	return nil
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {

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
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("Deleted cluster %s", name)

	return nil
}
