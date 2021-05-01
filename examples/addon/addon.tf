provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  #  token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as the environment variable NIRMATA_URL.
  #  url = ""
}


resource "nirmata_cluster_addons" "cluster_addon" {
  name                       = "addon1"
  cluster                    = "cluster-1"
  catalog                    = "default-addon-catalog"
  application                = "app"
  environment                = "env"
  namespace                  = "ns"
  channel                    = "channel"
  # labels                     = {cluster = "addon"}
}