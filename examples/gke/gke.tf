provider "nirmata" {
  // Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  // token = ""

  // Nirmata address. Defaults to https://nirmata.io and can be configured as the environment variable NIRMATA_URL.
  // url = ""
}

// A ClusterType contains reusable configuration to create clusters.
resource "nirmata_cluster_type_gke" "gke-cluster-type-1" {

  name                       = "gke-cluster-type-test"
  version                    = "1.17.13-gke.2001"
  credentials                = "gke-test"
  location_type              = "Zonal"
  zone                       = "us-central1-a"
  network                    = "default"
  subnetwork                 = "default"
  enable_cloud_run           = false
  enable_http_load_balancing = false
  allow_override_credentials = true
  channel                    = "REGULAR"

  system_metadata = {
    cluster = "gke"
  }

  cluster_field_override = {
    network    = "String"
    subnetwork = "String"
  }

  nodepool_field_override = {
    disk_size    = "Integer"
    machine_type = "String"
  }

  nodepools {
    machine_type             = "c2-standard-16"
    disk_size                = 110
    enable_preemptible_nodes = true
    service_account          = ""
    auto_upgrade             = true
    auto_repair              = true
    node_annotations = {
      node = "annotate"
    }
  }
}

// A nirmata_cluster is created using a cluster_type
resource "nirmata_cluster" "gke-cluster-1" {

  name                 = "gke-cluster-1"
  cluster_type         = nirmata_cluster_type_gke.gke-cluster-type-1.name
  node_count           = 1
  override_credentials = ""

  system_metadata = {
    cluster = "gke"
  }

  cluster_field_override = {
    network    = ""
    subnetwork = ""
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
