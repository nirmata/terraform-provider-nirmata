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

* `name` - (Required) a unique name for the cluster.
* `cloud` - (Optional) the cloud provider. Defaults to `Other`.
* `addons` - (Optional) a list of add-on services.
* `vault_auth` - (Optional) vault authentication configuration.
* `system_metadata` - (Optional) key-value pairs that will be provisioned as a ConfigMap called system-metadata-map in the cluster type.

## Nested Blocks

### addons

* `name` - (Required) a unique name for the add-on service
* `addon_selector` - (Required) the catalog application name
* `catalog` - (Required) the catalog name
* `channel` - (Required) the release channel
* `sequence_number` - (Optional) a sequence number to control installation order

### vault_auth

* `name` - (Required) a unique name
* `path` - (Required) the vault authentication path. The variable $(cluster.name) is allowed in the path for uniquenes.
* `addon_name` - (Required) the associated Vault Agent Injector add-on
* `credentials_name` - (Required) the Vault credentials to use 
* `roles` - (Required) a list of application roles to configure for add-on services
* `delete_auth_path` - (Optional) delete auth path on cluster delete

#### roles

* `name` - (Required) a unique name
* `service_account_name` - (Required) the allowed service account name
* `namespace` - (Required) the allowed namespace
* `policies` - (Required) the applied policies
