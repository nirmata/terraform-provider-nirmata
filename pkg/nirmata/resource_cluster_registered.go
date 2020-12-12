package nirmata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

	log.Printf("[DEBUG] - importing cluster %s with %+v", name, data)
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
		log.Printf("Failed to fetch controller YAML %s: %v \n", name, err)
		return err
	}

	yaml, yamlErr := getCtrlYAML(b)
	if yamlErr != nil {
		log.Printf("Failed to decode controller YAML %s: %v \n", name, yamlErr)
		return yamlErr
	}

	file, ferr := writeToTempFile([]byte(yaml))
	if ferr != nil {
		return fmt.Errorf("Failed to write temp file: %v", ferr)
	}

	defer os.Remove(file.Name())
	commandArgs := []string{"apply", "-f", file.Name()}
	bytes, execErr := exec.Command("kubectl", commandArgs...).CombinedOutput()
	if execErr != nil {
		return fmt.Errorf("Failed to execute command %v: %v %s", commandArgs, execErr, string(bytes))
	}

	clusterID := client.NewID(client.ServiceClusters, "KubernetesCluster", clusterUUID)
	state, waitErr := waitForClusterState(apiClient, d.Timeout(schema.TimeoutCreate), clusterID)
	if waitErr != nil {
		log.Printf("[ERROR] - failed to check cluster status. Error - %v", waitErr)
		return waitErr
	}

	if strings.EqualFold("failed", state) {
		status, err := getClusterStatus(apiClient, clusterID)
		if err != nil {
			log.Printf("[ERROR] - failed to retrieve cluster failure details: %v", err)
			return fmt.Errorf("cluster creation failed")
		}

		return fmt.Errorf("cluster creation failed: %s", status)
	}

	log.Printf("[INFO] registered cluster %s with ID %s", name, clusterID)
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

	return "", fmt.Errorf("Invalid controller YAML: %v", m)
}

func writeToTempFile(data []byte) (f *os.File, err error) {
	f, err = ioutil.TempFile(os.TempDir(), "temp-")
	if err != nil {
		return f, fmt.Errorf("Cannot create temporary file: %v", err)
	}

	if _, err = f.Write(data); err != nil {
		return f, fmt.Errorf("Failed to write temporary file %s: %v", f.Name(), err)
	}

	return
}
