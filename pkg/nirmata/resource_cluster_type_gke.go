package nirmata

import (
	"fmt"
	"log"
	"strings"
	"time"
	"strconv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/nirmata/go-client/pkg/client"
)

var gkeSchema = map[string]*schema.Schema{
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
		Optional: true,
	},
	"zone": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"location_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validateGKELocationType,
	},
	"node_locations": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional: true,
	},
	"enable_secrets_encryption": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"enable_workload_identity": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
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
	"start_time": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"duration": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "10",
	},
	"cluster_ipv4_cidr": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"services_ipv4_cidr": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"cloud_run": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"allow_override_credentials": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"enable_network_policy": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"http_load_balancing": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"enable_vertical_pod_autoscaling": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"horizontal_pod_autoscaling": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"enable_maintenance_policy": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"exclusion_timewindow": {
		Type:     schema.TypeMap,
		Optional:  true,
		Elem: &schema.Schema{
		  Type: schema.TypeString,
		},
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
	  "nodepool_field_override": {
		Type:     schema.TypeMap,
		Optional:  true,
		Elem: &schema.Schema{
		  Type: schema.TypeString,
		},
	  },
	  "nodepooltype": {
		Type:         schema.TypeSet,
		Required:     true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
			  "machinetype": {
				Type:      schema.TypeString,
        		Required:  true,
			  },
			  "disksize": {
				Type:      schema.TypeInt,
        		Required:  true,
			  },
			  "service_account": {
				Type:      schema.TypeString,
        		Optional:  true,
			  },
			  "enable_preemptible_nodes": {
				Type:      schema.TypeBool,
        		Optional:  true,
			  },
			  "node_annotations": {
				Type:     schema.TypeMap,
				Optional:  true,
				Elem: &schema.Schema{
				  Type: schema.TypeString,
				},
			  },
			  "node_labels": {
				Type:     schema.TypeMap,
				Optional:  true,
				Elem: &schema.Schema{
				  Type: schema.TypeString,
				},
			  },  
		  },
		},
	},
}

func resourceGkeClusterType() *schema.Resource {
	return &schema.Resource{
		Create: resourceGkeClusterTypeCreate,
		Read:   resourceGkeClusterTypeRead,
		Update: resourceGkeClusterTypeUpdate,
		Delete: resourceGkeClusterTypeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: gkeSchema,
	}
}

func resourceGkeClusterTypeCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	version := d.Get("version").(string)
	credentials := d.Get("credentials").(string)
	region := d.Get("region").(string)
	zone := d.Get("zone").(string)
	locationType := d.Get("location_type").(string)
	nodeLocations := d.Get("node_locations")
	enableSecretsEncryption := d.Get("enable_secrets_encryption").(bool)
	secretsEncryptionKey := d.Get("secrets_encryption_key").(string)
	enableWorkloadIdentity := d.Get("enable_workload_identity").(bool)
	workloadPool := d.Get("workload_pool").(string)
	network := d.Get("network").(string)
	subnetwork := d.Get("subnetwork").(string)
	nodepooltype := d.Get("nodepooltype").(*schema.Set).List()
	cloudRun := d.Get("cloud_run").(bool)
	allowOverrideCredentials := d.Get("allow_override_credentials").(bool)
	enableNetworkPolicy := d.Get("enable_network_policy").(bool)
	httpLoadBalancing := d.Get("http_load_balancing").(bool)
	enableVerticalPodAutoscaling := d.Get("enable_vertical_pod_autoscaling").(bool)
	horizontalPodAutoscaling := d.Get("horizontal_pod_autoscaling").(bool)
	enableMaintenancePolicy := d.Get("enable_maintenance_policy").(bool)
	clusterIpv4Cidr := d.Get("cluster_ipv4_cidr").(string)
	servicesIpv4Cidr := d.Get("services_ipv4_cidr").(string)
	duration := d.Get("duration").(string)
	startTime := d.Get("start_time").(string)
	systemMetadata := d.Get("system_metadata")
	exclusionTimeWindow := d.Get("exclusion_timewindow")
	clusterFieldOverride := d.Get("cluster_field_override")
	nodepoolFieldOverride := d.Get("nodepool_field_override")
	
	cloudCredID, err := apiClient.QueryByName(client.ServiceClusters, "CloudCredentials", credentials)
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return err
	}
	if locationType == "Regional" && region == "" {
		return fmt.Errorf("region is required when location_type is Regional")
	}

	if locationType == "Zonal" && zone == "" {
		return fmt.Errorf("zone is required when location_type is Zonal")
	}

	if enableSecretsEncryption && len(secretsEncryptionKey) == 0 {
		return fmt.Errorf("\nError - Encryption key is required if secrets encryption is enabled")
	}
	if enableWorkloadIdentity && len(workloadPool) == 0 {
		return fmt.Errorf("\nWorkload Pool is required if workload identity is enabled")
	}

	var addon []string
	if horizontalPodAutoscaling {
		addon = append(addon, "horizontalPodAutoscaling")
	}
	if httpLoadBalancing {
		addon = append(addon, "httpLoadBalancing")
	}
	if cloudRun {
		addon = append(addon, "cloudRunConfig")
	}
	fieldsToOverride :=  map[string]interface{} {
		"cluster" : clusterFieldOverride,
		"nodePool": nodepoolFieldOverride,
	}
	var nodeobjArr = make([]interface{}, 0)
	for key, node := range nodepooltype {
		element, ok := node.(map[string]interface{})
		if ok {
			nodePoolObj := map[string]interface{}{
				"modelIndex":      "NodePoolType",
				"name":            name + "-node-pool-"+ strconv.Itoa(key),
				"spec": map[string]interface{}{
					"modelIndex": "NodePoolSpec",
					"nodeLabels":  element["node_labels"],
					"nodeAnnotations":  element["node_annotations"],
					"gkeConfig": map[string]interface{}{
						"machineType": element["machinetype"],
						"diskSize":    element["disksize"],
						"serviceAccount":    element["service_account"],
						"enablePreemptibleNodes":    element["enable_preemptible_nodes"],
						"modelIndex":  "GkeNodePoolConfig",
					},
				},
			}
			nodeobjArr = append(nodeobjArr, nodePoolObj)
		}
	}
	clusterTypeData := map[string]interface{}{
		"name":        name,
		"description": "",
		"modelIndex":  "ClusterType",
		"spec": map[string]interface{}{
			"clusterMode": "providerManaged",
			"modelIndex":  "ClusterSpec",
			"version":     version,
			"cloud":       "googlecloudplatform",
			"systemMetadata": systemMetadata,
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
				"allowOverrideCredentials":    allowOverrideCredentials,
				"fieldsToOverride":    fieldsToOverride,
				"modelIndex":    "CloudConfigSpec",
				"gkeConfig": map[string]interface{}{
					"region":                  region,
					"zone":                    zone,
					"locationType":            locationType,
					"defaultNodeLocations":    nodeLocations,
					"enableSecretsEncryption": enableSecretsEncryption,
					"secretsEncryptionKey":    secretsEncryptionKey,
					"enableWorkloadIdentity":  enableWorkloadIdentity,
					"workloadPool":            workloadPool,
					"modelIndex":              "GkeClusterConfig",
					"network":                 network,
					"subnetwork":              subnetwork,
					"clusterIpv4Cidr": clusterIpv4Cidr,
                  	"servicesIpv4Cidr": servicesIpv4Cidr,
                  	"enableNetworkPolicy": enableNetworkPolicy,
                  	"enableMaintenancePolicy": enableMaintenancePolicy,
					"duration": duration,
					"startTime": startTime,
					"enableVerticalPodAutoscaling": enableVerticalPodAutoscaling,
					"exclusionTimeWindow": exclusionTimeWindow,
					"addons": addon,
				},
				"nodePoolTypes": nodeobjArr,
			},
		},
	}
	
	txn := make(map[string]interface{})
	var objArr = make([]interface{}, 0)
	objArr = append(objArr, clusterTypeData)
	txn["create"] = objArr
	data, err :=  apiClient.PostFromJSON(client.ServiceClusters, "txn", txn, nil)
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

