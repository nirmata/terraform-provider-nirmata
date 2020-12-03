package nirmata

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

var addonSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"addon_selector": {
		Type:     schema.TypeString,
		Required: true,
	},
	"catalog": {
		Type:     schema.TypeString,
		Required: true,
	},
	"channel": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"namespace": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"sequence_number": {
		Type:     schema.TypeInt,
		Required: true,
	},
}

func addOnsSchemaToAddOns(d *schema.ResourceData) map[string]interface{} {
	var addonsSpec []map[string]interface{}
	addonsSpec = append(addonsSpec, map[string]interface{}{
		"modelIndex":     "AddOnSpec",
		"name":           "kyverno",
		"addOnSelector":  "kyverno",
		"catalog":        "default-addon-catalog",
		"sequenceNumber": 1,
	})

	if _, ok := d.GetOk("addons"); ok {
		addons := d.Get("addons").([]interface{})
		for _, addon := range addons {
			element, ok := addon.(map[string]interface{})
			if ok {
				addonsSpec = append(addonsSpec, map[string]interface{}{
					"modelIndex":     "AddOnSpec",
					"name":           element["name"],
					"addOnSelector":  element["addon_selector"],
					"catalog":        element["catalog"],
					"channel":        element["channel"],
					"namespace":      element["namespace"],
					"sequenceNumber": element["sequence_number"],
				},
				)
			}
		}
	}

	return map[string]interface{}{
		"dns":        false,
		"modelIndex": "AddOns",
		"other":      addonsSpec,
	}
}
