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

* `name` - (Required) Enter a unique name to identify your application.
* `catalog` - (Required) Enter the name of the catalog.
* `application` - (Required) Enter the application name.
* `channel` - (Required) Enter the channel from which the application should be deployed.
* `environments` - (Required) Enter the list of environments to deploy an application.
* `version` - (Optional)  This field indicates the version of the application.