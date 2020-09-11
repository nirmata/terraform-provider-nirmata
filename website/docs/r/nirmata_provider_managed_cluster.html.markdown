---
subcategory: "ProviderManaged Cluster"
layout: "nirmata"
page_title: "Nirmata: nirmata_ProviderManaged_cluster"
description: |-
  
---

# Resource: nirmata_ProviderManaged_cluster

Manages a Provider Managed Cluster


## Example Usage

```hcl
resource "nirmata_ProviderManaged_cluster" "gke-cluster" {
  name = "gke-cluster-1"
  type_selector = nirmata_gke_clusterType.gke-cluster-type.name
  node_count = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) 
* `type_selector` - (Required) 
* `node_count` - (Optional) 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:


## Timeouts


## Import


```
$ terraform import nirmata_ProviderManaged_cluster.gke-cluster gke-cluster
```
