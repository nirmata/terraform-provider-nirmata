provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  
  #  token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as
  #  the environment variable NIRMATA_URL.
  
  #  url = ""
}

resource "nirmata_helm_application" "tf-helm-git" {
  name                = "tf-helm-app"
  repository          = "Bitnami"
  application         = "airflow"
  location            = "https://charts.bitnami.com/bitnami/airflow-0.0.1.tgz"
  app_version         = "1.10.3"
  chart_version       = "0.0.1"
  catalog             = ""

}
