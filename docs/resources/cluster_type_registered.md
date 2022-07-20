---
page_title: "nirmata_cluster_type_registered Resource"
---

# nirmata_cluster_type_registered Resource

A cluster type used to configure add-on services for existing clusters that are registered with Nirmata.

## Example Usage

```hcl

resource "nirmata_cluster_type_registered" "tf-registered-type-1" {
  name  = "tf-registered-cluster"
  cloud = "Other"
  system_metadata = {
    clustertype = "registered"
  }
  
  addons {
    name            = "vault-agent-injector"
    addon_selector  = "vault-agent-injector"
    catalog         = "default-catalog"
    channel         = "Stable"
    sequence_number = 15
  }

  vault_auth {
    name             = "vault-auth"
    path             = "nirmata/$(cluster.name)"
    addon_name       = "vault-agent-injector"
    credentials_name = "vault_access"
    delete_auth_path = false

    roles {
      name                 = "sample-role"
      service_account_name = "application-sample-sa"
      namespace            = "application-sample-ns"
      policies             = "application-sample-policy"
    }
  }
}

```

## Argument Reference

* `name` - (Required) Enter a unique name for the cluster.
* `cloud` - (Optional) Enter the cloud provider. Defaults to `Other`.
* `addons` - (Optional) Enter a list of add-on services.
* `vault_auth` - (Optional) Enter the vault authentication configuration.
* `system_metadata` - (Optional) The key-value pairs that will be provisioned as a ConfigMap called system-metadata-map in the cluster type.

## Nested Blocks

### addons

* `name` - (Required) Enter a unique name for the add-on service.
* `addon_selector` - (Required) Enter the catalog application name.
* `catalog` - (Required) Enter the catalog name.
* `channel` - (Required) Enter the channel from which the application should be deployed.
* `sequence_number` - (Optional) This field indicates a sequence number to control the installation order.

### vault_auth

* `name` - (Required) Enter a unique name.
* `path` - (Required) Enter the vault authentication path. The variable $(cluster.name) is allowed in the path for uniquenes.
* `addon_name` - (Required) Enter the associated Vault Agent Injector add-on.
* `credentials_name` - (Required) Enter the Vault credentials to be used. 
* `roles` - (Required) Enter a list of application roles to configure for add-on services.
* `delete_auth_path` - (Optional) This field indicates the delete authentication path on cluster delete.

#### roles

* `name` - (Required) Enter a unique name
* `service_account_name` - (Required) Enter the allowed service account name.
* `namespace` - (Required) Ener the allowed namespace.
* `policies` - (Required) Enter the applied policies.
