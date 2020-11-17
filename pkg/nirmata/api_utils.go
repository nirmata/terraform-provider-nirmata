package nirmata

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

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

// newUUID generates a random UUID
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits;
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random);
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// waitForClusterState waits until cluster is created or has failed
func waitForClusterState(apiClient client.Client, maxTime time.Duration, clusterID client.ID) (string, error) {
	for {
		newState, timeout, err := apiClient.WaitForStateChange(clusterID, "state", maxTime)
		if err != nil {
			return "", err
		}

		if timeout {
			return "", fmt.Errorf("timeout")
		}

		state := newState.(string)
		if state == "failed" || state == "ready" {
			return state, nil
		}
	}
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
