package nirmata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

var registeredClusterSchema = map[string]*schema.Schema{
	"name": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validateName,
	},
	"cluster_type": {
		Type:     schema.TypeString,
		Required: true,
	},
	"controller_yamls": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"state": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"controller_yamls_folder": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"delete_action": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "remove",
		ValidateFunc: validateDeleteAction,
	},
}

func resourceClusterRegistered() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterRegisteredCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterUpdate,
		Delete: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: registeredClusterSchema,
	}
}

func resourceClusterRegisteredCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	clusterType := d.Get("cluster_type").(string)

	typeID, err := apiClient.QueryByName(client.ServiceClusters, "ClusterType", clusterType)
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return err
	}
	spec, err := apiClient.GetRelation(typeID, "clusterSpecs")
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return err
	}

	mode := spec["clusterMode"]
	if mode != "discovered" {
		return err
	}

	data := map[string]interface{}{
		"name":         name,
		"mode":         mode,
		"typeSelector": clusterType,
	}

	log.Printf("[DEBUG] - registering cluster %s with %+v", name, data)
	clusterObj, err := apiClient.PostFromJSON(client.ServiceClusters, "KubernetesCluster", data, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to register cluster %s with data %v: %v", name, data, err)
		return err
	}

	clusterUUID := clusterObj["id"].(string)
	d.SetId(clusterUUID)

	clusterData, _ := apiClient.QueryByName(client.ServiceClusters, "KubernetesCluster", name)
	b, _, err := apiClient.GetURLWithID(clusterData, "controllerYAML")
	if err != nil {
		log.Printf("[ERROR] - Failed to fetch controller YAML %s: %v \n", name, err)
		return err
	}

	yaml, yamlErr := getCtrlYAML(b)
	if yamlErr != nil {
		log.Printf("[ERROR] - Failed to decode controller YAML %s: %v \n", name, yamlErr)
		return yamlErr
	}

	d.Set("controller_yamls", yaml)
	d.Set("state", clusterObj["state"])
	path, ferr := writeToTempDir([]byte(yaml))
	if ferr != nil {
		return fmt.Errorf("failed to write temp file: %v", ferr)
	}

	d.Set("controller_yamls_folder", path)
	return nil
}

func getCtrlYAML(b []byte) (string, error) {
	m := make(map[string]string)
	if err := json.Unmarshal(b, &m); err != nil {
		return "", err
	}

	for _, v := range m {
		return v, nil
	}

	return "", fmt.Errorf("invalid controller YAML: %v", m)
}

func writeToTempDir(data []byte) (path string, err error) {
	path, err = ioutil.TempDir("", "controller-")
	if err != nil {
		fmt.Println(err)
	}

	defer os.RemoveAll(path)
	result := strings.Split(string(data), "---")
	for i := range result {
		if result[i] != "" {
			fmt.Println(result[i])
			f, err := ioutil.TempFile(path, "temp-")
			if err != nil {
				return "", fmt.Errorf("cannot create temporary file: %v", err)
			}

			if _, err = f.Write([]byte(result[i])); err != nil {
				return "", fmt.Errorf("failed to write temporary file %s: %v", f.Name(), err)
			}

			f.Close()
		}
	}

	return
}
