// NOTE: this example needs to be applied in two phases, as the YAML file count
// is computed during the apply phase of the nirmata_cluster_registered resource.
//
// Steps:
//   terraform init
//   terraform plan 
//   terraform apply -target nirmata_cluster_registered.eks-registered
//   terraform plan
//   terraform apply

provider "nirmata" {
  # Nirmata address.
  url = "https://nirmata.io"
  // Nirmata API Key.  
  token =""

}


// create a new cluster and download the controller YAMLs
resource "nirmata_cluster_registered" "eks-registered" {
  name         = "eks-cluster"
  cluster_type = "default-add-ons"
}


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