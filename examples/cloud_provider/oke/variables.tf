variable unquoted {}

// The name of this cluster
// Required
variable "name" {
  name = ""
}

// The version of Kubernetes that should be used for this cluster.
// Required
variable "version" {
  version = "1.16.12-gke.3"
}

// Cloud credentials to use for this cluster
// Required
variable "credentials" {
  credentials = ""
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
