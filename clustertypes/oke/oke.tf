provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
}
resource "nirmata_oke_clusterType" "cluster-type-oke" {
  name       = "oke-cluster-type"
  version  = "v1.14.8"
  credentials = "oracle-nirmata"
  region= "us-phoenix-1"
  vmshape= "VM.Standard2.1"
}