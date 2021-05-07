package nirmata

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resourceRegisteredClusterType() *schema.Resource {
	return &schema.Resource{
		Create: resourceRegisterClusterTypeCreate,
		Read:   resourceRegisterClusterTypeRead,
		Update: resourceRegisterClusterTypeUpdate,
		Delete: resourceRegisterClusterTypeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateName,
			},
			"cloud": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addons": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: addonSchema,
				},
			},
			"vault_auth": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: vaultAuthSchema,
				},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceRegisterClusterTypeCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	cloud := d.Get("cloud").(string)
	labels := d.Get("labels")
	addons := addOnsSchemaToAddOns(d)

	clustertype := map[string]interface{}{
		"name":        name,
		"description": "",
		"modelIndex":  "ClusterType",
		"spec": map[string]interface{}{
			"clusterMode": "discovered",
			"modelIndex":  "ClusterSpec",
			"cloud":       cloud,
			"addons":      addons,
			"labels":      labels,
		},
	}
	if _, ok := d.GetOk("vault_auth"); ok {
		vl := d.Get("vault_auth").([]interface{})
		vault := vl[0].(map[string]interface{})
		clustertype["spec"].(map[string]interface{})["vault"] = vaultAuthSchemaToVaultAuthSpec(vault)
	}

	log.Printf("[DEBUG] - creating register cluster type %s with %+v", name, clustertype)
	data, err := apiClient.PostFromJSON(client.ServiceClusters, "ClusterType", clustertype, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create register cluster type %s with data %v: %v", name, clustertype, err)
		return err
	}

	UUID := data["id"].(string)
	d.SetId(UUID)
	log.Printf("[INFO] - created register cluster type %s %s", name, UUID)

	return nil
}

func resourceRegisterClusterTypeRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceClusters, "ClusterType", d.Id())
	_, err := apiClient.Get(id, &client.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] failed to retrieve cluster type details %s (%s): %v", name, id, err)
		return err
	}

	log.Printf("[INFO] - retrieved cluster type%s %s", name, id.UUID())
	return nil
}

func resourceRegisterClusterTypeUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := updateClusterTypeAddonAndVault(d, meta); err != nil {
		log.Printf("[ERROR] - failed to update cluster type add-on and vault auth with data : %v", err)
		return err
	}
	return nil
}

func resourceRegisterClusterTypeDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceClusters, "ClusterType", d.Id())

	if err := apiClient.Delete(id, nil); err != nil {
		return err
	}

	log.Printf("[INFO] - deleted cluster type %s %s", name, id.UUID())
	return nil
}
