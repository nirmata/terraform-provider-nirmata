package nirmata

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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
		Optional: true,
		Computed: true,
	},
	"controller_yamls_count": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"controller_ns_yamls_count": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"controller_crd_yamls_count": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"controller_deploy_yamls_count": {
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
		var name string
		if owner_type == "user" {
			name = d["email"].(string)
		} else {
			name = d["name"].(string)
		}
		if owner_name == name {
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
					return nil, fmt.Errorf("Invalid owner '%s/%s' provided", owner_type, owner_name)
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

func makeAccessControls(accessControlList []interface{}, users, teams []map[string]interface{}) ([]interface{}, error) {
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
					return nil, fmt.Errorf("Invalid entity '%s/%s' provided", entity_type, entity_name)
				}
				acOb := map[string]interface{}{
					"modelIndex": "AccessControl",
					"entityType": entity_type,
					"entityName": entity_name,
					"permission": elem["permission"],
					"entityId":   entity_id,
				}
				log.Printf("[DEBUG] - accessControlList: %v", acOb)
				acObArr = append(acObArr, acOb)
			}
		}
	}
	return acObArr, nil
}

/*
  - An example of accesscontrollist is this
  - "accessControlList": [
    {
    "ownerType": "user",
    "ownerId": "ce36b44d-74e0-417f-8b8f-7eb6f793f344",
    "ownerName": "user1@foo.com",
    "modelIndex": "AccessControlList",
    "accessControls": [
    {
    "entityType": "user",
    "entityId": "331977cc-2086-455b-bf4c-614cab868616",
    "permission": "admin",
    "entityName": "user2@foo.com",
    "modelIndex": "AccessControl"
    },
    {
    "entityType": "team",
    "entityId": "e2424edf-d7ee-4b86-9ec4-cde4659bd231",
    "permission": "edit",
    "entityName": "team1",
    "modelIndex": "AccessControl"
    },
    {
    "entityType": "user",
    "entityId": "1e71e1f4-102c-47a2-87f7-3398e0d92472",
    "permission": "view",
    "entityName": "user3@foo.com",
    "modelIndex": "AccessControl"
    }
    ]
    }
    ]
    This method creates the above structure
*/
func makeAccessControlList(d *schema.ResourceData, users, teams []map[string]interface{}) ([]interface{}, error) {
	var aclArr = make([]interface{}, 0)
	var accessControlListObj map[string]interface{}

	ownerInfo := d.Get("owner_info").([]interface{})
	log.Printf("[DEBUG] - owner_info %v", ownerInfo)

	accessControlListObj, err := makeAccessControlListObj(ownerInfo, users, teams)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] - accessControlListObj %v", accessControlListObj)

	accessControlList := d.Get("access_control_list").([]interface{})
	log.Printf("[DEBUG] - accessControlList %v", accessControlList)
	acs, err := makeAccessControls(accessControlList, users, teams)
	if err != nil {
		return nil, err
	}
	accessControlListObj["accessControls"] = acs
	log.Println("[DEBUG] - accessControlList", accessControlListObj)

	aclArr = append(aclArr, accessControlListObj)
	return aclArr, nil
}

func resourceClusterRegisteredCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	labels := d.Get("labels")
	endpoint := d.Get("endpoint").(string)
	clusterType := d.Get("cluster_type").(string)
	controller_yamls_folder := d.Get("controller_yamls_folder").(string)

	users, err := getUsers(apiClient)
	if err != nil {
		log.Printf("[ERROR] - failed to retrieve users: %v", err)
		return fmt.Errorf("Fetching users failed")
	}
	log.Println("[DEBUG] - users")
	for _, u := range users {
		for k, v := range u {
			log.Printf("[DEBUG] - key: %s - value: %s", k, v)
		}
	}

	teams, err := getTeams(apiClient)
	if err != nil {
		log.Printf("[ERROR] - failed to retrieve teams: %v", err)
		return fmt.Errorf("Fetching teams failed")
	}
	log.Println("[DEBUG] - teams")
	for _, t := range teams {
		for k, v := range t {
			log.Printf("[DEBUG] - key: %s - value: %s", k, v)
		}
	}

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

	accessControlList, err := makeAccessControlList(d, users, teams)
	if err != nil {
		log.Printf("[ERROR] - failed to create access control list")
		return err
	}
	data["accessControlList"] = accessControlList

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

	path, count, count_ns, count_crd, count_deploy, ferr := writeToTempDir(controller_yamls_folder, []byte(yaml))
	if ferr != nil {
		return fmt.Errorf("failed to write temporary files: %v", ferr)
	}

	log.Printf("[INFO] - wrote temporary YAMLs files to %s:", path)

	d.Set("controller_yamls_count", count)
	d.Set("controller_ns_yamls_count", count_ns)
	d.Set("controller_crd_yamls_count", count_crd)
	d.Set("controller_deploy_yamls_count", count_deploy)
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

func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func writeToTempDir(controller_yamls_folder string, data []byte) (path string, count, count_ns, count_crd, count_deploy int, err error) {
	if controller_yamls_folder != "" {
		if _, err := os.Stat(controller_yamls_folder); os.IsNotExist(err) {
			return "", 0, 0, 0, 0, fmt.Errorf("Folder '%s' does not exist", controller_yamls_folder)
		}
		empty, _ := IsDirEmpty(controller_yamls_folder)
		if !empty {
			return "", 0, 0, 0, 0, fmt.Errorf("Folder '%s' is not empty", controller_yamls_folder)
		}
		path = controller_yamls_folder
	} else {
		path, err = ioutil.TempDir(controller_yamls_folder, "controller-")
		if err != nil {
			fmt.Println(err)
		}
	}

	result := strings.Split(string(data), "---")
	count = 0
	count_ns = 0
	count_crd = 0
	count_deploy = 0
	for _, r := range result {
		if r == "" {
			continue
		}

		fileContents := fmt.Sprintf("%s", r)
		prefix := "temp%s"
		log.Printf("[DEBUG] fileContents %s", fileContents)
		if (strings.Contains(fileContents, "kind: Namespace") || strings.Contains(fileContents, "kind: \"Namespace\"")) {
			prefix = fmt.Sprintf(prefix, "-01-")
			count_ns += 1
		} else if (strings.Contains(fileContents, "kind: Deployment") || strings.Contains(fileContents, "kind: \"Deployment\"")) {
			prefix = fmt.Sprintf(prefix, "-03-")
			count_deploy += 1
		} else {
			prefix = fmt.Sprintf(prefix, "-02-")
			count_crd += 1
		}
		f, err := ioutil.TempFile(path, prefix)
		if err != nil {
			return "", 0, 0, 0, 0, fmt.Errorf("cannot create temporary file: %v", err)
		}

		if _, err = f.Write([]byte(r)); err != nil {
			return "", 0, 0, 0, 0, fmt.Errorf("failed to write temporary file %s: %v", f.Name(), err)
		}

		count += 1
		f.Close()
	}
	return
}
