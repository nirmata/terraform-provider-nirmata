package nirmata

import (
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

var clusterTypeaddOnMap = map[string]string{
	"vault_auth":    "vault",
	"addons":        "addons",
	"cluster_roles": "clusterRoles",
}
var vaultAuthMap = map[string]string{
	"name":             "name",
	"path":             "path",
	"delete_auth_path": "shouldDeleteAuthPath",
	"addon_name":       "addOnName",
	"roles":            "roles",
}
var enableIAMMap = map[string]string{
	"enable_iam_authentication": "enableIAMAuthentication",
	"enable_iam_authorization":  "enableIAMAuthorization",
}

func updateClusterTypeAddonAndVault(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	clusterTypeID := client.NewID(client.ServiceClusters, "ClusterType", d.Id())
	vaultAuthChanges := buildChanges(d, clusterTypeaddOnMap, "vault_auth")
	addonsChanges := buildChanges(d, clusterTypeaddOnMap, "addons")
	clusterRolesChanges := buildChanges(d, clusterTypeaddOnMap, "cluster_roles")
	enableIAMChanges := buildChanges(d, enableIAMMap, "enable_iam_authentication", "enable_iam_authorization")

	if len(addonsChanges) > 0 {
		if err := updateClusterTypeAddon(d, apiClient, clusterTypeID); err != nil {
			log.Printf("[ERROR] - failed to update cluster type add-on with data : %v", err)
			return err
		}
	}
	if len(vaultAuthChanges) > 0 {
		if err := updateVaultAddon(d, apiClient, clusterTypeID); err != nil {
			log.Printf("[ERROR] - failed to update cluster type vault with data : %v", err)
			return err
		}
	}
	if len(enableIAMChanges) > 0 {
		err := updateDescendant(apiClient, clusterTypeID, "IAMSpec", enableIAMChanges)
		if err != nil {
			return err
		}
	}

	if len(clusterRolesChanges) > 0 {
		if err := updateClusterRoles(d, apiClient, clusterTypeID); err != nil {
			log.Printf("[ERROR] - failed to update cluster type vault with data : %v", err)
			return err
		}
	}

	return nil
}

func updateClusterTypeAddon(d *schema.ResourceData, m interface{}, clusterTypeID client.ID) error {
	apiClient := m.(client.Client)
	fields := []string{"name", "id"}
	addOnSpecs, err := apiClient.GetDescendants(clusterTypeID, "addOnSpecs", &client.GetOptions{Fields: fields})
	if err != nil {
		log.Printf("[ERROR] - failed to retrieve %s from %v: %v", "addOnSpecs", clusterTypeID.Map(), err)
		return err
	}
	addOns, err := apiClient.GetDescendant(clusterTypeID, "AddOns", &client.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] - failed to retrieve %s from %v: %v", "addOnSpecs", clusterTypeID.Map(), err)
		return err
	}
	var newAddonNames []interface{}
	var createdAddons []map[string]interface{}
	var updatedAddons []map[string]interface{}
	var deletedAddons []map[string]interface{}
	addons := d.Get("addons").([]interface{})
	existsAddonNames := getAddonNames(addOnSpecs, false)
	for _, newName := range addons {
		a := newName.(map[string]interface{})
		newAddonNames = append(newAddonNames, a["name"])
	}
	createdAddonNames := getExitsAddonNames(newAddonNames, existsAddonNames)
	deletedAddonNames := getExitsAddonNames(existsAddonNames, newAddonNames)
	for _, cname := range createdAddonNames {
		for _, addon := range addons {
			element, ok := addon.(map[string]interface{})
			if ok {
				if cname == element["name"] {
					createdAddons = append(createdAddons, map[string]interface{}{
						"modelIndex":     "AddOnSpec",
						"name":           element["name"],
						"addOnSelector":  element["addon_selector"],
						"catalog":        element["catalog"],
						"channel":        element["channel"],
						"namespace":      element["namespace"],
						"sequenceNumber": element["sequence_number"],
						"parent":         addOns["id"],
					})
					break
				}
			}
		}
	}

	for _, oldAddon := range addOnSpecs {
		for _, addon := range addons {
			element, ok := addon.(map[string]interface{})
			if ok {
				if oldAddon["name"] == element["name"] {
					updatedAddons = append(updatedAddons, map[string]interface{}{
						"modelIndex":     "AddOnSpec",
						"name":           element["name"],
						"addOnSelector":  element["addon_selector"],
						"catalog":        element["catalog"],
						"channel":        element["channel"],
						"namespace":      element["namespace"],
						"sequenceNumber": element["sequence_number"],
						"id":             oldAddon["id"],
					})
					break
				}
			}
		}

		for _, dname := range deletedAddonNames {
			if oldAddon["name"] != "kyverno" {
				if oldAddon["name"] == dname {
					deletedAddons = append(deletedAddons, map[string]interface{}{
						"modelIndex": "AddOnSpec",
						"id":         oldAddon["id"],
					})
					break
				}
			}
		}
	}
	addOnTxn := make(map[string]interface{})
	addOnTxn["update"] = updatedAddons
	addOnTxn["create"] = createdAddons
	addOnTxn["delete"] = deletedAddons

	_, txnErr := apiClient.PostFromJSON(client.ServiceClusters, "txn", addOnTxn, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to update cluster ddon  with data : %v", txnErr)
		return err
	}
	return nil
}

