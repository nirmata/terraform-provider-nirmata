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

* `name` - (Required) unique name for the cluster add-on service.
* `cluster` - (Required) the kubernetes cluster.
* `catalog` - (Required) the catalog.
* `application` - (Required) the application.
* `namespace` - (Optional) defaults to the application name.
* `environment` - (Optional) defaults to the application name and cluster name.
* `channel` - (Required) the release channel
* `labels` - (Optional) labels to set on  cluster add-on.
