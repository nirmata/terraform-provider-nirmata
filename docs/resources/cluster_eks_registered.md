---
page_title: "nirmata_eks_cluster_registered Resource"
---

# nirmata_cluster_registered Resource

An existing eks cluster that is registered using local Kubernetes credentials.

## Example Usage

Register an existing EKS cluster using Kubernetes credentials. The new cluster is created and the controller YAMLs are downloaded to a temporary folder. The YAMLs are then applied to the existing eks cluster using the `kubectl` providers.

**NOTE:** this example needs to be applied in two phases, as the YAML file count is computed during the apply phase of the nirmata_cluster_registered resource. Steps:
1. terraform init
2. terraform plan 
3. terraform apply -target nirmata_cluster_registered.eks-registered
4. terraform plan
5. terraform apply

```hcl

resource "nirmata_cluster_registered" "eks-registered" {
  name         = "eks-cluster"
  cluster_type = "default-add-ons"
}

# Retrieve eks cluster information
provider "aws" {
  region = "us-west-1"
}

data "aws_eks_cluster" "cluster" {
  name = ""
}

provider "kubectl" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  exec {
    api_version = "client.authentication.k8s.io/v1alpha1"
    command     = "aws"
    args = [
      "eks",
      "get-token",
      "--cluster-name",
      data.aws_eks_cluster.cluster.name
    ]
  }
}

data "kubectl_filename_list" "manifests" {
  pattern = "${nirmata_cluster_registered.eks-registered.controller_yamls_folder}/*"
}

// apply the controller YAMLs
resource "kubectl_manifest" "test" {
  count     = nirmata_cluster_registered.eks-registered.controller_yamls_count
  yaml_body = file(element(data.kubectl_filename_list.manifests.matches, count.index))
}


```


## Argument Reference

* `name` - (Required) a unique name for the cluster.
* `cluster_type` - (Required) the cluster type to apply.
* `controller_yamls` - (Computed) the controller YAML
* `controller_yamls_folder` - (Computed) a local temporary folder with the controller YAML files
* `controller_yamls_count` - (Computed) the controller YAML file count
* `labels` - (Optional) labels to set on cluster.
* `delete_action` - (Optional) whether to delete or remove the cluster on destroy. Defaults to `remove`.




