provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  
  #  token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as
  #  the environment variable NIRMATA_URL.
  
  #  url = ""
}

#  An Environment is used to run application
resource "nirmata_environment" "tf-env" {
  name        = "environment-name"
  type        = "medium"
  cluster     = "cluster-name"
  namespace   = "namespace-name"
  environment_update_action   = "update" 
  labels      = { foo = "bar"}
  
}
