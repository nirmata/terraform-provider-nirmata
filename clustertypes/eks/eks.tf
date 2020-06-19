provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
}
resource "nirmata_eks_clusterType" "cluster-type-eks" {
  name       = "eks-cluster-type-tf"
  version  = "1.14"
  credentials = "eks"
  region= "us-east-2"
  vpcid= ""
  subneid= []
  clusterrolearn= ""
  securitygroups= []
  keyname= ""
  instancetypes= ["t2.small"]
  disksize= 120
}

