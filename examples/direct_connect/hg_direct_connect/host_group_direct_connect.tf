
provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
}

resource "nirmata_host_group_direct_connect" "dc-host-group" {
  name = var.name
}

// You will likely want to install the Nirmata agent via
// "${nirmata_host_group_direct_connect.dc-host-group.curl_script}"
// before creating the cluster.

resource "nirmata_cluster_direct_connect" "dc-cluster-1" {
  name       = var.direct_connect_cluster_name
  policy     = var.policy
  host_group = nirmata_host_group_direct_connect.dc-host-group.name
}

