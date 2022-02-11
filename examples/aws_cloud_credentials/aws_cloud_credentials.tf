provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  #  token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as the environment variable NIRMATA_URL.
   // url = ""
}


resource "nirmata_aws_cloud_credentials" "aws_cloud_credential" {
  name                      = "aws-credential"
  access_type               = "access_key"  // value are access_key or assume_role
  description               = "AWS Account"
  region                    = "us-west-1"
  access_key_id             = ""            // if access_type is access_key 
  secret_key                = ""            // if access_type is access_key
  # aws_role_arn            = ""            // if access_type is assume_role
}