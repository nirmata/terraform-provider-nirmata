---
page_title: "nirmata_cluster_type_gke Resource"
---

# nirmata_cluster_type_gke Resource

Represents a Google Kubernetes Engine (GKE) cluster type.

## Example Usage

Create a GKE cluster type 

```hcl
resource "nirmata_cluster_type_gke" "gke-us-west" {
  name = "gke-us-west"
  version = "1.16.13-gke.1"
  credentials = "gcp"
  location_type =  "Zonal"
  region = "us-central1"
  zone = "us-central1-a"
  node_locations = []
  enable_secrets_encryption = false
  enable_workload_identity = false
  workload_pool = ""
  secrets_encryption_key = ""
  network = "default"
  subnetwork = "default"
  machine_type = "e2-standard-2"
  disk_size = 60
}
```

## Argument Reference

* `name` - (Required) a unique name for the cluster.
* `version` - (Required) the GKE version (e.g. 1.16.12-gke.3)
* `credentials` - (Required) the cloud credentials to use.
* `location_type` - (Required) Regional or Zonal. A regional cluster has multiple replicas of the control plane, running in multiple zones within a given region. A zonal cluster has a single replica of the control plane running in a single zone.
* `region` - (Optional) the GCP region into which the cluster should be deployed (e.g. "us-central1"). Required when location_type is `Regional`.
* `zone` - (Optional) the GCP zone into which the cluster should be deployed (e.g. "us-central1-a"). Required if location_type is `Zonal`.
* `node_locations` - (Optional) the nodes should be deployed. Selecting more than one zone increases availability. (e.g. ["asia-east1-a"]). Required if location_type is `Regional`.
* `enable_secrets_encryption` - (Optional) enables envelope encryption for Kubernetes Secrets.
* `enable_workload_identity` - (Optional) Workload Identity is the recommended way to access Google Cloud services from applications running within GKE due to its improved security properties and manageability.
* `workload_pool` - (Optional) Workload Identity relies on a Workload Pool to aggregate identity across multiple clusters. Required if enable_secrets_encryption is true
* `secrets_encryption_key` - (Optional) the Resource ID of the key you want to use (e.g. projects/project-name/locations/global/keyRings/my-keyring/cryptoKeys/my-key). equired if enable_workload_identity is true.
* `network` - (Required) the cluster network (e.g. "default")
* `subnetwork` - (Required) the node subnetwork (e.g. "default")
* `machine_type` - (Required) the GCP machine type (e.g. "e2-standard-2")
* `disk_size` - (Required) the worker node disk size in GB (e.g. 60)
