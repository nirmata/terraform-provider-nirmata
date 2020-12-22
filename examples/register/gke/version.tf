terraform {
  required_version = ">= 0.13"
  required_providers {
      nirmata = {
      source  = "registry.terraform.io/nirmata/nirmata"
      version = "1.0.0"
    }
    
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }
  }
}