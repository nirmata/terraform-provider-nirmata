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

* `name` - (Required) a unique name for the application in catalog.
* `application` - (Required) the application name.
* `repository` - (Required)  the repository URL.
* `location` - (Required). the location of the chart. (https://charts.bitnami.com/bitnami/airflow-0.0.1.tgz)
* `app_version` - (Required). specify the version of the application. ("0.0.1")
* `chart_version` - (Required).the version of the chart. ("1.2.3")