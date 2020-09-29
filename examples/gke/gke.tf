provider "nirmata" {
  // Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  // token = ""

  // Nirmata address. Defaults to https://nirmata.io and can be configured as the environment variable NIRMATA_URL.
  // url = ""
}

// A ClusterType contains reusable configuration to create clusters.
resource "nirmata_cluster_type_gke" "gke-cluster-type-1" {

  // a unique name for the cluster type (e.g. eks-cluster)
  // Required
  name = "gke-cluster-type-1"

  // the GKE version (e.g. 1.17.9-gke.6300)
  // Required
  version = "1.17.9-gke.6300"

  // the GCP cloud credentials name configured in Nirmata (e.g. gcp-credentials)
  // Required
  credentials = "gke-test"

  //A regional cluster has multiple replicas of the control plane, running in multiple zones within a given region. A zonal cluster has a single replica of the control plane running in a single zone.
  // Required (Regional or Zonal)
  location_type =  "Zonal"

  // the GCP region into which the cluster should be deployed (e.g. "us-central1")
  // Required if location_type is Regional
  //region = "us-central1"

  // the GCP zone into which the cluster should be deployed (e.g. "us-central1-a")
  // Required if location_type is Zonal
  zone = "us-central1-a"

  // nodes should be deployed. Selecting more than one zone increases availability.  (e.g. ["asia-east1-a"])
  // Required if location_type is Regional
  //node_locations = []

  // Optional
  // Protect your Kubernetes Secrets with envelope encryption.
  // default set to false
  //enable_secrets_encryption = false

  // Workload Identity is the recommended way to access Google Cloud services from applications running within GKE due to its improved security properties and manageability.
  // Optional
  // default set to false
  //enable_workload_identity = false

  //Enter the Workload Pool for your project.
  // Workload Identity relies on a Workload Pool to aggregate identity across multiple clusters.
  // Required if enable_secrets_encryption is true
  //workload_pool = ""

  //Enter the Resource ID of the key you want to use (e.g. projects/project-name/locations/global/keyRings/my-keyring/cryptoKeys/my-key)
  // Required if enable_workload_identity is true
  //secrets_encryption_key = ""

  // Required (e.g. "default")
  network = "default"

  // Required (e.g. "default")
  subnetwork = "default"

  // the GCP machine type (e.g. "e2-standard-2")
  // Required
  machine_type = "e2-standard-2"

  // the worker node disk size in GB
  // Required
  disk_size = 60
}

// A nirmata_cluster is created using a cluster_type
resource "nirmata_cluster" "gke-cluster-1" {

  // a unique name for the Cluster
  // Required
  name = "gke-cluster-1"

  // the cluster type
  // Required
  cluster_type  =  nirmata_cluster_type_gke.gke-cluster-type-1.name

  // number of worker nodes
  // Required
  node_count = 1
}

output "cluster_type_name" {
  description = "ClusterType name"
  value       = nirmata_cluster_type_gke.gke-cluster-type-1.name
}

output "cluster_name" {
  description = "Cluster name"
  value       = nirmata_cluster.gke-cluster-1.name
}
