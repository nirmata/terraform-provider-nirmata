---
page_title: "nirmata_eks_cluster_imported Resource"
---

# nirmata_cluster_imported Resource

An existing cloud provider-managed cluster that is discovered and imported using cloud provider credentials.

## Example Usage

Import an existing EKS cluster called `eks-test` in `us-west-1`:

```hcl

resource "nirmata_cluster_imported" "eks-import" {
  name = "my-cluster-1"
  credentials = "eks-test-credentials"
  cluster_type  =  "eks-test"
  region = "us-west-1"
  delete_action = "remove"
  cluster = "eks"
  system_metadata = {
    cluster = "import"
  }
  labels = {foo = "bar"}
  endpoint = "kubernetes cluster API server url"
}

```

## Argument Reference

* `name` - (Required) Enter a unique name for the cluster.
* `credentials` - (Required) Enter the cloud credentials that is used to locate and import the cluster.
* `cluster_type` - (Required) Enter the cluster type to apply for the cluster.
* `cluster` - (Required) Enter the type of cluster (example, gke/eks).
* `region` - (Required) Enter the region of the cluster that is located in.
* `delete_action` - (Optional) This field indicates whether to delete or remove the cluster on destroy. The default value is `remove`.
* `system_metadata` - (Optional) This field indicates the key-value pairs that will be provisioned as a ConfigMap called system-metadata-map in the cluster.
* `labels` - (Optional) This field indicates the labels to be set on the cluster.
* `endpoint` - (Optional) This field indicates the url of the kubernetes cluster API server.
* `owner_info` - (Optional) The [owner_info](#owner_info) for this cluster, if it has to be overridden.
* `access_control_list` - (Optional) List of additional [ACLs](#access_control_list) for this cluster.

### owner_info
* `owner_type` - (Required) The type of the owner. Valid values are user or team.
* `owner_name` - (Required) The name of the user/team.

### access_control_list
* `entity_type` - (Required) The type of entity. Valid values are user or team.
* `permission` - (Required) The permission. Valid values are admin, edit, view.
* `name` - (Required) The name of the user/team.
