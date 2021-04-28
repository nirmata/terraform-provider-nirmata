resource "nirmata_cluster_registered" "gke-register" {
  name         = "gke-cluster"
  cluster_type = "default-addons-type"
}

data "google_client_config" "provider" {
}

variable "name" {
  default = "cluster-2"
}

variable "project" {
  default = "xxxx"
}

variable "location" {
  default = "us-central1-c"
}

data "google_container_cluster" "my_cluster" {
  name     = var.name
  project  = var.project
  location = var.location
}

provider "kubectl" {
  load_config_file       = false
  host                   = "https://${data.google_container_cluster.my_cluster.endpoint}"
  token                  = data.google_client_config.provider.access_token
  cluster_ca_certificate = base64decode(data.google_container_cluster.my_cluster.master_auth.0.cluster_ca_certificate)
}

// for split file and pass
data "kubectl_filename_list" "manifests" {
  pattern = "${nirmata_cluster_registered.gke-register.controller_yamls_folder}/*"
}

resource "kubectl_manifest" "test" {
  count     = nirmata_cluster_registered.gke-register.controller_yamls_count
  yaml_body = file(element(data.kubectl_filename_list.manifests.matches, count.index))
}
