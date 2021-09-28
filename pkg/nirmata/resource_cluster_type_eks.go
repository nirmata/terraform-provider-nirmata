package nirmata

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/nirmata/go-client/pkg/client"
)

func resourceEksClusterType() *schema.Resource {
	return &schema.Resource{
		Create: resourceEksClusterTypeCreate,
		Read:   resourceEksClusterTypeRead,
		Update: resourceEksClusterTypeUpdate,
		Delete: resourceEksClusterTypeDelete,
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
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"cluster_role_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_groups": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"log_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"enable_private_endpoint": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_secrets_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"kms_key_arn": {
				Type:     schema.TypeString,
				Optional: true, // required if enable_secrets_encryption = true
			},
			"enable_identity_provider": {
				Type:     schema.TypeBool,
				Optional: true,
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
			"addons": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: addonSchema,
				},
			},
			"vault_auth": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: vaultAuthSchema,
				},
			},
			"nodepools": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: eksNodePoolSchema,
				},
			},
			"nodepool_field_override": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_fargate": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"pod_execution_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnets": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"namespace_label_selectors": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"pod_label_selectors": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"auto_sync_namespaces": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

var eksNodePoolSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"instance_type": {
		Type:     schema.TypeString,
		Required: true,
	},
	"disk_size": {
		Type:         schema.TypeInt,
		Required:     true,
		ValidateFunc: validateEKSDiskSize,
	},
	"security_groups": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Required: true,
	},
	"iam_role": {
		Type:     schema.TypeString,
		Required: true,
	},
	"ssh_key_name": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"ami_type": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"image_id": {
		Type:     schema.TypeString,
		Optional: true,
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

func resourceEksClusterTypeCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	credentials := d.Get("credentials").(string)
	cloudCredID, err := apiClient.QueryByName(client.ServiceClusters, "CloudCredentials", credentials)
	if err != nil {
		log.Printf("Error - %v", err)
		return err
	}

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	region := d.Get("region").(string)
	securityGroups := d.Get("security_groups")
	clusterRoleArn := d.Get("cluster_role_arn").(string)
	vpcID := d.Get("vpc_id").(string)
	subnetID := d.Get("subnet_id")
	logTypes := d.Get("log_types")
	privateEndpointAccess := d.Get("enable_private_endpoint")
	enableSecretsEncryption := d.Get("enable_secrets_encryption")
	keyArn := d.Get("kms_key_arn")
	enableIdentityProvider := d.Get("enable_identity_provider")
	systemMetadata := d.Get("system_metadata")
	autoSyncNamespaces := d.Get("auto_sync_namespaces").(bool)

	// Cluster override fields
	allowOverrideCredentials := d.Get("allow_override_credentials").(bool)
	clusterFieldOverride := d.Get("cluster_field_override")
	nodepoolFieldOverride := d.Get("nodepool_field_override")

	// Enable Fragment Fields
	enableFargate := d.Get("enable_fargate").(bool)
	podExecutionRoleArn := d.Get("pod_execution_role_arn").(string)
	subnets := d.Get("subnet_id")
	namespaceLabelSelectors := d.Get("namespace_label_selectors")
	podLabelSelectors := d.Get("pod_label_selectors")

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
					"eksConfig": map[string]interface{}{
						"instanceType":   element["instance_type"],
						"diskSize":       element["disk_size"],
						"securityGroups": element["security_groups"],
						"nodeIamRole":    element["iam_role"],
						"keyName":        element["ssh_key_name"],
						"amiType":        element["ami_type"],
						"imageId":        element["image_id"],
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
			"cloud":              "aws",
			"systemMetadata":     systemMetadata,
			"autoSyncNamespaces": autoSyncNamespaces,
			"addons":             addons,
			"cloudConfigSpec": map[string]interface{}{
				"modelIndex":               "CloudConfigSpec",
				"credentials":              credential,
				"allowOverrideCredentials": allowOverrideCredentials,
				"fieldsToOverride":         fieldsToOverride,
				"eksConfig": map[string]interface{}{
					"region":                  region,
					"vpcId":                   vpcID,
					"subnetId":                subnetID,
					"clusterRoleArn":          clusterRoleArn,
					"securityGroups":          securityGroups,
					"logTypes":                logTypes,
					"enableFargate":           enableFargate,
					"privateEndpointAccess":   privateEndpointAccess,
					"enableIdentityProvider":  enableIdentityProvider,
					"enableSecretsEncryption": enableSecretsEncryption,
					"keyArn":                  keyArn,
				},
				"nodePoolTypes": nodeobjArr,
			},
		},
	}

	if enableFargate {
		fargate := map[string]interface{}{
			"podExecutionRoleArn":     podExecutionRoleArn,
			"subnets":                 subnets,
			"namespaceLabelSelectors": namespaceLabelSelectors,
			"podLabelSelectors":       podLabelSelectors,
		}
		clusterTypeData["spec"].(map[string]interface{})["cloudConfigSpec"].(map[string]interface{})["eksConfig"].(map[string]interface{})["fargateSettings"] = fargate
	}

	if _, ok := d.GetOk("vault_auth"); ok {
		vl := d.Get("vault_auth").([]interface{})
		vault := vl[0].(map[string]interface{})
		vaultAuth, vErr := vaultAuthSchemaToVaultAuthSpec(vault, apiClient)
		if vErr != nil {
			log.Printf("Vault Credential Name not found")
			return fmt.Errorf("vault credential name not found : %v", vErr)
		}
		clusterTypeData["spec"].(map[string]interface{})["vault"] = vaultAuth
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

func resourceEksClusterTypeRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

var eksAttributeMap = map[string]string{
	"version":                    "version",
	"region":                     "region",
	"vpc_id":                     "vpcId",
	"subnet_id":                  "subnetId",
	"cluster_role_arn":           "clusterRoleArn",
	"security_groups":            "securityGroups",
	"log_types":                  "logTypes",
	"enable_private_endpoint":    "privateEndpointAccess",
	"enable_secrets_encryption":  "enableSecretsEncryption",
	"kms_key_arn":                "keyArn",
	"enable_identity_provider":   "enableIdentityProvider",
	"system_metadata":            "systemMetadata",
	"allow_override_credentials": "allowOverrideCredentials",
	"enable_fargate":             "enableFargate",
	"pod_execution_role_arn":     "podExecutionRoleArn",
	"subnets":                    "subnets",
	"namespace_label_selectors":  "namespaceLabelSelectors",
	"pod_label_selectors":        "podLabelSelectors",
	"auto_sync_namespaces":       "autoSyncNamespaces",
	"nodepools":                  "nodepools",
}

func resourceEksClusterTypeUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	clusterTypeID := client.NewID(client.ServiceClusters, "ClusterType", d.Id())

	// update ClusterSpec
	clusterSpecChanges := buildChanges(d, eksAttributeMap, "version", "auto_sync_namespaces", "system_metadata")
	if len(clusterSpecChanges) > 0 {
		err := updateDescendant(apiClient, clusterTypeID, "ClusterSpec", clusterSpecChanges)
		if err != nil {
			return err
		}
	}
	if err := updateClusterTypeAddonAndVault(d, meta); err != nil {
		log.Printf("[ERROR] - failed to update cluster type add-on and vault auth with data : %v", err)
		return err
	}

	// update EksClusterConfig
	eksConfigChanges := buildChanges(d, eksAttributeMap,
		"region",
		"vpc_id",
		"subnet_id",
		"cluster_role_arn",
		"system_metadata",
		"auto_sync_namespaces",
		"allow_override_credentials",
		"security_groups",
		"log_types",
		"enable_private_endpoint",
		"enable_secrets_encryption",
		"kms_key_arn",
		"enable_identity_provider",
		"enable_fargate",
		"pod_execution_role_arn",
		"subnets",
		"namespace_label_selectors",
		"pod_label_selectors",
		"auto_sync_namespaces")

	if len(eksConfigChanges) > 0 {
		err := updateDescendant(apiClient, clusterTypeID, "EksClusterConfig", eksConfigChanges)
		if err != nil {
			return err
		}
	}
	nodePoolChanges := buildChanges(d, eksAttributeMap, "nodepools")

	if len(nodePoolChanges) > 0 {
		cloudConfigSpecData, err := apiClient.GetDescendant(clusterTypeID, "CloudConfigSpec", &client.GetOptions{})
		if err != nil {
			log.Printf("[ERROR] - failed to retrieve %s from %v: %v", "CloudConfigSpec", clusterTypeID.Map(), err)
			return err
		}
		parent := map[string]interface{}{
			"id":         cloudConfigSpecData["id"],
			"service":    "Cluster",
			"modelIndex": "CloudConfigSpec",
		}
		nodepools := d.Get("nodepools").([]interface{})
		var createdNodepool []map[string]interface{}
		for i, nodepool := range nodepools {
			element, ok := nodepool.(map[string]interface{})
			if ok {
				createdNodepool = append(createdNodepool, map[string]interface{}{
					"modelIndex": "NodePoolType",
					"name":       d.Get("name").(string) + "-node-pool-" + strconv.Itoa(i),
					"spec": map[string]interface{}{
						"modelIndex":      "NodePoolSpec",
						"nodeLabels":      element["node_labels"],
						"nodeAnnotations": element["node_annotations"],
						"eksConfig": map[string]interface{}{
							"instanceType":   element["instance_type"],
							"diskSize":       element["disk_size"],
							"securityGroups": element["security_groups"],
							"nodeIamRole":    element["iam_role"],
							"keyName":        element["ssh_key_name"],
							"amiType":        element["ami_type"],
							"imageId":        element["image_id"],
						},
					},
					"parent": parent,
				})
			}
		}
		txn := make(map[string]interface{})
		txn["delete"] = cloudConfigSpecData["nodePoolTypes"]
		txn["create"] = createdNodepool
		_, txnErr := apiClient.PostFromJSON(client.ServiceClusters, "txn", txn, nil)
		if txnErr != nil {
			log.Printf("[ERROR] - failed to update cluster type nodeool with data : %v", txnErr)
			return txnErr
		}
	}

	return nil
}

func resourceEksClusterTypeDelete(d *schema.ResourceData, meta interface{}) error {
	return deleteObj(d, meta, client.ServiceClusters, "ClusterType")
}
