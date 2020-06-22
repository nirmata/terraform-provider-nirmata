provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
}
resource "nirmata_eks_clusterType" "cluster-type-eks" {
  name       = "eks-cluster-type-tf"
  version  = "1.14" //The version of Kubernetes that should be used for this cluster.
  credentials = "" //cloud credentials that hosts this cluster
  region= "us-east-2" //The  region into which the cluster should be deployed
  vpcid= ""  // VPC enables you to launch AWS resources into a virtual network that you've defined
  subneid= []
  clusterrolearn= ""  //An role is an identity within your AWS account that has specific permissions
  securitygroups= []  //Security groups control communications within the Amazon EKS cluster including between the managed Kubernetes control plane and compute resources in your AWS account such as worker nodes and Fargate pods.
  keyname= ""  
  instancetypes= ["t2.small"]
  disksize= 120
}

