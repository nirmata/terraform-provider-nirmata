package nirmata

import (
	"crypto/rand"
	"fmt"
	"io"
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

// waitForState Wait until cluster is not created or failed
func waitForState(apiClient client.Client,maxTime time.Duration,name string) error {
	clusterTypeID, err := apiClient.QueryByName(client.ServiceClusters, "clustertypes", name)
	if err != nil {
		fmt.Printf("Error ", err)
		return err
	}

	for {
		_,newState,err := apiClient.WaitForStateChange(clusterTypeID,"state",maxTime,"Failed to get cluster status")
		if err != nil {
			continue
		}
		if newState.(string) == "failed" {
			return err
		}else if newState.(string) == "ready"{
			break;
		}
	}

	return nil
}