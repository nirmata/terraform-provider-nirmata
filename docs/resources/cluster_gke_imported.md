---
page_title: "nirmata_gke_cluster_imported Resource"
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
  cluster = "gke"
  delete_action = "remove"
  system_metadata = {
    cluster = "import"
  }
  labels = {foo = "bar"}
  endpoint = "kubernetes cluster API server url"

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

* `name` - (Required) Enter a unique name for the cluster.
* `credentials` - (Required) Enter the cloud credentials to locate and import the cluster.
* `cluster_type` - (Required) Enter the cluster type to apply.
* `region` - (Required) Enter the region where the cluster is located.
* `cluster` - (Required) Enter the type of cluster (gke/eks).
* `project` - (Required) Enter the project where the cluster is located.
* `delete_action` - (Optional) This field indicates whether to delete or remove the cluster on destroy. The default value is set to `remove`.
* `system_metadata` - (Optional) This field indicates the key-value pairs that will be provisioned as a ConfigMap called system-metadata-map in the cluster.
* `labels` - (Optional) This field indicates the labels set on cluster.
* `endpoint` - (Optional) This field indicates the url of the kubernetes cluster API server.
* `owner_info` - (Optional) The [owner_info](#owner_info) for this cluster, if it has to be overridden.
* `access_control_list` - (Optional) List of additional [ACLs](#access_control_list) for this cluster.
* `controller_yamls_folder` - (Optional) Location of folder where the controller files will be saved. default is a folder in `/tmp/` with prefix `controller-`


### owner_info
* `owner_type` - (Required) The type of the owner. Valid values are user or team.
* `owner_name` - (Required) The email of the user or the name of the team.

### access_control_list
* `entity_type` - (Required) The type of entity. Valid values are user or team.
* `permission` - (Required) The permission. Valid values are admin, edit, view.
* `name` - (Required) The email of the user or the name of the team.

### vault_auth

* `name` - (Required) Enter a unique name for vault authentication. 
* `path` - (Required) Enter the vault authentication path. The variable $(cluster.name) is allowed in the path for uniquenes.
* `addon_name` - (Required) Enter the associated Vault Agent Injector add-on.
* `credentials_name` - (Required) Enter the Vault credentials to be used. 
* `roles` - (Required) Enter a list of application roles to configure for the add-on services.
* `delete_auth_path` - (Optional) This field indicates the delete authentication path on cluster delete.

#### roles

* `name` - (Required) Enter a unique name for roles.
* `service_account_name` - (Required) Enter the allowed service account name.
* `namespace` - (Required) Enter the allowed namespace.
* `policies` - (Required) Enter the applied policies.






