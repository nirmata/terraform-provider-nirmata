provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  
  #  token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as
  #  the environment variable NIRMATA_URL.
  
  #  url = ""
}

resource "nirmata_run_application" "tf-catalog-run-app" {
  name                = "tf-run-app"
  application         = "application-name"
  catalog             = "catlog-name"
  ## version          = "version-name"
  channel             = "Rapid"
  environments        = ["env1","env2"]
 }
