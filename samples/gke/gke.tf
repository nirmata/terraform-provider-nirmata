provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
}

resource "nirmata_cluster_gke" "dc-gke-1" {
  name       = "dc-gh-1"
  disk_size  = 60
  node_type  = "n1-standard-2"
  node_count = 1
  region     = "us-central1-a"
}