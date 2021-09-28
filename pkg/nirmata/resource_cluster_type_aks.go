package nirmata

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"strings"

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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateName,
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
			"enable_private_cluster": {
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAksFields,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"system_metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_override_credentials": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"cluster_field_override": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"network_profile": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_service_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"docker_bridge_cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_sync_namespaces": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

var aksNodePoolSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"subnet_id": {
		Type:     schema.TypeString,
		Required: true,
	},
	"vms_size": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validateAksFields,
	},
	"vm_set_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validateAksFields,
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
	"network": {
		Type:     schema.TypeString,
		Required: false,
	},
	"os_type": {
		Type:     schema.TypeString,
		Required: false,
	},
	"node_annotations": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"node_labels": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
}

func resourceClusterTypeCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	credentials := d.Get("credentials").(string)
	region := d.Get("region").(string)
	resourceGroup := d.Get("resource_group").(string)
	workspaceID := d.Get("workspaceid").(string)
	httpsApplicationRouting := d.Get("httpsapplicationrouting").(bool)
	monitoring := d.Get("monitoring").(bool)
	systemMetadata := d.Get("system_metadata")
	enablePrivateCluster := d.Get("enable_private_cluster").(bool)
	networkProfile := d.Get("network_profile").(string)
	networkPolicy := d.Get("network_policy").(string)
	serviceCidr := d.Get("service_cidr").(string)
	dnsServiceIp := d.Get("dns_service_ip").(string)
	dockerBridgeCidr := d.Get("docker_bridge_cidr").(string)
	autoSyncNamespaces := d.Get("auto_sync_namespaces").(bool)

	// Cluster override fields
	allowOverrideCredentials := d.Get("allow_override_credentials").(bool)
	clusterFieldOverride := d.Get("cluster_field_override")
	nodepoolFieldOverride := d.Get("nodepool_field_override")

	cloudCredID, err := apiClient.QueryByName(client.ServiceClusters, "CloudCredentials", credentials)
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return err
	}

	fieldsToOverride := map[string]interface{}{
		"cluster":  clusterFieldOverride,
		"nodePool": nodepoolFieldOverride,
	}

	var nodeobjArr = make([]interface{}, 0)
	nodepools := d.Get("nodepools").([]interface{})
	for i, node := range nodepools {
		element, ok := node.(map[string]interface{})
		if ok {
			nodePoolObj := map[string]interface{}{
				"modelIndex": "NodePoolType",
				"name":       name + "-node-pool-" + strconv.Itoa(i),
				"spec": map[string]interface{}{
					"modelIndex":      "NodePoolSpec",
					"nodeLabels":      element["node_labels"],
					"nodeAnnotations": element["node_annotations"],
					"aksConfig": map[string]interface{}{
						"subnetId":  element["subnet_id"],
						"vmSize":    element["vmsize"],
						"vmSetType": element["vm_set_type"],
						"diskSize":  element["disk_size"],
						"network":   element["network"],
						"osType":    element["os_type"],
					},
				},
			}

			nodeobjArr = append(nodeobjArr, nodePoolObj)
		}
	}

	addons := addOnsSchemaToAddOns(d)
	credential := map[string]interface{}{
		"id":         cloudCredID.UUID(),
		"service":    "Cluster",
		"modelIndex": "CloudCredentials",
	}

	clusterTypeData := map[string]interface{}{
		"name":        name,
		"description": "",
		"modelIndex":  "ClusterType",
		"spec": map[string]interface{}{
			"clusterMode":        "providerManaged",
			"modelIndex":         "ClusterSpec",
			"version":            version,
			"cloud":              "azure",
			"systemMetadata":     systemMetadata,
			"autoSyncNamespaces": autoSyncNamespaces,
			"addons":             addons,
			"cloudConfigSpec": map[string]interface{}{
				"modelIndex":               "CloudConfigSpec",
				"credentials":              credential,
				"allowOverrideCredentials": allowOverrideCredentials,
				"fieldsToOverride":         fieldsToOverride,
				"aksConfig": map[string]interface{}{
					"region":                  region,
					"resourceGroup":           resourceGroup,
					"httpsApplicationRouting": httpsApplicationRouting,
					"enablePrivateCluster":    enablePrivateCluster,
					"monitoring":              monitoring,
					"workspaceId":             workspaceID,
					"modelIndex":              "AksClusterConfig",
					"networkProfile":          networkProfile,
					"serviceCidr":             serviceCidr,
					"dnsServiceIp":            dnsServiceIp,
					"dockerBridgeCidr":        dockerBridgeCidr,
					"networkPolicy":           networkPolicy,
					"networkPlugin":           "kubenet",
					"podCidr":                 "10.244.0.0/16",
				},
				"nodePoolTypes": nodeobjArr,
			},
		},
	}

	txn := make(map[string]interface{})
	var objArr = make([]interface{}, 0)
	objArr = append(objArr, clusterTypeData)
	txn["create"] = objArr
	data, err := apiClient.PostFromJSON(client.ServiceClusters, "txn", txn, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create cluster type  with data : %v", err)
		return err
	}

	obj, resultErr := extractCreateFromTxnResult(data, "ClusterType")
	if resultErr != nil {
		log.Printf("[ERROR] - %v", err)
		return resultErr
	}

	d.SetId(obj.ID().UUID())
	return nil
}

func resourceClusterTypeRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceClusterTypeUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := updateClusterTypeAddonAndVault(d, meta); err != nil {
		log.Printf("[ERROR] - failed to update cluster type add-on and vault auth with data : %v", err)
		return err
	}
	return nil
}

func resourceClusterTypeDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	name := d.Get("name").(string)

	id, err := apiClient.QueryByName(client.ServiceClusters, "clustertypes", name)
	if err != nil {
		log.Printf("ERROR - %v", err)
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

		log.Printf("[ERROR] - %v", err)
		return err
	}

	log.Printf("[INFO] Deleted cluster type %s", name)
	return nil
}
