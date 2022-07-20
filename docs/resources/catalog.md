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

* `name` - (Required) Enter a unique name for the catalog.
* `description` - (Optional) This field indicates the description for the catalog.
* `labels` - (Optional) This field indicates the labels to be set on the catalog.
