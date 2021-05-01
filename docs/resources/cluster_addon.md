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

* `name` - (Required) a unique name for the cluster addon.
* `cluster` - (Required) the host cluster.
* `catalog` - (Required) the catalog.
* `application` - (Required) the application.
* `namespace` - (Optional) Defaults to the application name.
* `environment` - (Optional) Defaults to the application name and cluster name.
* `channel` - (Required) the release channel
* `labels` - (Optional) Labels to set on  cluster addon.
* `service_name` - (Optional).
* `service_port` - (Optional).
* `service_scheme` - (Optional).
