---
subcategory: "ClusterType OKE"
layout: "nirmata"
page_title: "Nirmata: nirmata_oke_clusterType"
description: |-
  
---

# Resource: nirmata_oke_clusterType

Manages a OKE ClusterType 


## Example Usage

```hcl
resource "nirmata_oke_clusterType" "cluster-type-oke" {

  name = "oke-cluster-type"
  version = ""
  credentials = ""
  region = " "
  vm_shape = ""
}


```

## Argument Reference

The following arguments are supported:

* `name` - (Required) 
* `version` - (Required) 
* `credentials` - (Optional) 
* `region` - (Optional) 
* `vm_shape` - (Optional) 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:


## Timeouts


## Import


```
$ terraform import nirmata_oke_clusterType.cluster-type-oke cluster-type-oke
```
