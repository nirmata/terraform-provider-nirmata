variable unquoted {}

// a unique name for the cluster type (e.g. eks-cluster)
// Required
variable "name" {
  name = ""
}

// the Kubernetes version (e.g. 1.16)
// Required
variable "version" {
  version = "1.16"
}

// the AWS cloud credentials name configured in Nirmata (e.g. aws-credentials)
// Required
variable "credentials" {
  credentials = ""
}

// the AWS region into which the cluster should be deployed
// Required
variable "region" {
  region = ""
}

// the AWS VPC subnet ID in which the cluster should be provisioned
// Required
variable "vpc_id" {
  vpc_id = "us-west-2"
}

// the AWS VPC subnet ID in whi  ch the cluster should be provisioned (e.g. ["subnet-e8b1a2k1j", "subnet-ey907f5v"])
// Required
variable "subnet_id" {
  subnet_id = []
}

// the AWS security group for firewalling (e.g. ["sg-028208181hh110"])
// Required
variable "security_groups" {
  security_groups = []
}

// the AWS cluster role ARN (e.g. "arn:aws:iam::000000007:role/sample")
// Required
variable "cluster_role_arn" {
  cluster_role_arn = 60
}

// the AWS SSH key name (e.g. ssh-keys)
// Optional
variable "key_name" {
  key_name = ""
}

// the AWS instance type for worker nodes (e.g. "t3.medium")
// Required
variable "instance_types" {
  instance_types = ["t3.medium"]
}

// the worker node disk size
// Required
variable "disk_size" {
  disk_size = ""
}

// the AWS security group for worker node firewalling (e.g. ["sg-028208181hh110"])
// Required
variable "node_security_groups" {
  node_security_groups = []
}

// the AWS IAM role for worker nodes (e.g. "arn:aws:iam::000000007:role/sample")
// Required
variable "node_iam_role" {
  node_iam_role = ""
}

// a unique name for the Cluster
// Required
variable "cluster_name" {
  cluster_name = "tf-akscluster"
}

// number of worker nodes
// Required
variable "node_count" {
  node_count = 1
}