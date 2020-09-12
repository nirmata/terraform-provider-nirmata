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
  name = var.name

  // the Kubernetes version (e.g. 1.16)
  // Required
  version = var.version

  // the AWS cloud credentials name configured in Nirmata (e.g. aws-credentials)
  // Required
  // credentials = var.credentials

  // the AWS region into which the cluster should be deployed
  // Required
  region = var.region

  // the AWS VPC subnet ID in which the cluster should be provisioned
  // Required
  // vpc_id = var.vpc_id

  // the AWS VPC subnet ID in which the cluster should be provisioned (e.g. ["subnet-e8b1a2k1j", "subnet-ey907f5v"])
  // Required
  // subnet_id = var.subnet_id

  // the AWS security group for firewalling (e.g. ["sg-028208181hh110"])
  // Required
  // security_groups = var.security_groups

  // the AWS cluster role ARN (e.g. "arn:aws:iam::000000007:role/sample")
  // Required
  // cluster_role_arn = var.cluster_role_arn

  // the AWS SSH key name (e.g. ssh-keys)
  // Required
  // key_name= var.key_name

  // the AWS instance type for worker nodes (e.g. "t3.medium")
  instance_types = var.instance_types

  // the worker node disk size
  // Required
  disk_size = var.disk_size

  // the AWS security group for worker node firewalling (e.g. ["sg-028208181hh110"])
  // Required
  // node_security_groups = var.node_security_groups

  // the AWS IAM role for worker nodes (e.g. "arn:aws:iam::000000007:role/sample")
  // Required
  // node_iam_role = var.node_iam_role
}

// A Cluster is created using a ClusterType
resource "nirmata_ProviderManaged_cluster" "eks-cluster" {

  // a unique name for the Cluster
  // Required
  name = var.cluster_name

  // the cluster type
  // Required
  type_selector = nirmata_eks_clusterType.eks-cluster-type.name

  // number of worker nodes
  // Required
  node_count = var.node_count
}


