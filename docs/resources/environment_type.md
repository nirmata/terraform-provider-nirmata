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

* `name` - (Required) a unique name for the environment type.
* `is_default` - (Optional) use as the default environment type.
* `labels` - (Optional) labels to set on  environment type add-on.
* `resource_limits` - (Required) a map of resource limits for the environment. Commonly used resources include `cpu`, `memory`, and `storage`. Check the [Kubernetes docs](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/) for a complete reference.
