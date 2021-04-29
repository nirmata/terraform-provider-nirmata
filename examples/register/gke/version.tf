terraform {
  required_version = ">= 0.14"
  required_providers {
    nirmata = {
      source  = "nirmata/nirmata"
      version = "1.0.0-pre"
    }

    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }
  }
}
