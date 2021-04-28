---
page_title: "nirmata_cluster_registered Resource"
---

# nirmata_cluster_registered Resource

An existing cluster that is registered using local Kubernetes credentials.

## Example Usage

Register an existing GKE cluster using Kubernetes credentials. The new cluster is created and the controller YAMLs are downloaded to a temporary folder. The YAMLs are then applied to the existing cluster using the `kubectl` and `google_client_config` providers.


```hcl

resource "nirmata_cluster_registered" "gke-register" {
  name         = "gke-cluster"
  cluster_type = "default-addons-type"
}

data "google_client_config" "provider" {
}

variable "name" {
  default = "cluster-2"
}

variable "project" {
  default = "xxxx"
}

variable "location" {
  default = "us-central1-c"
}

data "google_container_cluster" "my_cluster" {
  name     = var.name
  project  = var.project
  location = var.location
}

provider "kubectl" {
  load_config_file       = false
  host                   = "https://${data.google_container_cluster.my_cluster.endpoint}"
  token                  = data.google_client_config.provider.access_token
  cluster_ca_certificate = base64decode(data.google_container_cluster.my_cluster.master_auth.0.cluster_ca_certificate)
}

// for split file and pass
data "kubectl_filename_list" "manifests" {
  pattern = "${nirmata_cluster_registered.gke-register.controller_yamls_folder}/*"
}

resource "kubectl_manifest" "test" {
  count     = length(data.kubectl_filename_list.manifests.matches)
  yaml_body = file(element(data.kubectl_filename_list.manifests.matches, count.index))
}

```

## Argument Reference

* `name` - (Required) a unique name for the cluster.
* `cluster_type` - (Required) the cluster type to apply.
* `controller_yamls` - (Computed) the controller YAML
* `controller_yamls_folder` - (Computed) a local temporary folder with the controller YAML files
* `delete_action` - (Optional) whether to delete or remove the cluster on destroy. Defaults to `remove`.



