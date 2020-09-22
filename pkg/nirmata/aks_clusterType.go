package nirmata

import (
	"fmt"
	"regexp"
	"time"

	guuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/nirmata/go-client/pkg/client"
)

func resourceAksClusterType() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterTypeCreate,
		Read:   resourceClusterTypeRead,
		Update: resourceClusterTypeUpdate,
		Delete: resourceClusterTypeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
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
			"https_application_routing": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"monitoring": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group": {
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

			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vms_ize": {
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
			"vm_set_type": {
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
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if v.(int) < 29 {
						errors = append(errors, fmt.Errorf(
							"%q The disk size must be grater than 29", k))
					}
					return
				},
			},
		},
	}
}

func resourceClusterTypeCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	clouduuid := guuid.New()
	nodepooluuid := guuid.New()

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	credentials := d.Get("credentials").(string)
	region := d.Get("region").(string)
	resourceGroup := d.Get("resource_group").(string)
	subnetID := d.Get("subnet_id").(string)
	vmSize := d.Get("vmsize").(string)
	vmSetType := d.Get("vmsettype").(string)
	workspaceID := d.Get("workspaceid").(string)
	httpsApplicationRouting := d.Get("httpsapplicationrouting").(bool)
	monitoring := d.Get("monitoring").(bool)
	diskSize := d.Get("disksize").(int)

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
			"cloud":       "azure",
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
				"aksConfig": map[string]interface{}{
					"region":                  region,
					"resourceGroup":           resourceGroup,
					"httpsApplicationRouting": httpsApplicationRouting,
					"monitoring":              monitoring,
					"workspaceId":             workspaceID,
					"modelIndex":              "AksClusterConfig",
					"networkProfile":          "basic",
					"serviceCidr":             "10.0.0.0/16",
					"dnsServiceIp":            "10.0.0.10",
					"dockerBridgeCidr":        "172.17.0.1/16",
					"networkPolicy":           "",
					"networkPlugin":           "kubenet",
					"podCidr":                 "10.244.0.0/16",
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
			"aksConfig": map[string]interface{}{
				"subnetId":   subnetID,
				"vmSize":     vmSize,
				"vmSetType":  vmSetType,
				"diskSize":   diskSize,
				"osType":     "Linux",
				"modelIndex": "AksNodePoolConfig",
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

	waitErr  := waitForState(apiClient,d.Timeout(schema.TimeoutCreate),name)
	if waitErr != nil {
		return waitErr
	}

	pmcID := data["id"].(string)
	d.SetId(pmcID)
	d.Set("curl_script", pmcID)
	return nil
}

func resourceClusterTypeRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceClusterTypeUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceClusterTypeDelete(d *schema.ResourceData, meta interface{}) error {
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
