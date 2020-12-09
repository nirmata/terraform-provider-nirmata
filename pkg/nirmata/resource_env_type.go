package nirmata

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
	"log"
)

func resourceEnvironmentType() *schema.Resource {
	return &schema.Resource{

		Create: resourceEnvironmentTypeCreate,
		Read:   resourceEnvironmentTypeRead,
		Update: resourceEnvironmentTypeUpdate,
		Delete: resourceEnvironmentTypeDelete,
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
				Required: true,
			},
			"resource_limits": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceEnvironmentTypeCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	isDefault := d.Get("is_default").(bool)
	resourceLimits := d.Get("resource_limits")

	data := map[string]interface{}{
		"name":         name,
		"resourceLimits": resourceLimits,
		"isDefault":  isDefault,
	}

	log.Printf("[DEBUG] - creating environment type %s with %+v", name, data)
	envTypeData, err := apiClient.PostFromJSON(client.ServiceEnvironments, "EnvironmentResourceType", data, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create environment type %s with data %v: %v", name, data, err)
		return err
	}
	envTypeUUID := envTypeData["id"].(string)
	d.SetId(envTypeUUID)
	log.Printf("[INFO] - created environment type %s %s", name, envTypeUUID)

	return nil
}

func resourceEnvironmentTypeRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceEnvironments, "EnvironmentResourceType", d.Id())

	_, err := apiClient.Get(id, &client.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] failed to retrieve environment type details %s (%s): %v", name, id, err)
		return err
	}

	log.Printf("[INFO] - retrieved environment type%s %s", name, id.UUID())
	return nil
}

var envTypeMap = map[string]string{
	"is_default":      "isDefault",
	"resource_limits": "resourceLimits" ,
	}

func resourceEnvironmentTypeUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	id := client.NewID(client.ServiceEnvironments, "EnvironmentResourceType", d.Id())
	envChanges := buildChanges(d, envTypeMap, "is_default","resource_limits")
	if len(envChanges) > 0 {
		envTypeData, err :=  apiClient.Get(id, &client.GetOptions{})
		if err != nil {
			log.Printf("[ERROR] - failed to retrieve %s from %v: %v", "EnvironmentResourceType", id.Map(), err)
			return err
		}
		d, plainErr := client.NewObject(envTypeData)
		if plainErr != nil {
			log.Printf("[ERROR] - failed to decode %s %v: %v", "EnvironmentResourceType", d, err)
			return err
		}
		_, plainErr = apiClient.PutWithIDFromJSON(d.ID(), envChanges)
		if plainErr != nil {
			log.Printf("[ERROR] - failed to update %s %v: %v", "EnvironmentResourceType", d.ID().Map(), err)
			return err
		}
		log.Printf("[DEBUG] updated %v %v", d.ID().Map(), envChanges)
	}
	return nil
}

func resourceEnvironmentTypeDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceEnvironments, "EnvironmentResourceType", d.Id())

	if err := apiClient.Delete(id, nil); err != nil {
		return err
	}

	log.Printf("[INFO] - deleted environment type %s %s", name, id.UUID())
	return nil
}
