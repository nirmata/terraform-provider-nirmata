provider "nirmata" {
  // Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  //
  // token = ""

  // Nirmata address. Defaults to https://nirmata.io and can be configured as
  // the environment variable NIRMATA_URL.
  //
  // url = ""
}

// An Environment Resource Type is used while run application

resource "nirmata_environment_type" "tf-env-type-1" {
  name             = "tf-env"
  resource_limits  = {
      cpu = "xx",
      memory = "xx"
      pod = "xx"
      storage = "xx"
  }
  is_default      = false
}
