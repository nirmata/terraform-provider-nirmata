provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  
   token = "L0DIoi44ma1FFDjNVqEtdZaZTprTlYDrxbDvY1CiElSiFLbaVBySRth/aiZcGS25BZpKf3kI6Cee3lYvEhcMaA=="

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as
  #  the environment variable NIRMATA_URL.
  
   url = "https://devtest5.nirmata.co/"
}

resource "nirmata_run_application" "tf-catalog-run-app" {
  name                = "tf-run-app-new"
  application         = "test-3"
  catalog             = "default-policy-catalog"
  version             = "auto-0.0.1-6f28035"
  channel             = "Rapid"
  environments        = ["tf-env", "tf-env-1"]
 }
