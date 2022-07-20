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

* `rollout_name` - (Required) Enter a unique name for rollout.
* `catalog` - (Required) Enter the name of the catalog.
* `application` - (Required) Enter the application name.
* `channel` - (Required) Enter the channel from which the application should be deployed.
* `version` - (Required)  Enter the version of the application.