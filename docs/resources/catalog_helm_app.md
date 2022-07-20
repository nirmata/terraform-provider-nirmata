---
page_title: "nirmata_helm_application Resource"
---

# nirmata_helm_application Resource

 Application is a group of workloads, routing and storage configurations.

## Example Usage

```hcl

resource "nirmata_helm_application" "tf-helm-git" {
  name                = "tf-helm-app"
  repository          = "Bitnami"
  application         = "airflow"
  location            = "https://charts.bitnami.com/bitnami/airflow-0.0.1.tgz"
  app_version         = "1.10.3"
  chart_version       = "0.0.1"
  catalog             = "test-catalog"

}


```

## Argument Reference

* `name` - (Required) Enter a unique name for the application in the catalog.
* `application` - (Required) Enter the application name.
* `repository` - (Required)  Enter the repository URL.
* `location` - (Required). Enter the location of the chart (for example, https://charts.bitnami.com/bitnami/airflow-0.0.1.tgz).
* `app_version` - (Required). Enter the version of the application (for example, "0.0.1").
* `chart_version` - (Required).Enter the version of the chart (for example, "1.2.3").