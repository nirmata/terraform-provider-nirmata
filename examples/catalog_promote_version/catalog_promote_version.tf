provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  
  #  token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as
  #  the environment variable NIRMATA_URL.
  
  #  url = ""
}

resource "nirmata_promote_version" "tf-catalog-promote-version" {
  rollout_name        = "tf-version"
  catalog             = ""
  application         = ""
  version             = ""
  channel             = "Rapid"
 }
