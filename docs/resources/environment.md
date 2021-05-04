---
page_title: "nirmata_environment Resource"
---

# nirmata_environment Resource

A virtual cluster backed by a namespace and security policies to allow sharing of cluster resources.

## Example Usage

```hcl

resource "nirmata_environment" "tf-env-1" {
  name        = "tf-env-1"
  type        = "medium"
  cluster     = "prod-demo"
  namespace   = "tf-ns-1"
}

```

## Argument Reference

* `name` - (Required) a unique name for the environment.
* `type` - (Required) the environnment type.
* `cluster` - (Required)  the kubernetes cluster.
* `namespace` - (Optional) the cluster namespace bound to this environment. Defaults to the environment name.