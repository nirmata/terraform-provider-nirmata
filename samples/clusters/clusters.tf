provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
}

resource "nirmata_ProviderManaged_cluster" "cluster-1" {
  name       = "pvm-cluster"
  type_selector  = "take"
  node_count = 1
}