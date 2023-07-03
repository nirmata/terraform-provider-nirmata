// NOTE: this example needs to be applied in two phases, as the YAML file count
// is computed during the apply phase of the nirmata_cluster_registered resource.
//
// Steps:
//   terraform init
//   terraform plan 
//   terraform apply -target nirmata_cluster_registered.gke-registered
//   terraform plan
//   terraform apply

// create a new cluster and download the controller YAMLs
resource "nirmata_cluster_registered" "gke-registered" {
  name         = "gke-cluster"
  cluster_type = "default-add-ons"
}

// fetch the GKE cluster details (requires external configuration)
data "google_client_config" "provider" {
}

data "google_container_cluster" "my_cluster" {
  name     = "cluster-1"
  project  = "nirmata-demo"
  location = "us-central1-c"
}

// configure kubectl with GKE access details
provider "kubectl" {
  load_config_file       = false
  host                   = "https://${data.google_container_cluster.my_cluster.endpoint}"
  token                  = data.google_client_config.provider.access_token
  cluster_ca_certificate = base64decode(data.google_container_cluster.my_cluster.master_auth.0.cluster_ca_certificate)
}

data "kubectl_filename_list" "namespace" {
   pattern = "${nirmata_cluster_registered.gke-registered.controller_yamls_folder}/temp-01-*"
}

data "kubectl_filename_list" "crd" {
   pattern = "${nirmata_cluster_registered.gke-registered.controller_yamls_folder}/temp-02-*"
}

data "kubectl_filename_list" "deployment" {
   pattern = "${nirmata_cluster_registered.gke-registered.controller_yamls_folder}/temp-03-*"
}


// Register Nirmata Cluster
resource "nirmata_cluster_registered" "gke-registered" {
  name         = var.nirmata_cluster_name
  cluster_type = var.nirmata_cluster_type
}

// apply the controller YAMLs
resource "kubectl_manifest" "namespace" {
  wait        = true
  count       = nirmata_cluster_registered.gke-registered.controller_ns_yamls_count
  yaml_body   = file(element(data.kubectl_filename_list.namespace.matches, count.index))
  apply_only  = true
}

resource "kubectl_manifest" "crd" {
  wait        = true
  count       = nirmata_cluster_registered.gke-registered.controller_crd_yamls_count
  yaml_body   = file(element(data.kubectl_filename_list.crd.matches, count.index))
  apply_only  = true
  depends_on  = [kubectl_manifest.namespace]
}

resource "kubectl_manifest" "deployment" {
  wait        = true
  count       = nirmata_cluster_registered.gke-registered.controller_deploy_yamls_count
  yaml_body   = file(element(data.kubectl_filename_list.deployment.matches, count.index))
  apply_only  = true
  depends_on  = [kubectl_manifest.crd]
}

