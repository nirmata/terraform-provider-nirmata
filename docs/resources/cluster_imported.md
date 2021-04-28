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
}

```

## Argument Reference

* `name` - (Required) a unique name for the cluster.
* `credentials` - (Required) the cloud credentials to use to locate and import the cluster.
* `cluster_type` - (Required) the cluster type to apply.
* `region` - (Required) the region the cluster is located in.
* `project` - (Required) the project the cluster is located in.
* `delete_action` - (Optional) whether to delete or remove the cluster on destroy. Defaults to `remove`.






