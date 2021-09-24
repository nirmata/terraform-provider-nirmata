provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  
  # token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as
  #  the environment variable NIRMATA_URL.
  
  #  url = ""
}


resource "nirmata_catalog_application" "tf-catalog-app" {
  name              = "tf-catalog-app"
  catalog           = ""
  yamls             = file("${path.module}/foo.yaml")
}

output "version" {
  value = nirmata_catalog_application.tf-catalog-app.version
}

output "release_name" {
  value = nirmata_catalog_application.tf-catalog-app.release_name
}