provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
  // NIRMATA_URL=https://mynirmata.local terraform <whatever>
}

resource "nirmata_host_group_direct_connect" "dc-host-group" {
  // this must not existing in Nirmata
  name = "baremetal-hg-1"
}
// This is fake resource to run the provisioner
// terraform really isn't designed to do this.
resource "null_resource" "node" {
    depends_on = [ nirmata_host_group_direct_connect.dc-host-group ]
    // This is for a modern Ubuntu see Nirmata doc for other distros
    provisioner "remote-exec" {
      inline = [ "sudo apt-get update",
        "echo sudo apt-get install -y docker.io",
        "${nirmata_host_group_direct_connect.dc-host-group.curl_script}"]
    }

  connection {
    type        = "ssh"
    user        = "ubuntu"
    // Use a password or key to access node
    //password    = "pass1234"
    private_key = file("~/.ssh/terraform")
    host        = "10.18.0.12"

  //You can add more provisioner to target more nodes
  }
}

resource "nirmata_cluster_direct_connect" "dc-cluster-1" {
  name = "baremetal-cluster-1"
  // This policy must exist in nirmata
  policy = "default-v1.16.0"
  host_group = nirmata_host_group_direct_connect.dc-host-group.name
  depends_on = [ null_resource.node ]
}
