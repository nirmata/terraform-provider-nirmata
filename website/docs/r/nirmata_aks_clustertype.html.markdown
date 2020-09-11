---
subcategory: "ClusterType AKS"
layout: "nirmata"
page_title: "Nirmata: nirmata_aks_clusterType"
description: |-
  
---

# Resource: nirmata_aks_clusterType

Manages a AKS ClusterType 


## Example Usage

```hcl
resource "nirmata_aks_clusterType" "aks-cluster-type" {
  name = "aks-tf-1"
  version = "1.17.7"
  credentials = ""
  region = "centralus"
  resource_group = ""
  subnet_id = ""
  vm_size = "Standard_D2_v3"
  vm_set_type = "VirtualMachineScaleSets"
  disk_size = 60
  https_application_routing = false
  monitoring = false
  workspace_id = ""
}

// A Cluster is created using a ClusterType
resource "nirmata_ProviderManaged_cluster" "aks-cluster" {
  name = "tf-akscluster"
  type_selector = nirmata_aks_clusterType.aks-cluster-type.name
  node_count = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) 
* `version` - (Required) 
* `credentials` - (Optional) 
* `region` - (Optional) 
* `resource_group` - (Optional) 
* `subnet_id` - (Optional)
* `vm_size` - (Optional)
* `vm_set_type` - (Optional) 
* `disk_size` - (Optional) 
* `https_application_routing` - (Optional) 
* `monitoring` - (Optional) 
* `workspace_id` - (Optional) 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:


## Timeouts


## Import


```
$ terraform import nirmata_aks_clusterType.aks-cluster-type aks-cluster-type
```
