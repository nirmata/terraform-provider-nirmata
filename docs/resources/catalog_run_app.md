---
page_title: "nirmata_run_application Resource"
---

# nirmata_run_application Resource

Deploy an application in environments.

## Example Usage

```hcl

resource "nirmata_run_application" "tf-catalog-run-app" {
  name                = "tf-run-app"
  application         = "application-name"
  catalog             = "catlog-name"
  version             = "version-name"
  channel             = "Rapid"
  environments        = ["env1","env2"]
 }

```

## Argument Reference

* `name` - (Required) A unique name to identify your application.
* `catalog` - (Required) the name of catalog.
* `application` - (Required) the application name.
* `channel` - (Required) The channel from which the application should be deployed.
* `environments` - (Required) the list of environments to deploy an application .
* `version` - (Optional)  the version for the application.