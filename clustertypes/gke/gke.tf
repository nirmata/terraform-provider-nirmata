provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
}
resource "nirmata_gke_clusterType" "cluster-type-gke" {
  name       = "gke-cluster-type"
  version  = "asia-east1-a"
  credentials = "gke-nirmata"
  region= "1.16.9-gke.2"
  machinetype = "e2-highcpu-16"
  disksize = 10
}
