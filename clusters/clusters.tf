provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
}

resource "nirmata_ProviderManaged_cluster" "cluster-1" {
  name       = "pvm-cluster-1"
  type_selector  = "autoscaling"
  node_count = 1
}