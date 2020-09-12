variable unquoted {}

// a unique name for the cluster type (e.g. az-cluster)
// Required
variable "name" {
  name = ""
}

// a unique name for the cluster type (e.g. eks-cluster)
// Required
variable "version" {
  version = "1.16.12-gke.3"
}

// the GCP cloud credentials name configured in Nirmata (e.g. gcp-credentials)
// Required
variable "credentials" {
  credentials = ""
}

// the GCP region into which the cluster should be deployed (e.g. "us-central1-b")
// Required
variable "region" {
  region = ""
}

// the GCP machine type (e.g. "e2-standard-2")
// Required
variable "machine_type" {
  machine_type = "us-west-2"
}

// the worker node disk size in GB
// Required
variable "disk_size" {
  disk_size = 60
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