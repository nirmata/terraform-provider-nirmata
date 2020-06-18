
provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
}

resource "nirmata_host_group_direct_connect" "dc-host-group" {
  name = "dc-hg-1"
}

// You will likely want to install the Nirmata agent via
// "${nirmata_host_group_direct_connect.dc-host-group.curl_script}"
// before creating the cluster.

resource "nirmata_cluster_direct_connect" "dc-cluster-1" {
  name = "dc-cluster-1"
  policy = "default-v1.16.0"
  host_group = nirmata_host_group_direct_connect.dc-host-group.name
}

output "agent_script" {
  description = "Nirmata agent install command"
  value       = nirmata_host_group_direct_connect.dc-host-group.curl_script
}
