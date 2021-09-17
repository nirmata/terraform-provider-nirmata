---
page_title: "nirmata_catalog Resource"
---

# nirmata_catalog Resource

## Example Usage

```hcl

resource "nirmata_catalog" "tf-catalog-1" {
  name              = "tf-catalog"
  description       = ""
  labels            = {}
}

```

## Argument Reference

* `name` - (Required) a unique name for the catalog.
* `description` - (Optional) description of catalog.
* `labels` - (Optional) labels to set on the add-on application's environment.
