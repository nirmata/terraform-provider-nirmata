---
page_title: "nirmata_kind_cluster_registered Resource"
---

# nirmata_cluster_registered Resource

An existing kind cluster that is registered using local Kubernetes credentials.

## Example Usage

Register an existing KIND cluster using Kubernetes credentials. The new cluster is created and the controller YAMLs are downloaded to a temporary folder. The YAMLs are then applied to the existing kind cluster using the `kubectl` providers.

**NOTE:** this example needs to be applied in two phases, as the YAML file count is computed during the apply phase of the nirmata_cluster_registered resource. Steps:
1. terraform init
2. terraform plan 
3. terraform apply -target nirmata_cluster_registered.kind-registered
4. terraform plan
5. terraform apply

```hcl

resource "nirmata_cluster_registered" "kind-registered" {
  name         = "kind-cluster"
  cluster_type = "default-add-ons"
}

variable "host" {
  type = string
}

variable "client_certificate" {
  type = string
}

variable "client_key" {
  type = string
}

variable "cluster_ca_certificate" {
  type = string
}

provider "kubectl" {
  host = var.host
  client_certificate     = base64decode(var.client_certificate)
  client_key             = base64decode(var.client_key)
  cluster_ca_certificate = base64decode(var.cluster_ca_certificate)
}

data "kubectl_filename_list" "manifests" {
  pattern = "${nirmata_cluster_registered.kind-registered.controller_yamls_folder}/*"
}

// apply the controller YAMLs
resource "kubectl_manifest" "test" {
  count     = nirmata_cluster_registered.kind-registered.controller_yamls_count
  yaml_body = file(element(data.kubectl_filename_list.manifests.matches, count.index))
}


```

## terraform.tfvars
```
host                   = "https://127.0.0.1:32768"
client_certificate     = "LS0tLS1CRUdJTiB..."
client_key             = "LS0tLS1CRUdJTiB..."
cluster_ca_certificate = "LS0tLS1CRUdJTiB..."
```


## Argument Reference

* `name` - (Required) Enter a unique name for the cluster.
* `cluster_type` - (Required) Enter the cluster type to be applied to the cluster.
* `labels` - (Optional) This field indicates the labels to be set on the cluster.
* `delete_action` - (Optional) This field indicates whether to delete or remove the cluster on destroy. The default value is `remove`.

* `host` -  clusters.cluster.server.
* `client_certificate` - users.user.client-certificate.
* `client_key` - users.user.client-key.
* `cluster_ca_certificate` - clusters.cluster.certificate-authority-data.



