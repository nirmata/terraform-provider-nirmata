package nirmata

import (
	"fmt"
	"regexp"

	guuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/nirmata/go-client/pkg/client"
)

func resourceGkeClusterType() *schema.Resource {
	return &schema.Resource{
		Create: resourceGkeClusterTypeCreate,
		Read:   resourceGkeClusterTypeRead,
		Update: resourceGkeClusterTypeUpdate,
		Delete: resourceGkeClusterTypeDelete,
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
			"machine_type": {
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
			"disk_size": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if v.(int) < 9 {
						errors = append(errors, fmt.Errorf(
							"%q The disk size must be grater than 9", k))
					}
					return
				},
			},
			"location_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default : "Regional",
			},
			"node_locations": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"enable_secrets_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				Default : false,
			},
			"enable_workload_identity": {
				Type:     schema.TypeBool,
				Optional: true,
				Default : false,
			},
			"secrets_encryption_key": {
				Type:     schema.TypeString,
				Optional: true,	
			},
			"workload_pool": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnetwork": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGkeClusterTypeCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	clouduuid := guuid.New()
	nodepooluuid := guuid.New()

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	credentials := d.Get("credentials").(string)
	region := d.Get("region").(string)
	machinetype := d.Get("machine_type").(string)
	diskSize := d.Get("disk_size").(int)
	locationType := d.Get("location_type").(string)
	nodeLocations := d.Get("node_locations")
	enableSecretsEncryption := d.Get("enable_secrets_encryption").(bool)
	secretsEncryptionKey := d.Get("secrets_encryption_key").(string)
	enableWorkloadIdentity := d.Get("enable_workload_identity").(bool)
	workloadPool := d.Get("workload_pool").(string)
	network := d.Get("network").(string)
	subnetwork := d.Get("subnetwork").(string)

	cloudCredID, err := apiClient.QueryByName(client.ServiceClusters, "CloudCredentials", credentials)
	fmt.Printf("Error - %v", cloudCredID)
	if err != nil {
		fmt.Printf("Error - %v", err)
		return err
	}
	var areaType = "zone"
	if locationType == "Regional" {
		areaType = "region"
	}
	clustertype := map[string]interface{}{
		"name":        name,
		"description": "",
		"modelIndex":  "ClusterType",
		"spec": map[string]interface{}{
			"clusterMode": "providerManaged",
			"modelIndex":  "ClusterSpec",
			"version":     version,
			"cloud":       "googlecloudplatform",
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
				"gkeConfig": map[string]interface{}{
					areaType:                  region,
					"locationType":            locationType,
					"defaultNodeLocations":    nodeLocations,
					"enableSecretsEncryption": enableSecretsEncryption,
					"secretsEncryptionKey":    secretsEncryptionKey,
					"enableWorkloadIdentity":  enableWorkloadIdentity,
					"workloadPool":            workloadPool,
					"modelIndex":              "GkeClusterConfig",
					"network":                 network,
					"subnetwork":              subnetwork,
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
			"gkeConfig": map[string]interface{}{
				"machineType": machinetype,
				"diskSize":    diskSize,
				"modelIndex":  "GkeNodePoolConfig",
			},
		},
	}
	txn := make(map[string]interface{})
	var objArr = make([]interface{}, 0)
	objArr = append(objArr, clustertype, nodepoolobj)
	txn["create"] = objArr
	data, err := apiClient.PostFromJSON(client.ServiceClusters, "txn", txn, nil)
	if err != nil {
		fmt.Printf("\nError - failed to create cluster type  with data : %v", err)
		return err
	}
	changeID := data["changeId"].(string)
	d.SetId(changeID)
	return nil
}

func resourceGkeClusterTypeRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceGkeClusterTypeUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceGkeClusterTypeDelete(d *schema.ResourceData, meta interface{}) error {
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
