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
  name         = "aks-cluster"
  cluster_type = "default-add-ons"
}

# Retrieve AKS cluster information
provider "azurerm" {
  features {}
}

data "azurerm_kubernetes_cluster" "cluster" {
  name                = "tf-test"
  resource_group_name = "nirmata-mustufa-poc"
}

provider "kubectl" {
  host = data.azurerm_kubernetes_cluster.cluster.kube_config.0.host

  client_certificate     = base64decode(data.azurerm_kubernetes_cluster.cluster.kube_config.0.client_certificate)
  client_key             = base64decode(data.azurerm_kubernetes_cluster.cluster.kube_config.0.client_key)
  cluster_ca_certificate = base64decode(data.azurerm_kubernetes_cluster.cluster.kube_config.0.cluster_ca_certificate)

}
data "kubectl_filename_list" "manifests" {
  pattern = "${nirmata_cluster_registered.aks-registered.controller_yamls_folder}/*"
}

// apply the controller YAMLs
resource "kubectl_manifest" "test" {
  count     = nirmata_cluster_registered.aks-registered.controller_yamls_count
  yaml_body = file(element(data.kubectl_filename_list.manifests.matches, count.index))
}