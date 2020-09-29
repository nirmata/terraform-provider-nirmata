---
page_title: "Nirmata Provider"
---

# Nirmata Provider

The Nirmata Provider automates Kubernetes cluster and workload management using [Nirmata](https://nirmata.com).

## Example Usage

```hcl
# configure the Nirmata provider
provider "nirmata" {
 
  # Nirmata address.
  url = "https://nirmata.io"

  // Nirmata API Key. Also configurable using the environment variable NIRMATA_TOKEN.
  token = var.nirmata.token

}
```

```hcl
# create a cluster using the Nirmata provider and an existing cluster type
resource "nirmata_cluster" "gke-1" {
  name = "gke-1"
  cluster_type = "gke-us-west"
  node_count = 1
}
```

## Argument Reference

* `url` - (Required) Nirmata API url. Also configurable using the `NIRMATA_URL` environment variable.
* `token` - (Optional/Sensitive) Nirmata API access token. Also configurable using the `NIRMATA_TOKEN` environment variable.
