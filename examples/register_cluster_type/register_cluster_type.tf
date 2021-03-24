provider "nirmata" {
  // Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  //
  // token = ""

  // Nirmata address. Defaults to https://nirmata.io and can be configured as
  // the environment variable NIRMATA_URL.
  //
  // url = ""
}

// An register Cluster type is used while  creating register clusters
// cloud : Other, GoogleCloudPlatform, AWS .....
resource "nirmata_cluster_type_register" "tf-register-type-1" {
  name             = "tf-register-cluster"
  cloud      = "" 
}
