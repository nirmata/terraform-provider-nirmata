provider "nirmata" {
}

resource "nirmata_host_group_direct_connect" "dc-gh-1" {
    name = "dc-gh-1"
}

resource "nirmata_cluster_gke" "dc-gke-1" {
    name = "dc-gh-1"
    disk_size = 60
    node_type = "n1-standard-2"
    node_count = 1
    region = "us-central1-a"
}