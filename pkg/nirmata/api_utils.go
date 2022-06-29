package nirmata

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func delete(d *schema.ResourceData, m interface{}, s client.Service, model string, params map[string]string) error {
	apiClient := m.(client.Client)
	id := clientID(d, s, model)
	if err := apiClient.Delete(id, params); err != nil {
		if !strings.Contains(err.Error(), "404") {
			return err
		}
	}

	d.SetId("")
	return nil
}

func clientID(d *schema.ResourceData, s client.Service, model string) client.ID {
	uuid := d.Id()
	return client.NewID(s, model, uuid)
}

// waitForClusterState waits until cluster is created or has failed
func waitForClusterState(apiClient client.Client, maxTime time.Duration, clusterID client.ID) (string, error) {
	states := []interface{}{"ready", "failed"}
	state, err := apiClient.WaitForStates(clusterID, "state", states, maxTime, "")
	if err != nil {
		return "", err
	}

	return state.(string), nil
}

// waitForRollloutState waits until rollout is created or has failed
func waitForRolloutState(apiClient client.Client, maxTime time.Duration, rolloutID client.ID) (string, error) {
	states := []interface{}{"completed", "failed"}
	state, err := apiClient.WaitForStates(rolloutID, "state", states, maxTime, "")
	if err != nil {
		return "", err
	}

	return state.(string), nil
}

func waitForConfigurationState(apiClient client.Client, maxTime time.Duration, clusterID client.ID) (string, error) {
	states := []interface{}{"completed", "failed"}
	state, err := apiClient.WaitForStates(clusterID, "configurationState", states, maxTime, "")
	if err != nil {
		return "", err
	}

	return state.(string), nil
}

func getRolloutStatus(api client.Client, rolloutID client.ID) (string, error) {
	rolloutData, err := api.Get(rolloutID, client.NewGetOptions(nil, nil))
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve  rollout details: %v", err)
		return "", err
	}
	statusMsg := rolloutData["errorInfo"].(string)
	return statusMsg, nil
}

func getGitUpstreamStatus(api client.Client, rolloutID client.ID) (string, error) {
	rolloutData, err := api.Get(rolloutID, client.NewGetOptions(nil, nil))
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve git upstream details: %v", err)
		return "", err
	}
	if rolloutData["lastGitSyncError"] == nil {
		return "", fmt.Errorf(" [ERROR] - git sync failed")
	}
	statusMsg := rolloutData["lastGitSyncError"].(string)
	return statusMsg, nil
}

func getPolicyGroupStatus(api client.Client, ID client.ID) (string, error) {
	data, err := api.Get(ID, client.NewGetOptions(nil, nil))
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve policy group  details: %v", err)
		return "", err
	}
	if data["syncError"] == nil {
		return "", fmt.Errorf(" [ERROR] - policy group sync failed")
	}
	statusMsg := data["syncError"].(string)
	return statusMsg, nil
}

func waitForDeletedState(apiClient client.Client, maxTime time.Duration, clusterID client.ID) (string, error) {
	states := []interface{}{"deleted"}
	state, err := apiClient.WaitForStates(clusterID, "state", states, maxTime, "")
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] cluster deleted")
			return "", nil
		}
		return "", err
	}

	return state.(string), nil
}

func getClusterStatus(api client.Client, clusterID client.ID) (string, error) {
	clusterData, err := api.Get(clusterID, client.NewGetOptions(nil, nil))
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve failed cluster details: %v", err)
		return "", err
	}

	statusList := clusterData["status"].([]interface{})
	statusMsg := ""
	for _, s := range statusList {
		statusMsg += s.(string)
	}

	return statusMsg, nil
}

func waitForGitSyncStatus(apiClient client.Client, maxTime time.Duration, rolloutID client.ID) (string, error) {
	states := []interface{}{"success", "failed"}
	state, err := apiClient.WaitForStates(rolloutID, "lastGitSyncStatus", states, maxTime, "")
	if err != nil {
		return "", err
	}

	return state.(string), nil
}

func waitForPolicySetSyncStatus(apiClient client.Client, maxTime time.Duration, ID client.ID) (string, error) {
	states := []interface{}{"success", "failed"}
	state, err := apiClient.WaitForStates(ID, "syncState", states, maxTime, "")
	if err != nil {
		return "", err
	}

	return state.(string), nil
}

func waitForAddonState(apiClient client.Client, maxTime time.Duration, addonID client.ID) (string, error) {
	states := []interface{}{"running", "failed"}

	state, err := apiClient.WaitForStates(addonID, "runningState", states, maxTime, "")
	if err != nil {
		return "", err
	}

	return state.(string), nil
}

// extracts an object from a TXN POST operation
func extractCreateFromTxnResult(data map[string]interface{}, modelIndex string) (client.Object, error) {
	if data["create"] == nil {
		return nil, fmt.Errorf("no created objects in result %s", data)
	}

	createList := data["create"].([]interface{})
	for _, c := range createList {
		cm := c.(map[string]interface{})
		mi := cm["modelIndex"].(string)
		if !strings.EqualFold(mi, modelIndex) {
			continue
		}

		createdObj, err := client.NewObject(cm)
		if err != nil {
			return nil, fmt.Errorf("failed to parse cluster type %v: %v", c, err)
		}

		return createdObj, nil
	}

	return nil, fmt.Errorf("modelIndex %s not found in create result", modelIndex)
}

func getNodePoolType(apiClient client.Client, clusterTypeID client.ID) (map[string]interface{}, error) {
	var err error
	nodePool, err := apiClient.GetDescendants(clusterTypeID, "NodePoolType", &client.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch NodePoolType for %+v: %v", clusterTypeID.Map(), err)
	}
	if len(nodePool) < 1 {
		return nil, fmt.Errorf("no NodePoolType found for %+v: %v", clusterTypeID.Map(), err)
	}

	nodePoolTypeID, err := client.NewObject(nodePool[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode NodePoolType  %+v: %v", nodePool[0], err)
	}
	return apiClient.Get(nodePoolTypeID.ID(), &client.GetOptions{nil, nil, client.OutputModeExportDetails})
}

func fetchID(apiClient client.Client, service client.Service, modelIndex, nameOrID string) (client.ID, error) {
	if isUUID(nameOrID) {
		return client.NewID(service, modelIndex, nameOrID), nil
	}

	return apiClient.QueryByName(service, modelIndex, nameOrID)
}

func isUUID(id string) bool {
	if _, err := uuid.Parse(id); err == nil {
		return true
	}

	return false
}

func deleteObj(d *schema.ResourceData, meta interface{}, service client.Service, modelIndex string) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)

	id, err := apiClient.QueryByName(service, modelIndex, name)
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return err
	}

	if err := apiClient.Delete(id, nil); err != nil {
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

func waitForPolicyDeploySetSyncStatus(apiClient client.Client, maxTime time.Duration, rolloutID client.ID) (string, error) {
	states := []interface{}{"completed", "failed"}
	state, err := apiClient.WaitForStates(rolloutID, "rolloutState", states, maxTime, "")
	if err != nil {
		return "", err
	}

	return state.(string), nil
}

func getPolicyDeployGroupStatus(api client.Client, ID client.ID) (string, error) {
	data, err := api.Get(ID, client.NewGetOptions(nil, nil))
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve deployed policy group details: %v", err)
		return "", err
	}
	if data["rolloutError"] == nil {
		return "", fmt.Errorf(" [ERROR] - policy group rollout failed")
	}
	statusMsg := data["rolloutError"].(string)
	return statusMsg, nil
}
