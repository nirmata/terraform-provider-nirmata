#terraform  nirmata provider 
terraform {
  required_providers {
    nirmata = {
      source = "nirmata/nirmata"
      version = "1.1.10-rc2"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "3.0.2"
    }
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }
  }
}