variable unquoted {}

// The name of this cluster
// Required
variable "name" {
  name = "dc-hg-1"
}

// The name of this cluster
// Required
variable "direct_connect_cluster_name" {
  direct_connect_cluster_name = "dc-cluster-1"
}


// policy name
// Required
variable "policy" {
  policy = "default-v1.16.0"
}


// The region into which the cluster should be deployed
// Required
variable "region" {
  region = ""
}

// The type of VM for worker nodes
// Required
variable "machine_type" {
  vm_shape = "us-west-2"
}
