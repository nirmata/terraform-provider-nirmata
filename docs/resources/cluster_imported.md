---
page_title: "nirmata_cluster_imported Resource"
---

# nirmata_cluster_imported Resource

An existing cloud provider managed cluster that is discovered and imported using cloud provider credentials.

## Example Usage

Import an existing GKE cluster called `gke-test` in `us-central1-c`:

```hcl

resource "nirmata_cluster_imported" "gke-import-1" {
  name = "my-cluster-1"
  credentials = "gke-test-credentials"
  cluster_type  =  "gke-test"
  region = "us-central1-c"
  project = "my-project"
  delete_action = "remove"
  system_metadata = {
    cluster = "import"
  }
  labels = {foo = "bar"}

  vault_auth {
    name             = "vault-auth"
    path             = "nirmata/$(cluster.name)"
    addon_name       = "vault-agent-injector"
    credentials_name = "vault_access"
    delete_auth_path = true

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
* `credentials` - (Required) the cloud credentials to use to locate and import the cluster.
* `cluster_type` - (Required) the cluster type to apply.
* `region` - (Required) the region the cluster is located in.
* `project` - (Required) the project the cluster is located in.
* `delete_action` - (Optional) whether to delete or remove the cluster on destroy. Defaults to `remove`.
* `system_metadata` - (Optional) key-value pairs that will be provisioned as a ConfigMap called system-metadata-map in the cluster.
* `labels` - (Optional) labels to set on  cluster.


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






