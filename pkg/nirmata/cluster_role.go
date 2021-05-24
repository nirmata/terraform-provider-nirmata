package nirmata

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

var clusterRoleSchema = map[string]*schema.Schema{
	"api_groups": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional: true,
	},
	"resources": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional: true,
	},
	"verbs": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional: true,
	},
}
