provider "nirmata" {
  // Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  // token = ""

  // Nirmata address. Defaults to https://nirmata.io and can be configured as the environment variable NIRMATA_URL.
  // url = ""
}

// A ClusterType contains reusable configuration to create clusters.
resource "nirmata_gke_clusterType" "gke-cluster-type" {

  // a unique name for the cluster type (e.g. eks-cluster)
  // Required
  name = "gke-cluster-type"

  // the GKE version (e.g. 1.16.12-gke.3)
  // Required
  version = "1.16.12-gke.3"

  // the GCP cloud credentials name configured in Nirmata (e.g. gcp-credentials)
  // Required
  // credentials = ""

  // the GCP region into which the cluster should be deployed (e.g. "us-central1-b")
  // Required
  region = "us-central1-b"

  // the GCP machine type (e.g. "e2-standard-2")
  // Required
  machine_type = "e2-standard-2"

  // the worker node disk size in GB
  // Required
  disk_size = 60
}

// A Cluster is created using a ClusterType
resource "nirmata_ProviderManaged_cluster" "gke-cluster" {

  // a unique name for the Cluster
  // Required
  name = "gke-cluster-1"

  // the cluster type
  // Required
  type_selector = nirmata_gke_clusterType.gke-cluster-type.name

  // number of worker nodes
  // Required
  node_count = 1
}

output "cluster_type_name" {
  description = "ClusterType name"
  value       = nirmata_gke_clusterType.gke-cluster-type.name
}

output "cluster_name" {
  description = "Cluster name"
  value       = nirmata_ProviderManaged_cluster.gke-cluster.name
}
