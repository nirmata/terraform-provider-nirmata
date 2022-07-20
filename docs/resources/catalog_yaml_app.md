---
page_title: "nirmata_catalog_application Resource"
---

# nirmata_catalog_application Resource

 Application is a group of workloads, routing and storage configurations.

## Example Usage

```hcl

resource "nirmata_catalog_application" "tf-catalog-app" {
  name              = "tf-catalog-app"
  catalog           = "test-catalog"
  yamls             = file("${path.module}/fo.yaml")
}

```

## Argument Reference

* `name` - (Required) Enter a unique name for the application in the catalog.
* `catalog` - (Required) Enter the name of the catalog.
* `yamls` - (Required)  Enter the path for the yaml file.
