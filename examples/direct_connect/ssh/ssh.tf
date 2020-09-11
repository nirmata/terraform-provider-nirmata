// Honestly this example exists mainly to demostrate using the remote-exec 
// provisioner to install the Nirmata agent. To use something like vsphere,
// digitalocean, or the like replace null_resource with your cloud provider 
// and include the depends_on, provisioner and connection.

provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
  // NIRMATA_URL=https://nirmata.local terraform <whatever>
}

resource "nirmata_host_group_direct_connect" "dc-host-group" {
  // This must not exist in Nirmata!
  name = "baremetal-hg-1"
}
// This is fake resource to run the provisioner.
// Noter that Terraform really isn't designed to do this, but it does work.
// Destroying this resource will NOT cleanup the node or remove.
resource "null_resource" "node" {
  // You need this to be sure curl_script variable exists.
  depends_on = [nirmata_host_group_direct_connect.dc-host-group]

  // This is for a modern Ubuntu.  See Nirmata docs for other distros. The last 
  // part of the inline is distro independent assuming things are setup correctly.
  provisioner "remote-exec" {
    inline = ["sudo apt-get update",
      "echo sudo apt-get install -y docker.io",
    "${nirmata_host_group_direct_connect.dc-host-group.curl_script}"]
  }
  connection {
    type = "ssh"
    user = "ubuntu"
    // Use a password or key to access node
    //password    = "pass1234"
    private_key = file("~/.ssh/terraform")
    // In a real resource remove the host or replace it with a reference to the node(s).
    host = "10.18.0.12"
  }

  // You can add more provisioners to target more nodes
  // Or with a real cloud provider resource just add a count parameter.

  // If you are serious about using the null_resource to manage servers consider
  // using a destroy time provisioner to clean up.
  // https://www.terraform.io/docs/provisioners/index.html#destroy-time-provisioners
  // https://github.com/nirmata/custom-scripts/blob/master/cleanup-cluster-agent.sh
}

resource "nirmata_cluster_direct_connect" "dc-cluster-1" {
  // This cluster must not exist in Nirmata.
  name = "baremetal-cluster-1"
  // This policy must exist in Nirmata.
  policy     = "default-v1.16.0"
  host_group = nirmata_host_group_direct_connect.dc-host-group.name
  // This depends must match the cloud provider resource.
  depends_on = [null_resource.node]
}
