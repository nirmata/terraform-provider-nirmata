package nirmata

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resourcePolicySet() *schema.Resource {
	return &schema.Resource{

		Create: resourcePolicySetCreate,
		Read:   resourcePolicySetRead,
		Update: resourcePolicySetUpdate,
		Delete: resourcePolicySetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateName,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"git_credentials": {
				Type:     schema.TypeString,
				Required: true,
			},
			"git_repository": {
				Type:     schema.TypeString,
				Required: true,
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
			"delete_from_cluster": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourcePolicySetCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	is_default := d.Get("is_default").(bool)
	git_credentials := d.Get("git_credentials").(string)
	git_repository := d.Get("git_repository").(string)
	git_directoryList := d.Get("git_directory_list")
	git_branch := d.Get("git_branch").(string)
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

	if fixed_kustomization && target_based_kustomization {
		return fmt.Errorf(" [ERROR] - select only one type of kustomization")
	}
	if fixed_kustomization || target_based_kustomization {
		if kustomization_file_path == "" {
			return fmt.Errorf(" [ERROR] - kustomization file path is required")
		}
	}
	appData := map[string]interface{}{
		"name":      name,
		"type":      "git",
		"isDefault": is_default,
		"gitUpstream": map[string]interface{}{
			"name":          name,
			"branch":        git_branch,
			"repository":    git_repository,
			"directoryList": git_directoryList,
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

	policyGroupData, err := apiClient.PostFromJSON(client.ServiceClusters, "PolicyGroup", appData, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create policy group %s with data %v: %v", name, appData, err)
		return err
	}

	policyGroupUUID := policyGroupData["id"].(string)
	log.Printf("[INFO] - created policy group  %s %s", name, policyGroupUUID)

	policyGroupID := client.NewID(client.ServiceClusters, "PolicyGroup", policyGroupUUID)

	policyGroupStatus, policyGroupStatusErr := apiClient.GetDescendant(policyGroupID, "PolicyGroupStatus", &client.GetOptions{})

	if policyGroupStatusErr != nil {
		log.Printf("[ERROR] - failed to create policy group with data : %v", policyGroupStatusErr)
		return policyGroupStatusErr
	}
	policyGroupStatusID := client.NewID(client.ServiceClusters, "PolicyGroupStatus", policyGroupStatus["id"].(string))

	state, waitErr := waitForPolicySetSyncStatus(apiClient, d.Timeout(schema.TimeoutCreate), policyGroupStatusID)
	if waitErr != nil {
		log.Printf("[ERROR] - failed to sync policy group. Error - %v", waitErr)
		return waitErr
	}

	if strings.EqualFold("failed", state) {
		status, err := getPolicyGroupStatus(apiClient, policyGroupStatusID)
		if err != nil {
			log.Printf("[ERROR] - failed to retrieve policy group sync details: %v", err)
			return fmt.Errorf(" [ERROR] - policy group sync failed")
		}
		return fmt.Errorf(" [ERROR] - policy group sync failed: %s", status)
	}

	d.SetId(policyGroupUUID)
	d.Set("policyGroupID", policyGroupUUID)
	log.Printf("[INFO] - created policy group %s %s", name, policyGroupUUID)

	return nil
}

func resourcePolicySetRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceClusters, "PolicyGroup", d.Id())

	_, err := apiClient.Get(id, &client.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] policy group %+v not found", id.Map())
			d.SetId("")
			return nil
		}

		log.Printf("[ERROR] failed to retrieve policy group details %s (%s): %v", name, id, err)
		return err
	}

	log.Printf("[INFO] - retrieved policy group %s %s", name, id.UUID())
	return nil
}

func resourcePolicySetUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourcePolicySetDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	deleteFromCluster := d.Get("delete_from_cluster").(bool)
	id := client.NewID(client.ServiceClusters, "PolicyGroup", d.Id())

	params := map[string]string{
		"deleteFromCluster": strconv.FormatBool(deleteFromCluster),
	}
	if err := apiClient.Delete(id, params); err != nil {
		return err
	}

	log.Printf("[INFO] - deleted policy group %s %s", name, id.UUID())
	return nil
}
