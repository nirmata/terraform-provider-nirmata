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
  labels = {foo = "bar"}
  environment_update_action   = "notify" 
}

```

## Argument Reference

* `name` - (Required) Enter a unique name for the environment.
* `type` - (Required) Enter the environnment type.
* `cluster` - (Required)  Enter the Kubernetes cluster.
* `labels` - (Optional) This field indicates the labels to be set on the add-on application's environment.
* `namespace` - (Optional) This field indicates the cluster namespace bound to this environment. It defaults to the environment name.
* `environment_update_action` - (Optional) By default, this value is set to notify and to update if changes need to apply automatically.