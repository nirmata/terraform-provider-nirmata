package nirmata

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/nirmata/go-client/pkg/client"
)

var gkeClusterTypeSchema = map[string]*schema.Schema{
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
	"cluster_ipv4_cidr": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"services_ipv4_cidr": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"enable_cloud_run": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"enable_network_policy": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"enable_http_load_balancing": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"enable_vertical_pod_autoscaling": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"enable_horizontal_pod_autoscaling": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"enable_maintenance_policy": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"maintenance_start_time": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"maintenance_duration": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "10",
	},
	"maintenance_recurrence": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"maintenance_exclusion_timewindow": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
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
	"nodepool_field_override": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"nodepools": {
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: gkeNodePoolSchema,
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
}

var gkeNodePoolSchema = map[string]*schema.Schema{
	"machine_type": {
		Type:     schema.TypeString,
		Required: true,
	},
	"disk_size": {
		Type:     schema.TypeInt,
		Required: true,
	},
	"service_account": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"enable_preemptible_nodes": {
		Type:     schema.TypeBool,
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

var addonSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"addon_selector": {
		Type:     schema.TypeString,
		Required: true,
	},
	"catalog": {
		Type:     schema.TypeString,
		Required: true,
	},
	"channel": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"namespace": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"sequence_number": {
		Type:     schema.TypeInt,
		Required: true,
	},
}

var vaultAuthSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"path": {
		Type:     schema.TypeString,
		Required: true,
	},
	"addon_name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"credentials_id": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"credentials_name": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"roles": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: vaultRoleSchema,
		},
	},
}

var vaultRoleSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"service_account_name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"namespace": {
		Type:     schema.TypeString,
		Required: true,
	},
	"policies": {
		Type:     schema.TypeString,
		Required: true,
	},
}

func vaultAuthSchemaToVaultAuthSpec(vaultAuthSchema map[string]interface{}) map[string]interface{} {
	vaultAuthSpec := map[string]interface{}{
		"modelIndex": "VaultKubernetesAuthSpec",
		"name":       vaultAuthSchema["name"],
		"path":       vaultAuthSchema["path"],
		"addOnName":  vaultAuthSchema["addon_name"],
	}

	var rolesSpec []map[string]interface{}
	if _, ok := vaultAuthSchema["roles"]; ok {
		roles := vaultAuthSchema["roles"].([]interface{})
		for _, role := range roles {
			element, ok := role.(map[string]interface{})
			if ok {
				rolesSpec = append(rolesSpec, map[string]interface{}{
					"modelIndex":         "VaultRole",
					"name":               element["name"],
					"serviceAccountName": element["service_account_name"],
					"namespace":          element["namespace"],
					"policies":           element["policies"],
				},
				)
			}
		}
		vaultAuthSpec["roles"] = rolesSpec
	}

	credentialSpec := map[string]interface{}{
		"modelIndex": "VaultCredentials",
	}

	if ci, ok := vaultAuthSchema["credentials_id"]; ok {
		credentialSpec["id"] = ci
	}

	if cn, ok := vaultAuthSchema["credentials_name"]; ok {
		credentialSpec["name"] = cn
	}

	vaultAuthSpec["credentials"] = credentialSpec

	return vaultAuthSpec
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
		Schema: gkeClusterTypeSchema,
	}
}

