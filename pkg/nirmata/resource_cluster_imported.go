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
}

func resourceClusterImported() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterImportedCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterUpdate,
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
	credentials := d.Get("credentials").(string)
	region := d.Get("region").(string)
	clusterType := d.Get("cluster_type").(string)

	cloudCredID, err := apiClient.QueryByName(client.ServiceClusters, "CloudCredentials", credentials)
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return err
	}

	data := map[string]interface{}{
		"mode":                "providerManaged",
		"clusterTypeSelector": clusterType,
		"credentialsRef":      cloudCredID.UUID(),
		"clusters": map[string]interface{}{
			name: map[string]interface{}{
				"name":   name,
				"region": region,
			},
		},
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
