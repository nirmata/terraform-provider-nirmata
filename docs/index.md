---
page_title: "Nirmata Provider"
---

# Nirmata Provider

The Nirmata Provider automates Kubernetes cluster and workload management using [Nirmata](https://nirmata.com).

## Example Usage

```hcl
# configure the Nirmata provider
provider "nirmata" {
 
  # Nirmata address.
  url = "https://nirmata.io"

  // Nirmata API Key. Also configurable using the environment variable NIRMATA_TOKEN.
  token = var.nirmata.token

}
```

```hcl
# create a cluster using the Nirmata provider
resource "nirmata_cluster" "eks-eu" {
  name = "eks-eu"
  cluster_type = "eks-eu-prod"
  labels  = {foo = "bar"}
   nodepools {
      node_count                = 1 
      enable_auto_scaling       = false
      min_count                 = 1
      max_count                 = 4
   }
   delete_action = "remove"
   owner_info   = {
      owner_type = "user"
      owner_name = "user_email"
   }
   access_control_list {
      entity_type = "user"
      permission  = "admin"
      name        = "user_email"
   }
   access_control_list {
      entity_type = "team"
      permission  = "edit"
      name        = "team_name"
   }
}
```

## Argument Reference

* `name` - (Required) a unique name for the cluster.
* `cluster_type` - (Required) the type of cluster to create.
* `nodepools` - A list of [nodepool](#nodepool) types.
* `labels` - (Optional) labels to set on cluster.
* `delete_action` - (Optional) if delete_action set to `remove`, cluster only get removed from the Nirmata not from the original provider and delete_action set to `delete` cluster deleted from nirmata as well as original provider.
* `endpoint` - (Optional) This field indicates the url of the kubernetes cluster API server.
* `owner_info` - (Optional) The [owner_info](#owner_info) for this cluster, if it has to be overridden.
* `access_control_list` - (Optional) List of additional [ACLs](#access_control_list) for this cluster.
## Nested Blocks

### nodepool

* `node_count` - (Required) the number of worker nodes for the cluster
* `enable_auto_scaling` - (Optional) Enable autoscaling for cluster. default valie is disable.
* `min_count` - (Optional) Set minimun node count value for cluster. `enable_auto_scaling` must set true to set min_count.
* `max_count` - (Optional) Set max node count value for cluster. `enable_auto_scaling` must set true to set max_count.

### owner_info
* `owner_type` - (Required) The type of the owner. Valid values are user or team.
* `owner_name` - (Required) The email of the user or the name of the team.

### access_control_list
* `entity_type` - (Required) The type of entity. Valid values are user or team.
* `permission` - (Required) The permission. Valid values are admin, edit, view.
* `name` - (Required) The email of the user or the name of the team.
