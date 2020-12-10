
provider "nirmata" {
  // Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  // token = ""

  // Nirmata address. Defaults to https://nirmata.io and can be configured as the environment variable NIRMATA_URL.
  // url = ""
}

// A nirmata_cluster created by importing an existing cluster


resource "nirmata_cluster_registered" "gke-register-1" {
  name = "my-cluster-1"
  cluster_type  =  "default-add-ons"
}
// A nirmata_cluster created by importing an existing cluster
data "google_client_config" "default" {
}

data "google_container_cluster" "my_cluster" {
 name        =  var.name
  project     = var.project
  location    = var.location
}



provider "kubernetes" {
  load_config_file = false
  host  = "https://${data.google_container_cluster.my_cluster.endpoint}"
  token = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(
    data.google_container_cluster.my_cluster.master_auth[0].cluster_ca_certificate,
  )
}
