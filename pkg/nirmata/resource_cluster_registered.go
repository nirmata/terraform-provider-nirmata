package nirmata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

var accessControlSchema = map[string]*schema.Schema{
	"entity_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validateEntity,
	},
	"permission": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validatePermission,
	},
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
}

var ownerInfoSchema = map[string]*schema.Schema{
	"owner_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validateEntity,
	},
	"owner_name": {
		Type:     schema.TypeString,
		Required: true,
	},
}

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
	"endpoint": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"owner_info": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: ownerInfoSchema,
		},
	},
	"access_control_list": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: accessControlSchema,
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

func fetchOwnerId(owner_type, owner_name string, users, teams []map[string]interface{}) (string, error) {
	var data []map[string]interface{}
	if owner_type == "user" {
		data = users
	} else {
		data = teams
	}
	for _, d := range data {
		if owner_name == d["name"] {
			id := d["id"].(string)
			log.Printf("[DEBUG] - Found match for owner_name %s with id %s", owner_name, id)
			return id, nil
		}
	}
	return "", fmt.Errorf("[ERROR] - id not found for %s/%s", owner_type, owner_name)
}

func makeAccessControlListObj(ownerInfo []interface{}, users, teams []map[string]interface{}) (map[string]interface{}, error) {
	var accessControlListObj map[string]interface{}
	// Set owner_type to user as the default when no owner_info is provided
	var owner_type string = "user"
	var owner_name string

	if ownerInfo != nil && len(ownerInfo) != 0 {
		for _, o := range ownerInfo {
			elem, ok := o.(map[string]interface{})
			if ok {
				owner_type = elem["owner_type"].(string)
				owner_name = elem["owner_name"].(string)
				owner_id, err := fetchOwnerId(owner_type, owner_name, users, teams)
				if err != nil {
					log.Printf("[ERROR] - No entity of type %s with name %s found in the cluster", owner_type, owner_name)
					return nil, fmt.Errorf("[ERROR] - Invalid owner '%s/%s' provided", owner_type, owner_name)
				}
				accessControlListObj = map[string]interface{}{
					"modelIndex": "AccessControlList",
					"ownerType":  owner_type,
					"ownerName":  owner_name,
					"ownerId":    owner_id,
				}
			}
		}
	} else {
		log.Printf("[DEBUG] - owner_info is empty")
		accessControlListObj = map[string]interface{}{
			"modelIndex": "AccessControlList",
			"ownerType":  owner_type,
		}
	}
	return accessControlListObj, nil
}

func makeAccessControlList(accessControlList []interface{}, users, teams []map[string]interface{}) ([]interface{}, error) {
	var acObArr = make([]interface{}, 0)

	if accessControlList != nil {
		for _, ac := range accessControlList {
			log.Printf("[DEBUG] - ac %v", ac)
			elem, ok := ac.(map[string]interface{})
			if ok {
				entity_type := elem["entity_type"].(string)
				entity_name := elem["name"].(string)
				entity_id, err := fetchOwnerId(entity_type, entity_name, users, teams)
				if err != nil {
					log.Printf("[ERROR] - No entity of type %s with name %s found in the cluster", entity_type, entity_name)
					return nil, fmt.Errorf("[ERROR] - Invalid entity '%s/%s' provided", entity_type, entity_name)
				}
				acOb := map[string]interface{}{
					"modelIndex": "AccessControl",
					"entityType": entity_type,
					"entityName": entity_name,
					"permission": elem["permission"],
					"entityId":   entity_id,
				}
				log.Printf("[DEBUG] - Printing accessControlList - acOb %v", acOb)
				acObArr = append(acObArr, acOb)
			}
		}
	}
	return acObArr, nil
}

func resourceClusterRegisteredCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	labels := d.Get("labels")
	endpoint := d.Get("endpoint").(string)
	clusterType := d.Get("cluster_type").(string)

	users, err := getUsers(apiClient)
	if err != nil {
		log.Printf("[ERROR] - failed to retrieve users: %v", err)
		return fmt.Errorf("[ERROR] - Fetching users failed")
	}
	log.Println("[DEBUG] - Printing users")
	for _, u := range users {
		for k, v := range u {
			log.Printf("[DEBUG] - key: %s - value: %s", k, v)
		}
	}

	teams, err := getTeams(apiClient)
	if err != nil {
		log.Printf("[ERROR] - failed to retrieve teams: %v", err)
		return fmt.Errorf("[ERROR] - Fetching teams failed")
	}
	log.Println("[DEBUG] - Printing teams")
	for _, t := range teams {
		for k, v := range t {
			log.Printf("[DEBUG] - key: %s - value: %s", k, v)
		}
	}

	ownerInfo := d.Get("owner_info").([]interface{})
	log.Printf("[DEBUG] - Printing owner_info from resource %v", ownerInfo)

	accessControlListObj, err := makeAccessControlListObj(ownerInfo, users, teams)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] - Printing accessControlListObj %v", accessControlListObj)

	accessControlList := d.Get("access_control_list").([]interface{})
	log.Printf("[DEBUG] - Printing accessControlList %v", accessControlList)
	acObArr, err := makeAccessControlList(accessControlList, users, teams)
	if err != nil {
		return err
	}
	accessControlListObj["accessControls"] = acObArr
	log.Println("[DEBUG] - Printing accessControlList", accessControlListObj)

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

	data["config"] = map[string]interface{}{
		"modelIndex":    "ClusterConfig",
		"version":       spec["version"],
		"cloudProvider": spec["cloud"],
		"endpoint":      endpoint,
	}

	var aclArr = make([]interface{}, 0)
	aclArr = append(aclArr, accessControlListObj)
	data["accessControlList"] = aclArr

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
	for _, r := range result {
		if r == "" {
			continue
		}

		f, err := ioutil.TempFile(path, "temp-")
		if err != nil {
			return "", 0, fmt.Errorf("cannot create temporary file: %v", err)
		}

		if _, err = f.Write([]byte(r)); err != nil {
			return "", 0, fmt.Errorf("failed to write temporary file %s: %v", f.Name(), err)
		}

		count += 1
		f.Close()
	}
	return
}
