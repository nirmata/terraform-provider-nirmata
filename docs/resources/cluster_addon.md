---
page_title: "nirmata_cluster_addons Resource"
---

# nirmata_cluster_addons Resource

You can deploy and manage add-ons across clusters.

## Example Usage

```hcl

resource "nirmata_cluster_addons" "cluster_addon" {
  name                       = "addon"
  cluster                    = "cluster-1"
  catalog                    = "default-addon-catalog"
  application                = "app"
  channel                    = "channel"
}

```

## Argument Reference

* `name` - (Required) Enter a unique name for the cluster add-on service.
* `cluster` - (Required) Enter the kubernetes cluster.
* `catalog` - (Required) Enter the catalog name.
* `application` - (Required) Enter the application name.
* `namespace` - (Optional) This field indicates the default value to the application name.
* `environment` - (Optional) This field indicates the defaults to the application name and the cluster name.
* `channel` - (Required) Enter the channel from which the application should be deployed.
* `labels` - (Optional) This field indicates the labels set on cluster add-on.
