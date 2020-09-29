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
  node_count = 3
}
```

## Argument Reference

* `name` - (Required) a unique name for the cluster.
* `cluster_type` - (Required) the type of cluster to create.
* `node_count` - (Required) the number of worker nodes for the cluster.


