package nirmata

import (
	"encoding/json"
	"fmt"
	"log"

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

	credentialID := ""
	if git_credentials != "" {
		credID, err := apiClient.QueryByName(client.ServiceEnvironments, "GitCredential", git_credentials)
		if err != nil {
			fmt.Printf("Error - %v", err)
			return err
		}
		credentialID = credID.UUID()
	}
	catID, err := apiClient.QueryByName(client.ServiceCatalogs, "Catalogs", catalog)

	appData := map[string]interface{}{
		"name":         name,
		"parent":       catID.Map(),
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

	data, marshalErr := json.Marshal(appData)
	if marshalErr != nil {
		fmt.Printf("Error - %v", err)
		return err
	}

	log.Printf("[DEBUG] - creating  application %s with %+v", name, appData)
	appId, err := apiClient.PostWithID(catID, "applications", data, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create  application %s with data %v: %v", name, appData, err)
		return err
	}
	catalogAppUUID := appId["id"].(string)
	d.SetId(catalogAppUUID)
	log.Printf("[INFO] - created application %s %s", name, catalogAppUUID)

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
	"git_repository":     "repository",
	"git_directory_list": "directoryList",
	"git_branch":         "branch",
	"git_include_list":   "includeList",
}

func resourceGitApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	id := client.NewID(client.ServiceCatalogs, "applications", d.Id())
	appChanges := buildChanges(d, catalogAppMap, "git_repository", "git_directory_list", "git_branch", "git_include_list")
	if len(appChanges) > 0 {
		gitUpstream, err := apiClient.GetDescendant(id, "GitUpstream", &client.GetOptions{})
		if err != nil {
			log.Printf("[ERROR] - failed to retrieve %s from %v: %v", "GitUpstream", id.Map(), err)
			return err
		}

		d, plainErr := client.NewObject(gitUpstream)
		if plainErr != nil {
			log.Printf("[ERROR] - failed to decode %s %v: %v", "GitUpstream", d, err)
			return err
		}

		_, plainErr = apiClient.PutWithIDFromJSON(d.ID(), appChanges)
		if plainErr != nil {
			log.Printf("[ERROR] - failed to update %s %v: %v", "GitUpstream", d.ID().Map(), err)
			return err
		}

		log.Printf("[DEBUG] updated %v %v", d.ID().Map(), appChanges)
		return nil
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
