
provider "nirmata" {
  // Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  // token = ""

  // Nirmata address. Defaults to https://nirmata.io and can be configured as the environment variable NIRMATA_URL.
  // url = ""
}

data "google_client_config" "provider" {
}

data "google_container_cluster" "my_cluster" {
  name        =  var.name
  project     = var.project
  location    = var.location
}

// A nirmata_cluster created by registered an existing GKE cluster
resource "nirmata_cluster_registered" "gke-register" {
  name = "gke-cluster-tf-n"
  cluster_type  =  "default-addons-type"
}