func resourceGkeClusterTypeCreate(d *schema.ResourceData, meta interface{}) error {
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
	nodepools := d.Get("nodepools").([]interface{})
	clusterIpv4Cidr := d.Get("cluster_ipv4_cidr").(string)
	servicesIpv4Cidr := d.Get("services_ipv4_cidr").(string)
	cloudRun := d.Get("enable_cloud_run").(bool)
	enableNetworkPolicy := d.Get("enable_network_policy").(bool)
	httpLoadBalancing := d.Get("enable_http_load_balancing").(bool)
	enableVerticalPodAutoscaling := d.Get("enable_vertical_pod_autoscaling").(bool)
	horizontalPodAutoscaling := d.Get("enable_horizontal_pod_autoscaling").(bool)
	enableMaintenancePolicy := d.Get("enable_maintenance_policy").(bool)
	duration := d.Get("maintenance_duration").(string)
	startTime := d.Get("maintenance_start_time").(string)
	recurrence := d.Get("maintenance_recurrence").(string)
	exclusionTimeWindow := d.Get("maintenance_exclusion_timewindow")
	systemMetadata := d.Get("system_metadata")
	allowOverrideCredentials := d.Get("allow_override_credentials").(bool)
	clusterFieldOverride := d.Get("cluster_field_override")
	nodepoolFieldOverride := d.Get("nodepool_field_override")

	apiClient := meta.(client.Client)
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

	var gkeAddons []string
	if horizontalPodAutoscaling {
		gkeAddons = append(gkeAddons, "horizontalPodAutoscaling")
	}

	if httpLoadBalancing {
		gkeAddons = append(gkeAddons, "httpLoadBalancing")
	}

	if cloudRun {
		gkeAddons = append(gkeAddons, "cloudRunConfig")
	}

	fieldsToOverride := map[string]interface{}{
		"cluster":  clusterFieldOverride,
		"nodePool": nodepoolFieldOverride,
	}

	var nodeobjArr = make([]interface{}, 0)
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
					"gkeConfig": map[string]interface{}{
						"machineType":            element["machine_type"],
						"diskSize":               element["disk_size"],
						"serviceAccount":         element["service_account"],
						"enablePreemptibleNodes": element["enable_preemptible_nodes"],
						"modelIndex":             "GkeNodePoolConfig",
					},
				},
			}

			nodeobjArr = append(nodeobjArr, nodePoolObj)
		}
	}

	var addonsSpec []map[string]interface{}

	addonsSpec = append(addonsSpec, map[string]interface{}{
		"modelIndex":     "AddOnSpec",
		"name":           "kyverno",
		"addOnSelector":  "kyverno",
		"catalog":        "default-addon-catalog",
		"sequenceNumber": 1,
	},
	)

	if _, ok := d.GetOk("addons"); ok {
		addons := d.Get("addons").([]interface{})
		for _, addon := range addons {
			element, ok := addon.(map[string]interface{})
			if ok {
				addonsSpec = append(addonsSpec, map[string]interface{}{
					"modelIndex":     "AddOnSpec",
					"name":           element["name"],
					"addOnSelector":  element["addon_selector"],
					"catalog":        element["catalog"],
					"channel":        element["channel"],
					"namespace":      element["namespace"],
					"sequenceNumber": element["sequence_number"],
				},
				)
			}
		}
	}

	clusterTypeData := map[string]interface{}{
		"name":        name,
		"description": "",
		"modelIndex":  "ClusterType",
		"spec": map[string]interface{}{
			"clusterMode":    "providerManaged",
			"modelIndex":     "ClusterSpec",
			"version":        version,
			"cloud":          "googlecloudplatform",
			"systemMetadata": systemMetadata,
			"addons": map[string]interface{}{
				"dns":        false,
				"modelIndex": "AddOns",
				"other":      addonsSpec,
			},
			"cloudConfigSpec": map[string]interface{}{
				"credentials":              cloudCredID.UUID(),
				"allowOverrideCredentials": allowOverrideCredentials,
				"fieldsToOverride":         fieldsToOverride,
				"modelIndex":               "CloudConfigSpec",
				"gkeConfig": map[string]interface{}{
					"modelIndex":                   "GkeClusterConfig",
					"region":                       region,
					"zone":                         zone,
					"locationType":                 locationType,
					"defaultNodeLocations":         nodeLocations,
					"enableSecretsEncryption":      enableSecretsEncryption,
					"secretsEncryptionKey":         secretsEncryptionKey,
					"enableWorkloadIdentity":       enableWorkloadIdentity,
					"workloadPool":                 workloadPool,
					"network":                      network,
					"subnetwork":                   subnetwork,
					"clusterIpv4Cidr":              clusterIpv4Cidr,
					"servicesIpv4Cidr":             servicesIpv4Cidr,
					"enableNetworkPolicy":          enableNetworkPolicy,
					"enableMaintenancePolicy":      enableMaintenancePolicy,
					"duration":                     duration,
					"startTime":                    startTime,
					"recurrence":                   recurrence,
					"enableVerticalPodAutoscaling": enableVerticalPodAutoscaling,
					"exclusionTimeWindow":          exclusionTimeWindow,
					"addons":                       gkeAddons,
				},
				"nodePoolTypes": nodeobjArr,
			},
		},
	}

	if _, ok := d.GetOk("vault_auth"); ok {
		vl := d.Get("vault_auth").([]interface{})
		vault := vl[0].(map[string]interface{})
		clusterTypeData["spec"].(map[string]interface{})["vault"] = vaultAuthSchemaToVaultAuthSpec(vault)
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

var gkeClusterTypePaths = map[string]string{
	"version":                          "spec[0].version",
	"region":                           "spec[0].cloudConfigSpec[0].gkeConfig[0].region",
	"network":                          "spec[0].cloudConfigSpec[0].gkeConfig[0].network",
	"subnetwork":                       "spec[0].cloudConfigSpec[0].gkeConfig[0].subnetwork",
	"zone":                             "spec[0].cloudConfigSpec[0].gkeConfig[0].zone",
	"location_type":                    "spec[0].cloudConfigSpec[0].gkeConfig[0].locationType",
	"node_locations":                   "spec[0].cloudConfigSpec[0].gkeConfig[0].defaultNodeLocations",
	"enable_workload_identity":         "spec[0].cloudConfigSpec[0].gkeConfig[0].enableWorkloadIdentity",
	"enable_secrets_encryption":        "spec[0].cloudConfigSpec[0].gkeConfig[0].enableSecretsEncryption",
	"secrets_encryption_key":           "spec[0].cloudConfigSpec[0].gkeConfig[0].secretsEncryptionKey",
	"workload_pool":                    "spec[0].cloudConfigSpec[0].gkeConfig[0].workloadPool",
	"maintenance_start_time":           "spec[0].cloudConfigSpec[0].gkeConfig[0].startTime",
	"maintenance_recurrence":           "spec[0].cloudConfigSpec[0].gkeConfig[0].recurrence",
	"maintenance_duration":             "spec[0].cloudConfigSpec[0].gkeConfig[0].duration",
	"maintenance_exclusion_timewindow": "spec[0].cloudConfigSpec[0].gkeConfig[0].exclusionTimeWindow",
	"enable_maintenance_policy":        "spec[0].cloudConfigSpec[0].gkeConfig[0].enableMaintenancePolicy",
	"cluster_ipv4_cidr":                "spec[0].cloudConfigSpec[0].gkeConfig[0].clusterIpv4Cidr",
	"services_ipv4_cidr":               "spec[0].cloudConfigSpec[0].gkeConfig[0].servicesIpv4Cidr",
	"enable_vertical_pod_autoscaling":  "spec[0].cloudConfigSpec[0].gkeConfig[0].enableVerticalPodAutoscaling",
	"enable_network_policy":            "spec[0].cloudConfigSpec[0].gkeConfig[0].enableNetworkPolicy",
}

var nodePoolTypePaths = map[string]string{
	"machine_type":             "spec[0].gkeConfig[0].machineType",
	"disk_size":                "spec[0].gkeConfig[0].diskSize",
	"enable_preemptible_nodes": "spec[0].gkeConfig[0].enablePreemptibleNodes",
	"service_account":          "spec[0].gkeConfig[0].serviceAccount",
	"node_labels":              "spec[0].nodeLabels",
	"node_annotations":         "spec[0].nodeAnnotations",
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
		s := gkeClusterTypeSchema[field]
		err = updateData(field, d, s, path, clusterTypeData)
		if err != nil {
			return fmt.Errorf("failed to update field %s from %s: %v", field, path, err)
		}
	}

	return nil
}

