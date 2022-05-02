#terraform  nirmata provider 
terraform {
  required_providers {

    nirmata = {
      source = "nirmata/nirmata"
      version = "1.1.7-rc8"
    }

    aws = {
        source  = "hashicorp/aws"
        version = ">= 3.20.0"
      }
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }
  }
}