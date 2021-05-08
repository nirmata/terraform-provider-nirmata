provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  
  # token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as
  #  the environment variable NIRMATA_URL.
  
  # url = ""
}

#  An registered Cluster type is used while creating registered clusters

resource "nirmata_cluster_type_registered" "tf-registered-type-1" {
  name  = "tf-registered-cluster"
  cloud = "GoogleCloudPlatform"

  labels = {
    foo = "bar"
  }

  vault_auth {
    name             = ""
    path             = ""
    addon_name       = ""
    credentials_id   = ""
    delete_auth_path = false

    roles {
      name                 = "nginx"
      service_account_name = "default"
      namespace            = "nginx"
      policies             = ""
    }
  }

  addons {
    name            = "cert-manager"
    addon_selector  = "cert-manager"
    catalog         = "default-catalog"
    channel         = "Stable"
    sequence_number = 2
  }

}
