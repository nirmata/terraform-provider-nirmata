---
subcategory: "ClusterType DirectConnect"
layout: "nirmata"
page_title: "Nirmata: nirmata_cluster_direct_connect"
description: |-
  
---

# Resource: nirmata_cluster_direct_connect

Manages a Cluster Direct Connect


## Example Usage

```hcl
resource "nirmata_cluster_direct_connect" "dc-cluster-1" {
  name       = "dc-cluster-1"
  policy     = "default-v1.16.0"
  host_group = nirmata_host_group_direct_connect.dc-host-group.name
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) 
* `policy` - (Required) 
* `host_group` - (Required) 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:


## Timeouts


## Import


```
$ terraform import nirmata_cluster_direct_connect.dc-cluster-1 dc-cluster-1
```
