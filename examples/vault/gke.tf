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
    node_annotations = {
      node = "annotate"
    }
  }

  addons {
    name            = "vault-agent-injector"
    addon_selector  = "vault-agent-injector"
    catalog         = "default-catalog"
    channel         = "Stable"
    sequence_number = 15
  }

  vault_auth {
    name             = "gke-vault-auth"
    path             = "nirmata/$(cluster.name)"
    addon_name       = "vault-agent-injector"
    credentials_name = "vault_access"

    roles {
      name                 = "sample-role"
      service_account_name = "application-sample-sa"
      namespace            = "application-sample-ns"
      policies             = "application-sample-policy"
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
