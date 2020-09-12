---
subcategory: "ClusterType GKE"
layout: "nirmata"
page_title: "Nirmata: nirmata_gke_clusterType"
description: |-
  
---

# Resource: nirmata_gke_clusterType

Manages a EKS ClusterType 


## Example Usage

```hcl
// A ClusterType contains reusable configuration to create clusters.
resource "nirmata_gke_clusterType" "gke-cluster-type" {
  name = "gke-cluster-type"
  version = "1.16.12-gke.3"
  credentials = ""
  region = "us-central1-b"
  machine_type = "e2-standard-2"
  disk_size = 60
}

// A Cluster is created using a ClusterType
resource "nirmata_ProviderManaged_cluster" "gke-cluster" {
  name = "gke-cluster-1"
  type_selector = nirmata_gke_clusterType.gke-cluster-type.name
  node_count = 1
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) 
* `version` - (Required) 
* `credentials` - (Optional) 
* `region` - (Optional) 
* `vpc_id` - (Optional) 
* `machine_type` - (Optional)
* `disk_size` - (Optional)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:


## Timeouts


## Import


```
$ terraform import nirmata_gke_clusterType.gke-cluster-type gke-cluster-type
```
