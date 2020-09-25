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
  version = "1.16.13-gke.1" 

  // the GCP cloud credentials name configured in Nirmata (e.g. gcp-credentials)
  // Required
  //credentials = ""
  
  // the GCP region into which the cluster should be deployed (e.g. "us-central1-b")
  // Required
  region= "us-central1-b"

  // the GCP machine type (e.g. "e2-standard-2")
  // Required
  //machine_type = ""
  
  // the worker node disk size in GB
  // Required
  disk_size = 60

  //A regional cluster has multiple replicas of the control plane, running in multiple zones within a given region. A zonal cluster has a single replica of the control plane running in a single zone.
  // Optional
  // default set to Regional
  //e.g. (Regional,Zonal)
  //location_type=  "Regional" 
  
  // nodes should be deployed. Selecting more than one zone increases availability.  (e.g. ["asia-east1-a"])
  // Required 
  //node_locations = []

  // Optional
  // Protect your Kubernetes Secrets with envelope encryption.
  // default set to false
  //enable_secrets_encryption = false

  // Workload Identity is the recommended way to access Google Cloud services from applications running within GKE due to its improved security properties and manageability.
  // Optional 
  // default set to false
  //enable_workload_identity = false

  //Enter the Workload Pool for your project. Workload Identity relies on a Workload Pool to aggregate identity across multiple clusters.
  // Required if enable_secrets_encryption is true
  //workload_pool = ""

  //Enter the Resource ID of the key you want to use (e.g. projects/project-name/locations/global/keyRings/my-keyring/cryptoKeys/my-key)
  // Required if enable_workload_identity is true
  //secrets_encryption_key = ""

  // Required (e.g. "default")
  //network = ""

 // Required (e.g. "default")
  //subnetwork = ""

  
 
}

// A Cluster is created using a ClusterType
resource "nirmata_ProviderManaged_cluster" "gke-cluster" {

  // a unique name for the Cluster
  // Required
  name = "gke-cluster"
  
  // the cluster type
  // Required
  type_selector  =  nirmata_gke_clusterType.gke-cluster-type.name
  
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
