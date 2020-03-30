package nirmata

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/nirmata/go-client/pkg/client"
)

func resourceClusterGke() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterGkeCreate,
		Read:   resourceClusterGkeRead,
		Update: resourceClusterGkeUpdate,
		Delete: resourceClusterGkeDelete,
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
			"disk_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					_ = v.(string)
					//TODO : ADD Logic
					return
				},
			},
			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					_ = v.(int)
					//TODO : ADD Logic
					return
				},
			},

			"cloud_provider_flag": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					_ = v.(string)
					//TODO : ADD Logic
					return
				},
			},
			"kubernetes_version": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					_ = v.(string)
					//TODO : ADD Logic
					return
				},
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceClusterGkeCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	name := d.Get("name").(string)
	diskSize := d.Get("disk_size")
	nodeType := d.Get("node_type").(string)
	nodeCount := d.Get("node_count").(int)
	region := d.Get("region").(string)
	kubernetesVersion := d.Get("kubernetes_version").(string)
	flagCloudProvider := d.Get("cloud_provider_flag").(string)

	cpID, err := getCloudProviderID(apiClient, "GoogleCloudPlatform", flagCloudProvider)

	if err != nil {
		fmt.Println(err.Error())
	}
	hg := map[string]interface{}{
		"name":         name,
		"orchestrator": "kubernetes",
		"mode":         "providerManaged",
		"upstreamType": "git",
		"kubernetesCluster": map[string]interface{}{
			"modelIndex": "KubernetesCluster",
			"clusterConfig": map[string]interface{}{
				"modelIndex":       "K8sClusterConfig",
				"cloudProviderRef": cpID.Map(),
				"version":          &kubernetesVersion,
				"nodeCount":        &nodeCount,
				"cloudProvider":    "GoogleCloudPlatform",
				"providerK8sClusterConfig": map[string]interface{}{
					"modelIndex":  "ProviderK8sClusterConfig",
					"region":      &region,
					"diskSize":    &diskSize,
					"machineType": &nodeType,
				},
			},
		},
	}

	data, err := apiClient.PostFromJSON(client.ServiceClusters, "hostClusters", hg, nil)
	if err != nil {
		return err
	}

	hgID := data["id"].(string)
	d.SetId(hgID)
	return nil
}

func resourceClusterGkeRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceClusterGkeUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceClusterGkeDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	name := d.Get("name").(string)

	id, err := apiClient.QueryByName(client.ServiceClusters, "hostClusters", name)
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

func getCloudProviderID(api client.Client, cpType string, flagCloudProvider string) (client.ID, error) {
	if flagCloudProvider != "" {
		return api.QueryByName(client.ServiceConfig, "cloudProviders", flagCloudProvider)
	}
	opts := client.NewGetModelID(nil, client.NewQuery().FieldEqualsValue("type", cpType))
	cIDs, err := api.GetCollection(client.ServiceConfig, "CloudProviders", opts)
	if err != nil {
		return nil, err
	}

	if len(cIDs) > 1 {
		names := make([]string, len(cIDs))
		for i, c := range cIDs {
			names[i] = c["name"].(string)
		}

		return nil, fmt.Errorf("Flag --cloud-provider <name> is required.\nAvailable Cloud Providers: %s", strings.Join(names, ", "))
	}

	cpObj, err := client.NewObject(cIDs[0])
	if err != nil {
		return nil, err
	}

	return cpObj.ID(), nil
}
