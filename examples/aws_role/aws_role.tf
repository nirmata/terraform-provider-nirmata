provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
    token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as the environment variable NIRMATA_URL.
    url = ""
}


resource "nirmata_aws_role_credentials" "aws_role" {
  name                      = ""
  region                    = ""
  aws_access_key_id         = ""
  aws_secret_key            = ""
  # aws_role_arn              = ""
  # aws_external_id           = ""
}