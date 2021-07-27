provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  
  #  token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as
  #  the environment variable NIRMATA_URL.
  
  #  url = ""
}


resource "nirmata_catalog" "tf-catalog-1" {
  name              = "tf-cat"
  description       = ""
  labels            = {}
}
