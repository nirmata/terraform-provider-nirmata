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
  credentials = "gke-nirmata"

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

  // Optional
  cloud_run = true

  // Optional
  allow_override_credentials = true

  
  // One Nodepooltype Required
  // Machine Type (Required): the GCP machine type (e.g. "e2-standard-2")
  // Disk Size (Required):  the worker node disk size in GB
  // Service Account (Optional): The service account to be used to call Google Cloud APIs.
  // Node_labels (Optional)
  // Node Annotations (Optional)
  nodepooltype  {    
    machinetype = "c2-standard-16"
    disksize= 110
    enable_preemptible_nodes  =  false
    service_account = ""
    node_annotations = {
       node = "annotate"
    }
  }  

  nodepooltype  {       
    machinetype = "c2-standard-14"
    disksize = 120
    enable_preemptible_nodes  =  false
    service_account = ""
    node_labels= {
        node= "label"
    }
  } 

  // Cluster IPv4 CIDR (Optional) : Pod CIDR Range
  cluster_ipv4_cidr = ""

  // Services IPv4 CIDR (Optional) : Kubernetes Service Address Range
  services_ipv4_cidr = ""

  // Enable Network Policy(Optional)
  enable_network_policy = false

  // Enable HTTP Load Balancing (Optional)
  http_load_balancing = false

  //Enable Horizontal Pod Autoscaling (Optional)
  enable_vertical_pod_autoscaling =  false

  //Enable Vertical Pod Autoscaling (Optional)
  horizontal_pod_autoscaling = true

  // Enable Maintenance Policy (Optional)
  enable_maintenance_policy = false

  // Maintenance Duration (Daily) Start Time
  start_time = ""

  // Maintenance sDuration in Hours 
  duration = "10"

  // Maintenance Exclusions (Optional)
  exclusion_timewindow = {}

  //System metadata will be available on the cluster in a config map.
  // Optional
  system_metadata = {
    cluster = "gke"
  }

 cluster_field_override = {
    network = "String"
    subnetwork = "String"
  }

  nodepool_field_override = {
    disksize = "Integer"
    machinetype = "String"
  }
  
}

// A nirmata_cluster is created using a cluster_type
resource "nirmata_cluster" "gke-cluster-1" {

  // a unique name for the Cluster
  // Required
  name = "gke-cluster-1"

  // the cluster type
  // Required
  cluster_type  =  nirmata_cluster_type_gke.gke-cluster-type-1.name

  // Set the desired number of nodes that the group should launch with initially.
  // Required
  node_count = 1

  // Optional
  // update the credential.
  override_credentials = ""

  //System metadata will be available on the cluster in a config map.
  // Optional
  system_metadata = {
    cluster = "gke"
  }
}

output "cluster_type_name" {
  description = "ClusterType name"
  value       = nirmata_cluster_type_gke.gke-cluster-type-1.name
}

output "cluster_name" {
  description = "Cluster name"
  value       = nirmata_cluster.gke-cluster-1.name
}
