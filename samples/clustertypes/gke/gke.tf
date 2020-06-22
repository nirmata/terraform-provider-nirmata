provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
}
resource "nirmata_gke_clusterType" "cluster-type-gke" {
  name       = "gke-cluster-type"
  version  = "" //The version of Kubernetes that should be used for this cluster.
  credentials = "" //cloud credentials that hosts this cluster
  region= "" //The  region into which the cluster should be deployed
  machinetype = ""
  disksize = 10
}
