package nirmata

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

var importedClusterSchema = map[string]*schema.Schema{
	"name": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validateName,
	},
	"cluster": {
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
	"cluster_type": {
		Type:     schema.TypeString,
		Required: true,
	},
	"project": {
		Type:     schema.TypeString,
		Required: true,
	},
	"delete_action": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "remove",
		ValidateFunc: validateDeleteAction,
	},
	"system_metadata": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
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
	"labels": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
}

func resourceClusterImported() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterImportedCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterImportUpdate,
		Delete: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: importedClusterSchema,
	}
}

func resourceClusterImportedCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	labels := d.Get("labels")
	cluster := d.Get("cluster").(string)
	credentials := d.Get("credentials").(string)
	region := d.Get("region").(string)
	clusterType := d.Get("cluster_type").(string)
	project := d.Get("project").(string)
	systemMetadata := d.Get("system_metadata")
	deleteAction := d.Get("delete_action").(string)
	if deleteAction == "" {
		d.Set("delete_action", "remove")
	}

	cloudCredID, err := apiClient.QueryByName(client.ServiceClusters, "CloudCredentials", credentials)
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return err
	}

	clusterJson := map[string]interface{}{}
	if cluster == "gke" {
		clusterJson = map[string]interface{}{
			name: map[string]interface{}{
				"name":           name,
				"region":         region,
				"project":        project,
				"id":             name,
				"systemMetadata": systemMetadata,
			},
		}
	} else if cluster == "eks" {
		b, _, err := api.GetURLWithID(credentialsID, "fetchClusters?region="+region)
		if err != nil {
			fmt.Println(err)
			return err
		}
		data := map[string]interface{}{}
		if err := json.Unmarshal(b, &data); err != nil {
			return err
		}
		nodeapps := extractCollection(data["clusters"])
		for _, v := range nodeapps {
			if name == v["name"] {
				clusterJson = map[string]interface{}{
					name: v,
				}
			}
		}
	}

	data := map[string]interface{}{
		"mode":                "providerManaged",
		"labels":              labels,
		"systemMetadata":      systemMetadata,
		"clusterTypeSelector": clusterType,
		"credentialsRef":      cloudCredID.UUID(),
		"clusters":            clusterJson,
	}

	log.Printf("[DEBUG] - importing cluster %s with %+v", name, data)
	action, err := apiClient.PostFromJSON(client.ServiceClusters, "ImportClustersAction", data, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to import cluster %s with data %v: %v", name, data, err)
		return fmt.Errorf("failed to import cluster %s: %v", name, err.Error())
	}

	actionObj, objErr := client.NewObject(action)
	if objErr != nil {
		log.Printf("[ERROR] - failed to convert %v: %v", action, err)
		return err
	}

	actionID := actionObj.ID()
	state, actionObj, waitErr := waitForImportClustersAction(apiClient, d.Timeout(schema.TimeoutCreate), actionID)
	if waitErr != nil {
		log.Printf("[ERROR] - failed to import cluster. Error - %v", waitErr)
		return waitErr
	}

	if strings.EqualFold("failed", state) {
		log.Printf("[ERROR] - failed to import cluster - %v", actionObj.Data())
		progress := actionObj.Data()["progress"]
		return fmt.Errorf("cluster import failed: %s", getJSON(progress))
	}

	if state == "success" {
		clustersField := actionObj.GetString("clustersField")
		log.Printf("Got cluster data: %+v", clustersField)
	}

	clusterID, err := apiClient.QueryByName(client.ServiceClusters, "KubernetesCluster", name)
	if err != nil {
		log.Printf("[ERROR] - failed to fetch cluster %v: %v", name, err)
		return waitErr
	}

	d.SetId(clusterID.UUID())
	log.Printf("[INFO] registered cluster %s with ID %s", name, clusterID)
	return nil
}

var clustervaultMap = map[string]string{
	"vault_auth": "vault",
	"labels":     "labels",
}

func resourceClusterImportUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	clusterID := client.NewID(client.ServiceClusters, "KubernetesCluster", d.Id())
	vaultAuthChanges := buildChanges(d, clustervaultMap, "vault_auth")
	labelsChanges := buildChanges(d, clustervaultMap, "labels")
	if len(vaultAuthChanges) > 0 {
		if err := updateVaultAddon(d, apiClient, clusterID); err != nil {
			log.Printf("[ERROR] - failed to update cluster  vault with data : %v", err)
			return err
		}
	}

	if len(labelsChanges) > 0 {
		if err := updateClusterLabels(d, apiClient); err != nil {
			log.Printf("[ERROR] - failed to update labels with data : %v", err)
			return err
		}
	}

	return nil
}
func waitForImportClustersAction(apiClient client.Client, maxTime time.Duration, actionID client.ID) (string, client.Object, error) {
	states := []interface{}{"success", "failed"}
	stateRaw, err := apiClient.WaitForStates(actionID, "status", states, maxTime, "")
	if err != nil {
		log.Printf("[ERROR] - failed check states: %v", err)
		return "", nil, err
	}

	state := stateRaw.(string)

	data, err := apiClient.Get(actionID, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to retrieve %v: %v", actionID.Map(), err)
		return "", nil, err
	}

	actionObj, objErr := client.NewObject(data)
	if objErr != nil {
		log.Printf("[ERROR] - failed to convert %v: %v", data, err)
		return "", nil, err
	}

	return state, actionObj, nil
}

func getJSON(data interface{}) string {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		// ...
	}

	return string(jsonBytes)
}
