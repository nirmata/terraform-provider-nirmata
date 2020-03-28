package nirmata

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNirmataCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceNirmataClusterCreate,
		Read:   resourceNirmataClusterRead,
		Update: resourceNirmataClusterUpdate,
		Delete: resourceNirmataClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			/*
				The UniqueID could be used as the Id(), but none of the API
				calls allow specifying a user by the UniqueID: they require the
				name. The only way to locate a user by UniqueID is to list them
				all and that would make this provider unnecessarily complex
				and inefficient. Still, there are other reasons one might want
				the UniqueID, so we can make it available.
			*/
			"unique_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: false,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) > 64 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 64 characters", k))
					}
					if !regexp.MustCompile(`^[\w+=,.@-]*$`).MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must match [\\w+=,.@-]", k))
					}
					return
				},
			},
			"cloud_provider": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "gke" {
						errors = append(errors, fmt.Errorf(
							"Currently we only support GKE", k))
					}
					return
				},
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					_ = v.(string)
					//TODO : ADD Logic
					return
				},
			},
			"node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					_ = v.(string)
					//TODO : ADD Logic
					return
				},
			},
			"region": {
				Type:     schema.TypeString,
				Optional: false,
			},
		},
	}
}

func resourceNirmataClusterCreate(d *schema.ResourceData, meta interface{}) error {
	
	return nil
}

func resourceNirmataClusterRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceNirmataClusterUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceNirmataClusterDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}
