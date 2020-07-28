provider "nirmata" {
  // Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  // token = ""

  // Nirmata address. Defaults to https://nirmata.io and can be configured as the environment variable NIRMATA_URL.
  // url = ""
}

// A ClusterType contains reusable configuration to create clusters.
resource "nirmata_eks_clusterType" "eks-cluster-type" {

  // a unique name for the cluster type (e.g. eks-cluster)
  // Required
  name = "eks_us-west-2_1.16"

  // the Kubernetes version (e.g. 1.16)
  // Required
  version = "1.16"
  
  // the AWS cloud credentials name configured in Nirmata (e.g. aws-credentials)
  // Required
  // credentials = ""
  
  // the AWS region into which the cluster should be deployed
  // Required
  region= "us-west-2"
  
  // the AWS VPC subnet ID in which the cluster should be provisioned
  // Required
  // vpcid= ""
  
  // the AWS VPC subnet ID in which the cluster should be provisioned (e.g. ["subnet-e8b1a2k1j", "subnet-ey907f5v"])
  // Required
  // subnetid= []
    
  // the AWS security group to manage communications (e.g. ["sg-028208181hh110"])
  // Required
  // securitygroups= [] 
  
  // the AWS cluster role ARN (e.g. "arn:aws:iam::000000007:role/sample")
  // Required
  // clusterrolearn= ""

  // the AWS SSH key name (e.g. ssh-keys)
  // Required
  // keyname= ""
  
  // the AWS instance type for worker nodes (e.g. "t3.medium")
  instancetypes= ["t3.medium"]
  
  // the worker node disk size
  // Required
  disksize = 60
}

// A Cluster is created using a ClusterType
resource "nirmata_ProviderManaged_cluster" "eks-cluster" {

  // a unique name for the Cluster
  // Required
  name       = "eks-cluster-1"

  // the cluster type
  // Required
  type_selector  =  nirmata_eks_clusterType.eks-cluster-type.name
  
  // number of worker nodes
  // Required
  node_count = 1
}

output "cluster_type_name" {
  description = "ClusterType name"
  value       = nirmata_eks_clusterType.eks-cluster-type.name
}

output "cluster_name" {
  description = "Cluster name"
  value       = nirmata_ProviderManaged_cluster.eks-cluster.name
}
