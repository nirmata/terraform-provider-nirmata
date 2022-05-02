terraform {
  required_providers {
    nirmata = {
      source = "nirmata/nirmata"
      version = "1.1.7-rc8"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }
  }
}