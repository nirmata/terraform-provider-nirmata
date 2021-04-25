// This example first creates an EKS cluster type and then creates
// a EKS cluster using that type.


// 1. Create a EKS cluster type
resource "nirmata_cluster_type_eks" "eks-cluster-type-1" {
  name                      = "tf-eks-cluster-type-1"
  version                   = "1.18"
  credentials               = "aws-xxxxx"
  region                    = "us-west-2"
  vpc_id                    = "vpc-xxxxxxxx"
  subnet_id                 = ["subnet-xxxxxxxx", "subnet-xxxxxxxx"]
  security_groups           = ["sg-xxxxxxxxxxxxxxxx"]
  cluster_role_arn          = "arn:aws:iam::xxxxxxxx:role/xxxxxxxx"
  enable_private_endpoint   = true
  enable_identity_provider  = true
  auto_sync_namespaces       = false
  // enable_secrets_encryption = true
  // kms_key_arn = ""
  // log_types = ""

  //enable_fargate = true
  //pod_execution_role_arn = ""
  //subnets= ["sg-xxxxxxxxxxxxxxxx","sg-xxxxxxxxxxxxxxxx"]
  //namespace_label_selectors = {}
  //pod_label_selectors = {}

  nodepools {
    name                = "default"
    instance_type       = "t3.medium"
    disk_size           = 60
    ssh_key_name        = "xxxxxxxx"
    security_groups     = ["sg-xxxxxxxxxxxxxxxx"]
    iam_role            = "arn:aws:iam::xxxxxxxx:role/eks-xxxxxxxx"
    // ami_type = ""
    // image_id = ""
  }
}

// 2. Create a cluster
resource "nirmata_cluster" "eks-cluster-1" {
  name                 = "eks-cluster-1"
  cluster_type         = nirmata_cluster_type_eks.eks-cluster-type-1.name
  node_count           = 1
}

output "cluster_type_name" {
  description = "ClusterType name"
  value       = nirmata_cluster_type_eks.eks-cluster-type-1.name
}

output "cluster_name" {
  description = "Cluster name"
  value       = nirmata_cluster.eks-cluster-1.name
}
