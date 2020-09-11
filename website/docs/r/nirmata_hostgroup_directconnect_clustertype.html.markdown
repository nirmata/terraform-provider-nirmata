---
subcategory: "ClusterType DirectConnect"
layout: "nirmata"
page_title: "Nirmata: nirmata_host_group_direct_connect"
description: |-
  
---

# Resource: nirmata_host_group_direct_connect

Manages a Host Group ClusterType 


## Example Usage

```hcl
resource "nirmata_host_group_direct_connect" "dc-host-group" {
  name = "dc-hg-1"
}

resource "nirmata_cluster_direct_connect" "dc-cluster-1" {
  name       = "dc-cluster-1"
  policy     = "default-v1.16.0"
  host_group = nirmata_host_group_direct_connect.dc-host-group.name
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:


## Timeouts


## Import


```
$ terraform import nirmata_host_group_direct_connect.dc-host-group dc-host-group
```
