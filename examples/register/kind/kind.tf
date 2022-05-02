// NOTE: this example needs to be applied in two phases, as the YAML file count
// is computed during the apply phase of the nirmata_cluster_registered resource.
//
// Steps:
//   terraform init
//   terraform plan 
//   terraform apply -target nirmata_cluster_registered.kind-registered
//   terraform plan
//   terraform apply

// create a new cluster and download the controller YAMLs
resource "nirmata_cluster_registered" "kind-registered" {
  name         = "kind-cluster"
  cluster_type = "default-add-ons"
}

variable "host" {
  type = string
}
variable "client_certificate" {
  type = string
}
variable "client_key" {
  type = string
}
variable "cluster_ca_certificate" {
  type = string
}

provider "kubectl" {
  host = var.host
  client_certificate     = base64decode(var.client_certificate)
  client_key             = base64decode(var.client_key)
  cluster_ca_certificate = base64decode(var.cluster_ca_certificate)
}

data "kubectl_filename_list" "manifests" {
  pattern = "${nirmata_cluster_registered.kind-registered.controller_yamls_folder}/*"
}
// apply the controller YAMLs
resource "kubectl_manifest" "test" {
  count     = nirmata_cluster_registered.kind-registered.controller_yamls_count
  yaml_body = file(element(data.kubectl_filename_list.manifests.matches, count.index))
}