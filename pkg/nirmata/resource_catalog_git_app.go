package nirmata

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resourceGitApplication() *schema.Resource {
	return &schema.Resource{

		Create: resourceGitApplicationCreate,
		Read:   resourceGitApplicationRead,
		Update: resourceGitApplicationUpdate,
		Delete: resourceGitApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateName,
			},
			"catalog": {
				Type:     schema.TypeString,
				Required: true,
			},
			"git_credentials": {
				Type:     schema.TypeString,
				Required: true,
			},
			"git_repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"git_branch": {
				Type:     schema.TypeString,
				Required: true,
			},
			"git_directory_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"git_include_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"release_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_kustomization": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"target_based_kustomization": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"kustomization_file_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGitApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	catalog := d.Get("catalog").(string)
	namespace := d.Get("namespace").(string)
	git_credentials := d.Get("git_credentials").(string)
	git_repository := d.Get("git_repository").(string)
	git_directoryList := d.Get("git_directory_list")
	git_branch := d.Get("git_branch").(string)
	git_includeList := d.Get("git_include_list")
	fixed_kustomization := d.Get("fixed_kustomization").(bool)
	target_based_kustomization := d.Get("target_based_kustomization").(bool)
	kustomization_file_path := d.Get("kustomization_file_path").(string)

	credentialID := ""
	if git_credentials != "" {
		credID, err := apiClient.QueryByName(client.ServiceEnvironments, "GitCredential", git_credentials)
		if err != nil {
			fmt.Printf("Error - %v", err)
			return err
		}
		credentialID = credID.UUID()
	}
	catID, cerr := apiClient.QueryByName(client.ServiceCatalogs, "Catalogs", catalog)
	if cerr != nil {
		log.Printf("[ERROR] - failed to find catalog with name : %v", catalog)
		return cerr
	}
	if namespace == "" {
		namespace = name
	}
	if fixed_kustomization && target_based_kustomization {
		return fmt.Errorf(" [ERROR] - select only one type of kustomization")
	}
	if fixed_kustomization || target_based_kustomization {
		if kustomization_file_path == "" {
			return fmt.Errorf(" [ERROR] - kustomization file path is required")
		}
	}

	appData := map[string]interface{}{
		"name":         name,
		"parent":       catID.UUID(),
		"upstreamType": "git",
		"namespace":    namespace,
		"gitUpstream": map[string]interface{}{
			"name":          name,
			"branch":        git_branch,
			"repository":    git_repository,
			"directoryList": git_directoryList,
			"includeList":   git_includeList,
			"gitCredential": map[string]interface{}{
				"service":    "Catalog",
				"modelIndex": "GitUpstream",
				"id":         credentialID,
			},
		},
	}

	if target_based_kustomization {
		appData["patchConfig"] = map[string]interface{}{
			"overlayFile": kustomization_file_path,
		}
	}

	if fixed_kustomization {
		appData["kustomizeConfig"] = map[string]interface{}{
			"overlayFile": kustomization_file_path,
		}
	}

	data, marshalErr := json.Marshal(appData)
	if marshalErr != nil {
		fmt.Printf("Error - %v", marshalErr)
		return marshalErr
	}
	log.Printf("[DEBUG] - creating  application %s with %+v", name, appData)
	appId, err := apiClient.PostWithID(catID, "applications", data, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create  application %s with data %v: %v", name, appData, err)
		return err
	}
	fields := []string{"version", "name"}
	catalogGitAppUUID := appId["id"].(string)
	gitAppID := client.NewID(client.ServiceCatalogs, "Applications", catalogGitAppUUID)
	gitUpstream, gitUpstreamErr := apiClient.GetDescendant(gitAppID, "GitUpstream", &client.GetOptions{})
	if gitUpstreamErr != nil {
		log.Printf("[ERROR] - failed to create git upstream with data : %v", gitUpstreamErr)
		return gitUpstreamErr
	}
	gitUpstreamID := client.NewID(client.ServiceCatalogs, "GitUpstream", gitUpstream["id"].(string))

	state, waitErr := waitForGitSyncStatus(apiClient, d.Timeout(schema.TimeoutCreate), gitUpstreamID)
	if waitErr != nil {
		log.Printf("[ERROR] - failed to get git sync status. Error - %v", waitErr)
		return waitErr
	}

	if strings.EqualFold("failed", state) {
		status, err := getGitUpstreamStatus(apiClient, gitUpstreamID)
		if err != nil {
			log.Printf("[ERROR] - failed to retrieve git sync  details: %v", err)
			return fmt.Errorf(" [ERROR] - git sync failed")
		}
		return fmt.Errorf(" [ERROR] - git sync failed: %s", status)
	}

	version, versionErr := apiClient.GetDescendant(gitAppID, "Version", &client.GetOptions{Fields: fields})
	if versionErr != nil {
		log.Printf("Error version not found - %v", versionErr)
		return versionErr
	}

	d.SetId(catalogGitAppUUID)
	d.Set("version", version["version"].(string))
	d.Set("release_name", version["name"].(string))
	log.Printf("[INFO] - created application %s %s", name, catalogGitAppUUID)

	return nil
}

func resourceGitApplicationRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceCatalogs, "applications", d.Id())

	_, err := apiClient.Get(id, &client.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] failed to retrieve application detail %s (%s): %v", name, id, err)
		return err
	}

	log.Printf("[INFO] - retrieved application %s %s", name, id.UUID())
	return nil
}

var catalogAppMap = map[string]string{
	"git_repository":             "repository",
	"git_directory_list":         "directoryList",
	"git_branch":                 "branch",
	"git_include_list":           "includeList",
	"fixed_kustomization":        "kustomizeConfig",
	"target_based_kustomization": "patchConfig",
}

func resourceGitApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	id := client.NewID(client.ServiceCatalogs, "applications", d.Id())
	appChanges := buildChanges(d, catalogAppMap, "git_repository", "git_directory_list", "git_branch", "git_include_list")
	fixed_kustomization := buildChanges(d, catalogAppMap, "fixed_kustomization")
	target_based_kustomization := buildChanges(d, catalogAppMap, "target_based_kustomization")
	applicationD, err := apiClient.Get(id, &client.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] failed to retrieve application detail %s : %v", id, err)
		return err
	}
	if len(appChanges) > 0 {
		gitUpstream, err := apiClient.GetDescendant(id, "GitUpstream", &client.GetOptions{})
		if err != nil {
			log.Printf("[ERROR] - failed to retrieve %s from %v: %v", "git upstream", id.Map(), err)
			return err
		}

		d, plainErr := client.NewObject(gitUpstream)
		if plainErr != nil {
			log.Printf("[ERROR] - failed to decode %s %v: %v", "git upstream", d, err)
			return err
		}

		_, plainErr = apiClient.PutWithIDFromJSON(d.ID(), appChanges)
		if plainErr != nil {
			log.Printf("[ERROR] - failed to update %s %v: %v", "git upstream", d.ID().Map(), err)
			return err
		}

		log.Printf("[DEBUG] updated %v %v", d.ID().Map(), appChanges)
	}
	if len(fixed_kustomization) > 0 || len(target_based_kustomization) > 0 {
		txn := make(map[string]interface{})
		modelIndex := "kustomizeConfig"
		if d.Get("fixed_kustomization").(bool) {
			txn["delete"] = applicationD["patchConfig"]
		}

		if d.Get("target_based_kustomization").(bool) {
			txn["delete"] = applicationD["kustomizeConfig"]
			modelIndex = "patchConfig"
		}
		parent := map[string]interface{}{
			"id":            d.Id(),
			"modelIndex":    "Application",
			"service":       "Catalog",
			"childRelation": modelIndex,
		}
		var createdKustomizationObj []map[string]interface{}
		createdKustomizationObj = append(createdKustomizationObj, map[string]interface{}{
			"overlayFile": d.Get("kustomization_file_path").(string),
			"parent":      parent,
			"modelIndex":  modelIndex,
		})
		txn["create"] = createdKustomizationObj
		_, txnErr := apiClient.PostFromJSON(client.ServiceCatalogs, "txn", txn, nil)
		if txnErr != nil {
			log.Printf("[ERROR] - failed to update cluster role with data : %v", txnErr)
			return txnErr
		}
	}

	return nil
}

func resourceGitApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceCatalogs, "Application", d.Id())

	if err := apiClient.Delete(id, nil); err != nil {
		return err
	}

	log.Printf("[INFO] - deleted application %s %s", name, id.UUID())
	return nil
}
