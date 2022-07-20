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

* `name` - (Required) Enter a unique name for the cluster.
* `cluster_type` - (Required) Enter the type of cluster to create.
* `nodepools` - Indicates a list of [nodepool](#nodepool) types.
* `labels` - (Optional) This field indicates the labels to be set on the cluster.
* `delete_action` - (Optional) This field indicates that if delete_action is set to `remove`, then the cluster gets removed from the Nirmata platform and not from the original provider. If delete_action is set to `delete`, then the cluster gets deleted from the Nirmata platform as well as from the original provider.
* `creation_timeout_minutes` - (Optional) This field is set to maximum time to create a cluster.

## Nested Blocks

### nodepool

* `node_count` - (Required) Enter the number of worker nodes for the cluster
* `enable_auto_scaling` - (Optional) This field indicates to enable autoscaling for the cluster. The default value is "disable".
* `min_count` - (Optional) This field indicates to set the minimum node count value for the cluster. To set this value, you must set  `enable_auto_scaling` to true.
* `max_count` - (Optional) This field indicates the set max node count value for the cluster. To set this value, you must set `enable_auto_scaling` to true.