var gkeClusterTypePaths = map[string]string{
	"version":                   	"spec[0].version",
	"region":                    	"spec[0].cloudConfigSpec[0].gkeConfig[0].region",
	"network":                   	"spec[0].cloudConfigSpec[0].gkeConfig[0].network",
	"subnetwork":                	"spec[0].cloudConfigSpec[0].gkeConfig[0].subnetwork",
	"zone":                      	"spec[0].cloudConfigSpec[0].gkeConfig[0].zone",
	"location_type":             	"spec[0].cloudConfigSpec[0].gkeConfig[0].locationType",
	"node_locations":            	"spec[0].cloudConfigSpec[0].gkeConfig[0].defaultNodeLocations",
	"enable_workload_identity":  	"spec[0].cloudConfigSpec[0].gkeConfig[0].enableWorkloadIdentity",
	"enable_secrets_encryption": 	"spec[0].cloudConfigSpec[0].gkeConfig[0].enableSecretsEncryption",
	"secrets_encryption_key":    	"spec[0].cloudConfigSpec[0].gkeConfig[0].secretsEncryptionKey",
	"workload_pool":             	"spec[0].cloudConfigSpec[0].gkeConfig[0].workloadPool",
	"start_time":             	  	"spec[0].cloudConfigSpec[0].gkeConfig[0].startTime",
	"exclusion_timewindow":         "spec[0].cloudConfigSpec[0].gkeConfig[0].exclusionTimeWindow",
	"enable_maintenance_policy":     "spec[0].cloudConfigSpec[0].gkeConfig[0].enableMaintenancePolicy",
	"cluster_ipv4_cidr":             "spec[0].cloudConfigSpec[0].gkeConfig[0].clusterIpv4Cidr",
	"services_ipv4_cidr":            "spec[0].cloudConfigSpec[0].gkeConfig[0].servicesIpv4Cidr",
	"enable_vertical_pod_autoscaling":"spec[0].cloudConfigSpec[0].gkeConfig[0].enableVerticalPodAutoscaling",
	"duration":             		 "spec[0].cloudConfigSpec[0].gkeConfig[0].duration",
	"enable_network_policy":         "spec[0].cloudConfigSpec[0].gkeConfig[0].enableNetworkPolicy",
}

var nodePoolTypePaths = map[string]string{
	"machinetype": 			"spec[0].gkeConfig[0].machineType",
	"disksize":    			"spec[0].gkeConfig[0].diskSize",
	"enable_preemptible_nodes": "spec[0].gkeConfig[0].enablePreemptibleNodes",
	"service_account":    		"spec[0].gkeConfig[0].serviceAccount",
	"node_labels":    			"spec[0].nodeLabels",
	"node_annotations":    		"spec[0].nodeAnnotations",
}

func resourceGkeClusterTypeRead(d *schema.ResourceData, meta interface{}) (err error) {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	clusterTypeID := client.NewID(client.ServiceClusters, "ClusterType", d.Id())

	clusterTypeData, err := apiClient.Get(clusterTypeID, &client.GetOptions{Mode: client.OutputModeExport})
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] cluster type %+v not found", clusterTypeID.Map())
			d.SetId("")
			return nil
		}

		log.Printf("[ERROR] - failed to retrieve cluster details %s (%s): %v", name, clusterTypeID.UUID(), err)
		return err
	}

	for field, path := range gkeClusterTypePaths {
		s := gkeSchema[field]
		err = updateData(field, d, s, path, clusterTypeData)
		if err != nil {
			return fmt.Errorf("failed to update field %s from %s: %v", field, path, err)
		}
	}

	return nil
}

var gkeAttributeMap = map[string]string{
	"version":                   	"version",
	"region":                    	"region",
	"network":                   	"network",
	"subnetwork":                	"subnetwork",
	"zone":                      	"zone",
	"location_type":             	"locationType",
	"node_locations":            	"defaultNodeLocations",
	"enable_workload_identity":  	"enableWorkloadIdentity",
	"enable_secrets_encryption": 	"enableSecretsEncryption",
	"secrets_encryption_key":    	"secretsEncryptionKey",
	"workload_pool":             	"workloadPool",
	"start_time":             	  	"startTime",
	"exclusion_timewindow":         "exclusionTimeWindow",
	"enable_maintenance_policy":     "enableMaintenancePolicy",
	"cluster_ipv4_cidr":             "clusterIpv4Cidr",
	"services_ipv4_cidr":            "servicesIpv4Cidr",
	"enable_vertical_pod_autoscaling":"enableVerticalPodAutoscaling",
	"duration":             		 "duration",
	"enable_network_policy":         "enableNetworkPolicy",
} 

