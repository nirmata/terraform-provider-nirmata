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
    name             = "gke-vaults"
    path             = "nirmata/$(cluster.nyt)"
    addon_name       = "vault-test"
    credentials_id   = "a30ab89f-cbb6-488d-96f4-86d5b6439be9"
    delete_auth_path = false

    roles {
      name                 = "nginx"
      service_account_name = "default"
      namespace            = "nginx"
      policies             = "get-secrets"
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
