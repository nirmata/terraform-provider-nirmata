provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
}

resource "nirmata_oke_clusterType" "cluster-type-oke" {

  // The name of this cluster
  // Required
  name       = "oke-cluster-type"

  // The version of Kubernetes that should be used for this cluster.
  // Required
  version  = ""

  // Cloud credentials to use for this cluster
  // Required
  credentials = ""

  // The region into which the cluster should be deployed
  // Required
  region= " "

  // The type of VM for worker nodes
  // Required
  vm_shape= ""
}
