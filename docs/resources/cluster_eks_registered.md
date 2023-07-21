---
page_title: "nirmata_eks_cluster_registered Resource"
---

# nirmata_cluster_registered Resource

An existing eks cluster that is registered using the local Kubernetes credentials.

## Example Usage

Register an existing EKS cluster using Kubernetes credentials. The new cluster is created and the controller YAMLs are downloaded to a temporary folder. The YAMLs are then applied to the existing eks cluster using the `kubectl` and `aws` providers.

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
  endpoint     = "kubernetes cluster API server url"
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

data "kubectl_filename_list" "namespace" {
   pattern = "${nirmata_cluster_registered.eks-registered.controller_yamls_folder}/temp-01-*"
}

data "kubectl_filename_list" "sa" {
   pattern = "${nirmata_cluster_registered.eks-registered.controller_yamls_folder}/temp-02-*"
}

data "kubectl_filename_list" "crd" {
   pattern = "${nirmata_cluster_registered.eks-registered.controller_yamls_folder}/temp-03-*"
}

data "kubectl_filename_list" "deployment" {
   pattern = "${nirmata_cluster_registered.eks-registered.controller_yamls_folder}/temp-04-*"
}


// Register Nirmata Cluster
resource "nirmata_cluster_registered" "eks-registered" {
  name         = var.nirmata_cluster_name
  cluster_type = var.nirmata_cluster_type
}

// apply the controller YAMLs
resource "kubectl_manifest" "namespace" {
  wait        = true
  count       = nirmata_cluster_registered.eks-registered.controller_ns_yamls_count
  yaml_body   = file(element(data.kubectl_filename_list.namespace.matches, count.index))
  apply_only  = true
  depends_on  = [nirmata_cluster_registered.eks-registered]
}

resource "kubectl_manifest" "sa" {
  wait        = true
  count       = nirmata_cluster_registered.eks-registered.controller_sa_yamls_count
  yaml_body   = file(element(data.kubectl_filename_list.sa.matches, count.index))
  apply_only  = true
  depends_on  = [kubectl_manifest.namespace]
}

resource "kubectl_manifest" "crd" {
  wait        = true
  count       = nirmata_cluster_registered.eks-registered.controller_crd_yamls_count
  yaml_body   = file(element(data.kubectl_filename_list.crd.matches, count.index))
  apply_only  = true
  depends_on  = [kubectl_manifest.sa]
}

resource "kubectl_manifest" "deployment" {
  wait        = true
  count       = nirmata_cluster_registered.eks-registered.controller_deploy_yamls_count
  yaml_body   = file(element(data.kubectl_filename_list.deployment.matches, count.index))
  apply_only  = true
  depends_on  = [kubectl_manifest.crd]
}



```


## Argument Reference

* `name` - (Required) Enter a unique name for the cluster.
* `cluster_type` - (Required) Enter the cluster type to be applied for the cluster.
* `labels` - (Optional) This field indicates the labels to be set on the cluster.
* `delete_action` - (Optional) This field indicates whether to delete or remove the cluster on destroy. The default value is `remove`.
* `endpoint` - (Optional) This field indicates the url of the kubernetes cluster API server.
* `owner_info` - (Optional) The [owner_info](#owner_info) for this cluster, if it has to be overridden.
* `access_control_list` - (Optional) List of additional [ACLs](#access_control_list) for this cluster.
* `controller_yamls_folder` - (Optional) Location of folder where the controller files will be saved. default is a folder in `/tmp/` with prefix `controller-`

### owner_info
* `owner_type` - (Required) The type of the owner. Valid values are user or team.
* `owner_name` - (Required) The email of the user or the name of the team.

### access_control_list
* `entity_type` - (Required) The type of entity. Valid values are user or team.
* `permission` - (Required) The permission. Valid values are admin, edit, view.
* `name` - (Required) The email of the user or the name of the team.
