package nirmata

import (
	"fmt"
	"regexp"

	guuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/nirmata/go-client/pkg/client"
)

func resourceOkeClusterType() *schema.Resource {
	return &schema.Resource{
		Create: resourceOkeClusterTypeCreate,
		Read:   resourceOkeClusterTypeRead,
		Update: resourceOkeClusterTypeUpdate,
		Delete: resourceOkeClusterTypeDelete,
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
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"credentials": {
				Type:     schema.TypeString,
				Required: true,
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vmshape": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !regexp.MustCompile(`^[\w+=,.@-]*$`).MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must match [\\w+=,.@-]", k))
					}
					return
				},
			},
		},
	}
}

func resourceOkeClusterTypeCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	clouduuid := guuid.New()
	nodepooluuid := guuid.New()

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	credentials := d.Get("credentials").(string)
	region := d.Get("region").(string)
	vmshape := d.Get("vmshape").(string)

	cloudCredID, err := apiClient.QueryByName(client.ServiceClusters, "CloudCredentials", credentials)
	fmt.Printf("Error - %v", cloudCredID)
	if err != nil {
		fmt.Printf("Error - %v", err)
		return err
	}

	clustertype := map[string]interface{}{
		"name":        name,
		"description": "",
		"modelIndex":  "ClusterType",
		"spec": map[string]interface{}{
			"clusterMode": "providerManaged",
			"modelIndex":  "ClusterSpec",
			"version":     version,
			"cloud":       "oraclecloudservices",
			"addons": map[string]interface{}{
				"dns":        false,
				"modelIndex": "AddOns",
				"addons": map[string]interface{}{
					"name":          "kyverno",
					"addOnSelector": "kyverno",
					"catalog":       "default-addon-catalog",
				},
			},
			"cloudConfigSpec": map[string]interface{}{
				"credentials":   cloudCredID.UUID(),
				"id":            clouduuid,
				"modelIndex":    "CloudConfigSpec",
				"nodePoolTypes": nodepooluuid,
				"okeConfig": map[string]interface{}{
					"region":     region,
					"modelIndex": "OkeClusterConfig",
				},
			},
		},
	}

	nodepoolobj := map[string]interface{}{
		"id":              nodepooluuid,
		"modelIndex":      "NodePoolType",
		"name":            name + "-default-node-pool-type",
		"cloudConfigSpec": clouduuid,
		"spec": map[string]interface{}{
			"modelIndex": "NodePoolSpec",
			"okeConfig": map[string]interface{}{
				"vmshape":    vmshape,
				"modelIndex": "OkeNodePoolConfig",
			},
		},
	}

	data, err := apiClient.PostFromJSON(client.ServiceClusters, "clustertypes", clustertype, nil)
	if err != nil {
		return err
	}

	_, nerr := apiClient.PostFromJSON(client.ServiceClusters, "nodepooltypes", nodepoolobj, nil)
	if nerr != nil {
		return err
	}

	pmcID := data["id"].(string)
	d.SetId(pmcID)
	return nil
}

func resourceOkeClusterTypeRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceOkeClusterTypeUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceOkeClusterTypeDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	name := d.Get("name").(string)

	id, err := apiClient.QueryByName(client.ServiceClusters, "clustertypes", name)
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

	fmt.Printf("Deleted cluster type %s", name)

	return nil
}
