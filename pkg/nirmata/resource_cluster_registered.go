package nirmata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"os"

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
	"controller_yamls_count": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"delete_action": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validateDeleteAction,
	},
	"labels": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
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
    f, err := os.Create("/tmp/tf.log")
    defer f.Close()
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	labels := d.Get("labels")
	clusterType := d.Get("cluster_type").(string)

	deleteAction := d.Get("delete_action").(string)
	if deleteAction == "" {
		d.Set("delete_action", "remove")
	}

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
		"labels":       labels,
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

	clusterData, err := apiClient.QueryByName(client.ServiceClusters, "KubernetesCluster", name)
	if err != nil {
	    errString := fmt.Sprintf("[ERROR] - %v", err)
	    f.WriteString(errString)
	    f.Sync()
    	log.Printf("[ERROR] - %v", err)
    	return err
    }
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

	path, count, ferr := writeToTempDir([]byte(yaml))
	if ferr != nil {
		return fmt.Errorf("failed to write temporary files: %v", ferr)
	}

	log.Printf("[INFO] - wrote temporary YAMLs files to %s:", path)

	d.Set("controller_yamls_count", count)
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

func writeToTempDir(data []byte) (path string, count int, err error) {
	path, err = ioutil.TempDir("", "controller-")
	if err != nil {
		fmt.Println(err)
	}

	result := strings.Split(string(data), "---")
	count = 0
	for i := range result {
		if result[i] != "" {
			fmt.Println(result[i])
			f, err := ioutil.TempFile(path, "temp-")
			if err != nil {
				return "", 0, fmt.Errorf("cannot create temporary file: %v", err)
			}

			if _, err = f.Write([]byte(result[i])); err != nil {
				return "", 0, fmt.Errorf("failed to write temporary file %s: %v", f.Name(), err)
			}

			count += 1
			f.Close()
		}
	}

	return
}
