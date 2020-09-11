---
subcategory: "ClusterType EKS"
layout: "nirmata"
page_title: "Nirmata: nirmata_eks_clusterType"
description: |-
  
---

# Resource: nirmata_eks_clusterType

Manages a EKS ClusterType 


## Example Usage

```hcl
resource "nirmata_eks_clusterType" "eks-cluster-type" {
  name = "eks_us-west-2_1.16"
  version = "1.16"
  credentials = ""
  region = "us-west-2"
  vpc_id= ""
  subnet_id= []
  security_groups= []
  cluster_role_arn= ""
  key_name= ""
  instance_types = ["t3.medium"]
  disk_size = 60
  node_security_groups = []
  node_iam_role = ""
}


resource "nirmata_ProviderManaged_cluster" "eks-cluster" {
  name = "eks-cluster-1"
  type_selector = nirmata_eks_clusterType.eks-cluster-type.name
  node_count = 1
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) 
* `version` - (Required) 
* `credentials` - (Optional) 
* `region` - (Optional) 
* `vpc_id` - (Optional) 
* `subnet_id` - (Optional)
* `security_groups` - (Optional)
* `cluster_role_arn` - (Optional) 
* `key_name` - (Optional) 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:


## Timeouts


## Import


```
$ terraform import nirmata_eks_clusterType.eks-cluster-type eks-cluster-type
```
