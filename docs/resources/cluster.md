---
page_title: "nirmata_cluster Resource"
---

# nirmata_cluster Resource

Represents a cluster. 

## Example Usage

Create a cluster using an available cluster_type

```hcl
resource "nirmata_cluster" "eks-eu-1" {
  name = "eks-eu-1"
  cluster_type = "eks-eu-prod"
  labels  = {foo = "bar"}
   nodepools {
      node_count                = 1 
      enable_auto_scaling       = false
      min_count                 = 1
      max_count                 = 4
   }
   delete_action = "remove"
   creation_timeout_minutes = 30
}

```

## Argument Reference

* `name` - (Required) a unique name for the cluster.
* `cluster_type` - (Required) the type of cluster to create.
* `nodepools` - A list of [nodepool](#nodepool) types.
* `labels` - (Optional) labels to set on cluster.
* `delete_action` - (Optional) if delete_action set to `remove`, cluster only get removed from the Nirmata not from the original provider and delete_action set to `delete` cluster deleted from nirmata as well as original provider.
*`creation_timeout_minutes` - (Optional) set maximum time to create cluster.

## Nested Blocks

### nodepool

* `node_count` - (Required) the number of worker nodes for the cluster
* `enable_auto_scaling` - (Optional) Enable autoscaling for cluster. default valie is disable.
* `min_count` - (Optional) Set minimun node count value for cluster. `enable_auto_scaling` must set true to set min_count.
* `max_count` - (Optional) Set max node count value for cluster. `enable_auto_scaling` must set true to set max_count.

