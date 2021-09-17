---
page_title: "nirmata_run_application Resource"
---

# nirmata_run_application Resource

Deploy an application in environments.

## Example Usage

```hcl

resource "nirmata_run_application" "tf-catalog-run-app" {
  name                = "tf-run-app"
  catalog             = ""
  application         = ""
  version             = ""
  channel             = "Rapid"
  environments        = []
 }

```

## Argument Reference

* `name` - (Required) A unique name to identify your application.
* `catalog` - (Required) the name of catalog.
* `application` - (Required) the application name.
* `channel` - (Required) The channel from which the application should be deployed.
* `environments` - (Required) the list of environments to deploy an application .
* `version` - (Required)  the version for the application.