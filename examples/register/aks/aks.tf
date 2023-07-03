// NOTE: this example needs to be applied in two phases, as the YAML file count
// is computed during the apply phase of the nirmata_cluster_registered resource.
//
// Steps:
//   terraform init
//   terraform plan 
//   terraform apply -target nirmata_cluster_registered.aks-registered
//   terraform plan
//   terraform apply

// create a new cluster and download the controller YAMLs
resource "nirmata_cluster_registered" "aks-registered" {
  name         = "aks-cluster-test-03"
  cluster_type = "default-add-ons"
  controller_yamls_folder = "/Users/atul/Work/git/nirmata/terraform-provider-nirmata/examples/register/aks/yamls"
}

provider "nirmata" {

  // Nirmata address.
  url = "https://nirmata.io"

  // Nirmata API Key. Also configurable using the environment variable NIRMATA_TOKEN.
  token = "im+0Qe/LaV6V0CYkjc2gdTMdzKl23gIGSTZVbnnKXLeRAJld9ivWA2vlCNo6pbV6+Cm4FO5E/NU2kPCz1Bajrw=="
}

# Retrieve AKS cluster information
provider "azurerm" {
  features {}
}

data "azurerm_kubernetes_cluster" "cluster" {
  name                = "atul-test-02"
  resource_group_name = "atul-test-0_group"
}

provider "kubectl" {
  host = data.azurerm_kubernetes_cluster.cluster.kube_config.0.host

  client_certificate     = base64decode(data.azurerm_kubernetes_cluster.cluster.kube_config.0.client_certificate)
  client_key             = base64decode(data.azurerm_kubernetes_cluster.cluster.kube_config.0.client_key)
  cluster_ca_certificate = base64decode(data.azurerm_kubernetes_cluster.cluster.kube_config.0.cluster_ca_certificate)

}

data "kubectl_filename_list" "namespace" {
   pattern = "${nirmata_cluster_registered.aks-registered.controller_yamls_folder}/temp-01-*"
}

data "kubectl_filename_list" "crd" {
   pattern = "${nirmata_cluster_registered.aks-registered.controller_yamls_folder}/temp-02-*"
}

data "kubectl_filename_list" "deployment" {
   pattern = "${nirmata_cluster_registered.aks-registered.controller_yamls_folder}/temp-03-*"
}


// Register Nirmata Cluster
resource "nirmata_cluster_registered" "aks-registered" {
  name         = var.nirmata_cluster_name
  cluster_type = var.nirmata_cluster_type
}

// apply the controller YAMLs
resource "kubectl_manifest" "namespace" {
  wait        = true
  count       = nirmata_cluster_registered.aks-registered.controller_ns_yamls_count
  yaml_body   = file(element(data.kubectl_filename_list.namespace.matches, count.index))
  apply_only  = true
}

resource "kubectl_manifest" "crd" {
  wait        = true
  count       = nirmata_cluster_registered.aks-registered.controller_crd_yamls_count
  yaml_body   = file(element(data.kubectl_filename_list.crd.matches, count.index))
  apply_only  = true
  depends_on  = [kubectl_manifest.namespace]
}

resource "kubectl_manifest" "deployment" {
  wait        = true
  count       = nirmata_cluster_registered.aks-registered.controller_deploy_yamls_count
  yaml_body   = file(element(data.kubectl_filename_list.deployment.matches, count.index))
  apply_only  = true
  depends_on  = [kubectl_manifest.crd]
}

