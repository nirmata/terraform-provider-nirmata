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
  name = "eks_us-west2-jim-test"

  // the Kubernetes version (e.g. 1.16)
  // Required
  version = "1.16"
  
  // the AWS cloud credentials name configured in Nirmata (e.g. aws-credentials)
  // Required
  credentials = "aws-cloud"
  
  // the AWS region into which the cluster should be deployed
  // Required
  region= "us-west-2"
  
  // the AWS VPC subnet ID in which the cluster should be provisioned
  // Required
  vpc_id= "vpc-05f45863"
  
  // the AWS VPC subnet ID in which the cluster should be provisioned (e.g. ["subnet-e8b1a2k1j", "subnet-ey907f5v"])
  // Required
  subnet_id= ["subnet-e7b3a2ae", "subnet-3e33b965"]
    
  // the AWS security group to manage communications (e.g. ["sg-028208181hh110"])
  // Required
  security_groups= ["sg-00e7a90524265bfb0"]
  
  // the AWS cluster role ARN (e.g. "arn:aws:iam::000000007:role/sample")
  // Required
  cluster_role_arn= "arn:aws:iam::094919933512:role/eks-role"

  // the AWS SSH key name (e.g. ssh-keys)
  // Required
  key_name= "nirmata-west-1-062014"
  
  // the AWS instance type for worker nodes (e.g. "t3.medium")
  instance_types= ["t3.medium"]
  
  // the worker node disk size
  // Required
  disk_size = 60

  node_security_groups = ["sg-00e7a90524265bfb0"]

  node_iam_role = "arn:aws:iam::094919933512:role/eks-arundathi-worker-node"
}

// A Cluster is created using a ClusterType
resource "nirmata_ProviderManaged_cluster" "eks-cluster" {

  // a unique name for the Cluster
  // Required
  name       = "eks-cluster-jim"

  // the cluster type
  // Required
  cluster_type  =  nirmata_eks_clusterType.eks-cluster-type.name
  
  // number of worker nodes
  // Required
  node_count = 2
}

output "cluster_type_name" {
  description = "ClusterType name"
  value       = nirmata_eks_clusterType.eks-cluster-type.name
}

output "cluster_name" {
  description = "Cluster name"
  value       = nirmata_ProviderManaged_cluster.eks-cluster.name
}
