provider "nirmata" {
  // Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  //
  // token = ""

  // Nirmata address. Defaults to https://nirmata.io and can be configured as
  // the environment variable NIRMATA_URL.
  //
  // url = ""
}

// An Environment is used to run application
resource "nirmata_environment" "tf-env-1" {
  name        = "tf-env-1"
  type        = "medium"
  cluster     = "prod-demo"
  namespace   = "tf-ns-1"
}
