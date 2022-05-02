---
page_title: "nirmata_aks_cluster_registered Resource"
---

# nirmata_cluster_registered Resource

An existing aks cluster that is registered using local Kubernetes credentials.

## Example Usage

Register an existing AKS cluster using Kubernetes credentials. The new cluster is created and the controller YAMLs are downloaded to a temporary folder. The YAMLs are then applied to the existing aks cluster using the `kubectl` providers.

**NOTE:** this example needs to be applied in two phases, as the YAML file count is computed during the apply phase of the nirmata_cluster_registered resource. Steps:
1. terraform init
2. terraform plan 
3. terraform apply -target nirmata_cluster_registered.aks-registered
4. terraform plan
5. terraform apply

```hcl

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


```


## Argument Reference

* `name` - (Required) a unique name for the cluster.
* `cluster_type` - (Required) the cluster type to apply.
* `controller_yamls` - (Computed) the controller YAML
* `controller_yamls_folder` - (Computed) a local temporary folder with the controller YAML files
* `controller_yamls_count` - (Computed) the controller YAML file count
* `labels` - (Optional) labels to set on cluster.
* `delete_action` - (Optional) whether to delete or remove the cluster on destroy. Defaults to `remove`.




