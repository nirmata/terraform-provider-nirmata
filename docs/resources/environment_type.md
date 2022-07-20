---
page_title: "nirmata_environment_type Resource"
---

# nirmata_environment_type Resource

A reusable configuration set for ennvironments.

## Example Usage

```hcl

resource "nirmata_environment_type" "tf-env-type-1" {
  name             = "tf-env"
  is_default      = false
  resource_limits  = {
      cpu = "500m",
      memory = "4Gi"
      pod = "50"
      storage = "1Gi"
  }
  labels = {foo = "bar"}
}

```

## Argument Reference

* `name` - (Required) Enter a unique name for the environment type.
* `is_default` - (Optional) This field uses the default environment type.
* `labels` - (Optional) This field indicates that the labels to set on the add-on application's environment type.
* `resource_limits` - (Required) Enter a map of resource limits for the environment. The commonly used resources include `cpu`, `memory`, and `storage`. Refer the [Kubernetes docs](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/) for a complete reference.
