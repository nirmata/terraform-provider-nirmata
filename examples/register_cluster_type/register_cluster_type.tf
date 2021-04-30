provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  
  #  token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as
  #  the environment variable NIRMATA_URL.
  
  #  url = ""
}

#  An registered Cluster type is used while creating registered clusters
resource "nirmata_cluster_type_registered" "tf-registered-type-1" {
  name  = "tf-registered-cluster"
  cloud = "Other"
}