func updateClusterRoles(d *schema.ResourceData, m interface{}, clusterTypeID client.ID) error {
	apiClient := m.(client.Client)
	iamSpecData, err := apiClient.GetDescendant(clusterTypeID, "IAMSpec", &client.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] - failed to retrieve %s from %v: %v", "IAMSpec", clusterTypeID.Map(), err)
		return err
	}

	iamSpec, plainErr := client.NewObject(iamSpecData)
	if plainErr != nil {
		log.Printf("[ERROR] - failed to decode %s %v: %v", iamSpecData, d, err)
		return err
	}

	clusterRoles, err := apiClient.GetRelation(iamSpec.ID(), "clusterRoles")
	if err != nil {
		log.Printf("[ERROR] - failed to retrieve %s from %v: %v", "addOnSpecs", clusterTypeID.Map(), err)
		return err
	}
	parent := map[string]interface{}{
		"id":         clusterRoles["id"],
		"service":    "Cluster",
		"modelIndex": "ClusterRole",
	}
	cRoles := d.Get("cluster_roles").([]interface{})
	var createdAddons []map[string]interface{}
	for _, addon := range cRoles {
		element, ok := addon.(map[string]interface{})
		if ok {

			createdAddons = append(createdAddons, map[string]interface{}{
				"modelIndex": "PolicyRule",
				"apiGroups":  element["api_groups"],
				"resources":  element["resources"],
				"verbs":      element["verbs"],
				"parent":     parent,
			})
		}
	}
	txn := make(map[string]interface{})
	txn["delete"] = clusterRoles["rules"]
	txn["create"] = createdAddons
	_, txnErr := apiClient.PostFromJSON(client.ServiceClusters, "txn", txn, nil)
	if txnErr != nil {
		log.Printf("[ERROR] - failed to update cluster role with data : %v", txnErr)
		return txnErr
	}
	return nil
}

func updateVaultAddon(d *schema.ResourceData, m interface{}, clusterTypeID client.ID) error {
	apiClient := m.(client.Client)
	vaultSpecData, err := apiClient.GetDescendant(clusterTypeID, "VaultKubernetesAuthSpec", &client.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] - failed to retrieve %s from %v: %v", "VaultKubernetesAuthSpec", clusterTypeID.Map(), err)
		return err
	}

	vaultAuth := d.Get("vault_auth").([]interface{})
	vault := vaultAuth[0].(map[string]interface{})
	changes := map[string]interface{}{}
	var rolesSpec []map[string]interface{}

	parent := map[string]interface{}{
		"id":         vaultSpecData["id"],
		"service":    "Cluster",
		"modelIndex": "VaultKubernetesAuthSpec",
	}

	for k, s := range vault {
		if k == "roles" {
			roles := s.([]interface{})
			for _, role := range roles {
				element := role.(map[string]interface{})
				rolesSpec = append(rolesSpec, map[string]interface{}{
					"modelIndex":         "VaultRole",
					"name":               element["name"],
					"serviceAccountName": element["service_account_name"],
					"namespace":          element["namespace"],
					"policies":           element["policies"],
					"parent":             parent,
				})
			}
		}
		name := vaultAuthMap[k]
		changes[name] = s
	}
	txn := make(map[string]interface{})
	txn["delete"] = vaultSpecData["roles"]
	txn["create"] = rolesSpec

	_, txnErr := apiClient.PostFromJSON(client.ServiceClusters, "txn", txn, nil)
	if txnErr != nil {
		log.Printf("[ERROR] - failed to create cluster type  with data : %v", txnErr)
		return txnErr
	}

	vaultErr := updateDescendant(apiClient, clusterTypeID, "VaultKubernetesAuthSpec", changes)
	if vaultErr != nil {
		return vaultErr
	}
	return nil
}

func updateDescendant(apiClient client.Client, id client.ID, descendant string, changes map[string]interface{}) error {

	clusterSpecData, err := apiClient.GetDescendant(id, descendant, &client.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] - failed to retrieve %s from %v: %v", descendant, id.Map(), err)
		return err
	}

	d, plainErr := client.NewObject(clusterSpecData)
	if plainErr != nil {
		log.Printf("[ERROR] - failed to decode %s %v: %v", descendant, d, err)
		return err
	}

	_, plainErr = apiClient.PutWithIDFromJSON(d.ID(), changes)
	if plainErr != nil {
		log.Printf("[ERROR] - failed to update %s %v: %v", descendant, d.ID().Map(), err)
		return err
	}

	log.Printf("[DEBUG] updated %v %v", d.ID().Map(), changes)
	return nil
}

func itemExists(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)
	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

func getExitsAddonNames(n, o []interface{}) (names []interface{}) {
	for _, n := range n {
		if !itemExists(o, n) {
			names = append(names, n)
		}
	}
	return names
}

func buildChanges(d *schema.ResourceData, nameMap map[string]string, attributes ...string) map[string]interface{} {
	changes := map[string]interface{}{}
	for _, a := range attributes {
		if d.HasChange(a) {
			name := nameMap[a]
			changes[name] = d.Get(a)
		}
	}

	return changes
}

func getAddonNames(addon []map[string]interface{}, isInterface bool) (names []interface{}) {
	for _, n := range addon {
		names = append(names, n["name"])
	}
	return names
}