var nodePoolAttributeMap = map[string]string{
	"machinetype": "machineType",
	"disksize":    "diskSize",
	"enable_preemptible_nodes": "enablePreemptibleNodes",
	"service_account":    		"serviceAccount",
	"node_labels":    			"nodeLabels",
	"node_annotations":    		"nodeAnnotations",
}

func resourceGkeClusterTypeUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	apiClient := meta.(client.Client)
	clusterTypeID := client.NewID(client.ServiceClusters, "ClusterType", d.Id())

	// update ClusterSpec
	clusterSpecChanges := buildChanges(d, gkeAttributeMap, "version")
	if len(clusterSpecChanges) > 0 {
		err := updateDescendant(apiClient, clusterTypeID, "ClusterSpec", clusterSpecChanges)
		if err != nil {
			return err
		}
	}

	// update GkeClusterConfig
	gkeConfigChanges := buildChanges(d, gkeAttributeMap, "region",
		"network", "subnetwork", "zone", "location_type", "node_locations",
		"enable_workload_identity", "enable_secrets_encryption", "secrets_encryption_key",
		"workload_pool", "start_time", "exclusion_timewindow", "enable_maintenance_policy" ,"cluster_ipv4_cidr",
		"services_ipv4_cidr", "enable_vertical_pod_autoscaling", "duration", "enable_network_policy")

	if len(gkeConfigChanges) > 0 {
		err := updateDescendant(apiClient, clusterTypeID, "GkeClusterConfig", gkeConfigChanges)
		if err != nil {
			return err
		}
	}

	// update NodePool
	nodePoolChanges := buildChanges(d, nodePoolAttributeMap, "machinetype", "disksize", "enable_preemptible_nodes","service_account","node_labels","node_annotations")
	if len(nodePoolChanges) > 0 {
		nodePoolData, err := getNodePoolType(apiClient, clusterTypeID)
		if err != nil {
			return err
		}

		npo, err := client.NewObject(nodePoolData)
		if err != nil {
			log.Printf("[ERROR] - failed to decode node pool %v: %v", nodePoolData, err)
			return err
		}

		err = updateDescendant(apiClient, npo.ID(), "GkeNodePoolConfig", nodePoolChanges)
		if err != nil {
			return err
		}
	}

	return nil
}

func buildChanges(d *schema.ResourceData, nameMap map[string]string, attributes ...string) map[string]interface{} {
	changes := map[string]interface{}{}
	for _, a := range attributes {
		if d.HasChange(a) {
			name := nameMap[a]
			changes[name] = d.Get(a)
		}
	}

	return changes
}

func updateDescendant(apiClient client.Client, id client.ID, descendant string, changes map[string]interface{}) error {
	clusterSpecData, err := apiClient.GetDescendant(id, descendant, &client.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] - failed to retrieve %s from %v: %v", descendant, id.Map(), err)
		return err
	}

	d, plainErr := client.NewObject(clusterSpecData)
	if plainErr != nil {
		log.Printf("[ERROR] - failed to decode %s %v: %v", descendant, d, err)
		return err
	}

	_, plainErr = apiClient.PutWithIDFromJSON(d.ID(), changes)
	if plainErr != nil {
		log.Printf("[ERROR] - failed to update %s %v: %v", descendant, d.ID().Map(), err)
		return err
	}

	log.Printf("[DEBUG] updated %v %v", d.ID().Map(), changes)
	return nil
}

func resourceGkeClusterTypeDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id, err := apiClient.QueryByName(client.ServiceClusters, "clustertypes", name)
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return err
	}

	params := map[string]string{
		"action": "delete",
	}

	if err := apiClient.Delete(id, params); err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] - %v not found: %v", id.Map(), err)
			d.SetId("")
			return nil
		}

		log.Printf("[ERROR] - %v", err)
		return err
	}

	log.Printf("[INFO] Deleted cluster type %s", name)
	return nil
}
