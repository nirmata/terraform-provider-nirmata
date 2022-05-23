provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  #  token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as the environment variable NIRMATA_URL.
  #  url = ""
}

#  A nirmata_cluster created by importing an existing cluster
resource "nirmata_cluster_imported" "eks-import" {
  name = "my-cluster-1"
  credentials = "eks-test"
  cluster_type  =  "eks-test"
  region = "us-west-1"
  delete_action = "remove"
}