var gkeAttributeMap = map[string]string{
	"version":                         "version",
	"region":                          "region",
	"network":                         "network",
	"subnetwork":                      "subnetwork",
	"zone":                            "zone",
	"location_type":                   "locationType",
	"node_locations":                  "defaultNodeLocations",
	"enable_workload_identity":        "enableWorkloadIdentity",
	"enable_secrets_encryption":       "enableSecretsEncryption",
	"secrets_encryption_key":          "secretsEncryptionKey",
	"workload_pool":                   "workloadPool",
	"start_time":                      "startTime",
	"exclusion_timewindow":            "exclusionTimeWindow",
	"enable_maintenance_policy":       "enableMaintenancePolicy",
	"cluster_ipv4_cidr":               "clusterIpv4Cidr",
	"services_ipv4_cidr":              "servicesIpv4Cidr",
	"enable_vertical_pod_autoscaling": "enableVerticalPodAutoscaling",
	"duration":                        "duration",
	"enable_network_policy":           "enableNetworkPolicy",
}

var nodePoolAttributeMap = map[string]string{
	"machine_type":             "machineType",
	"disksize":                 "diskSize",
	"enable_preemptible_nodes": "enablePreemptibleNodes",
	"service_account":          "serviceAccount",
	"node_labels":              "nodeLabels",
	"node_annotations":         "nodeAnnotations",
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
	gkeConfigChanges := buildChanges(d, gkeAttributeMap,
		"region",
		"network",
		"subnetwork",
		"zone",
		"location_type",
		"node_locations",
		"enable_workload_identity",
		"enable_secrets_encryption",
		"secrets_encryption_key",
		"workload_pool",
		"maintenance_sstart_time",
		"maintenance_exclusion_timewindow",
		"enable_maintenance_policy",
		"cluster_ipv4_cidr",
		"services_ipv4_cidr",
		"enable_vertical_pod_autoscaling",
		"maintenance_duration",
		"enable_network_policy")

	if len(gkeConfigChanges) > 0 {
		err := updateDescendant(apiClient, clusterTypeID, "GkeClusterConfig", gkeConfigChanges)
		if err != nil {
			return err
		}
	}

	// update NodePool
	nodePoolChanges := buildChanges(d, nodePoolAttributeMap,
		"machine_type",
		"disk_size",
		"enable_preemptible_nodes",
		"service_account",
		"node_labels",
		"node_annotations")

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
