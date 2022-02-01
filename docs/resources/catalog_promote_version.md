---
page_title: "nirmata_promote_version Resource"
---

# nirmata_promote_version Resource

 Promote release version .

## Example Usage

```hcl

resource "nirmata_promote_version" "tf-catalog-promote-version" {
  rollout_name        = "tf-version"
  catalog             = "test-catalog"
  application         = "test-application"
  version             = "version"
  channel             = "Rapid"
 }


```

## Argument Reference

* `rollout_name` - (Required) A unique name for rollout.
* `catalog` - (Required) the name of catalog.
* `application` - (Required) the application name.
* `channel` - (Required) The channel from which the application should be deployed.
* `version` - (Required)  the version for the application.