---
page_title: "nirmata_catalog_application Resource"
---

# nirmata_catalog_application Resource

 Application is a group of workloads, routing and storage configurations.

## Example Usage

```hcl

resource "nirmata_catalog_application" "tf-catalog-app" {
  name              = "tf-catalog-app"
  catalog           = ""
  yamls             = file("${path.module}/fo.yaml")
}

```

## Argument Reference

* `name` - (Required) a unique name for the application in catalog.
* `catalog` - (Required) the name of catalog.
* `yamls` - (Required)  path for yaml file.
