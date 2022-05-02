---
page_title: "nirmata_cluster_registered Resource"
---

# nirmata_cluster_registered Resource

An existing cluster that is registered using local Kubernetes credentials.

## Example Usage

Register an existing GKE cluster using Kubernetes credentials. The new cluster is created and the controller YAMLs are downloaded to a temporary folder. The YAMLs are then applied to the existing cluster using the `kubectl` and `google_client_config` providers.

**NOTE:** this example needs to be applied in two phases, as the YAML file count is computed during the apply phase of the nirmata_cluster_registered resource. Steps:
1. terraform init
2. terraform plan 
3. terraform apply -target nirmata_cluster_registered.gke-registered
4. terraform plan
5. terraform apply

```hcl

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
  project  = "xxx"
  location = "us-central1-c"
}

// configure kubectl with GKE access details
provider "kubectl" {
  load_config_file       = false
  host                   = "https://${data.google_container_cluster.my_cluster.endpoint}"
  token                  = data.google_client_config.provider.access_token
  cluster_ca_certificate = base64decode(data.google_container_cluster.my_cluster.master_auth.0.cluster_ca_certificate)
}

// read YAMLs from folder
data "kubectl_filename_list" "manifests" {
  pattern = "${nirmata_cluster_registered.gke-registered.controller_yamls_folder}/*"
}

// apply the controller YAMLs
resource "kubectl_manifest" "test" {
  count     = nirmata_cluster_registered.gke-registered.controller_yamls_count
  yaml_body = file(element(data.kubectl_filename_list.manifests.matches, count.index))
}


```

## Argument Reference

* `name` - (Required) a unique name for the cluster.
* `cluster_type` - (Required) the cluster type to apply.
* `labels` - (Optional) labels to set on cluster.
* `delete_action` - (Optional) whether to delete or remove the cluster on destroy. Defaults to `remove`.



