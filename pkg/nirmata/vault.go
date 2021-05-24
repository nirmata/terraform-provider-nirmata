package nirmata

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

var vaultAuthSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"path": {
		Type:     schema.TypeString,
		Required: true,
	},
	"delete_auth_path": {
		Type:     schema.TypeBool,
		Optional: true,
	},
	"addon_name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"credentials_id": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"credentials_name": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"roles": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: vaultRoleSchema,
		},
	},
}

var vaultRoleSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"service_account_name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"namespace": {
		Type:     schema.TypeString,
		Required: true,
	},
	"policies": {
		Type:     schema.TypeString,
		Required: true,
	},
}

func vaultAuthSchemaToVaultAuthSpec(vaultAuthSchema map[string]interface{}, m interface{}) map[string]interface{} {
	apiClient := m.(client.Client)
	vaultAuthSpec := map[string]interface{}{
		"modelIndex":           "VaultKubernetesAuthSpec",
		"name":                 vaultAuthSchema["name"],
		"path":                 vaultAuthSchema["path"],
		"shouldDeleteAuthPath": vaultAuthSchema["delete_auth_path"],
		"addOnName":            vaultAuthSchema["addon_name"],
	}

	var rolesSpec []map[string]interface{}
	if _, ok := vaultAuthSchema["roles"]; ok {
		roles := vaultAuthSchema["roles"].([]interface{})
		for _, role := range roles {
			element, ok := role.(map[string]interface{})
			if ok {
				rolesSpec = append(rolesSpec, map[string]interface{}{
					"modelIndex":         "VaultRole",
					"name":               element["name"],
					"serviceAccountName": element["service_account_name"],
					"namespace":          element["namespace"],
					"policies":           element["policies"],
				},
				)
			}
		}
		vaultAuthSpec["roles"] = rolesSpec
	}

	credentialSpec := map[string]interface{}{
		"modelIndex": "VaultCredentials",
	}

	if ci, ok := vaultAuthSchema["credentials_id"]; ok {
		credentialSpec["id"] = ci
	}

	if cn, ok := vaultAuthSchema["credentials_name"]; ok {
		name := vaultAuthSchema["credentials_name"].(string)
		vaultID, err := apiClient.QueryByName(client.ServiceClusters, "VaultCredentials", name)
		if err != nil {
			log.Printf("[ERROR] - %v", err)

		}
		credentialSpec["name"] = cn
		credentialSpec["id"] = vaultID.UUID()
	}

	vaultAuthSpec["credentials"] = credentialSpec

	return vaultAuthSpec
}
