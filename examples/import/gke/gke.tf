provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  #  token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as the environment variable NIRMATA_URL.
  #  url = ""
}

#  A nirmata_cluster created by importing an existing cluster
resource "nirmata_cluster_imported" "gke-import-1" {
  name = "my-cluster-1"
  credentials = "gke-test"
  cluster_type  =  "gke-test"
  region = "us-central1-c"
  cluster = "gke"
  project = "my-project"
  delete_action = "remove"
  vault_auth {
    name             = ""
    path             = ""
    addon_name       = ""
    credentials_id   = ""
    delete_auth_path = false

    roles {
      name                 = "nginx"
      service_account_name = "default"
      namespace            = "nginx"
      policies             = ""
    }
  }
}
