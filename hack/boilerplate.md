# Resources

```
package nirmata

import (
	
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/nirmata/terraform-provider-nirmaa/pkg/nirmata"
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNirmataClusterName,
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

```

# Datasource

```
package nirmata

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNirmataClusterAlias() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNirmataClusterAliasRead,

		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNirmataCluster(d *schema.ResourceData, meta interface{}) error {

	return nil
}

```