provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  #  token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as the environment variable NIRMATA_URL.
  #  url = ""
}


resource "nirmata_cluster_addons" "cluster_addon" {
  name                       = "addon1"
  cluster                    = "my-cluster-1"
  application                = "cert-manager"
  environment                = "cert-manager"
  catalog                    = "default-addon-catalog"
  namespace                  = "cert-manager"
  channel                    = "Rapid"
  # labels                     = {cluster = "addon"}
  # service_name               ="service_name"
  # service_scheme             ="service_scheme"
  # service_port               = "service_port"
}