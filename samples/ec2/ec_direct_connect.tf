
variable "awsprops" {
    type = map
    default = {
    region = "us-west-1"
    vpc = "vpc-00012345678909876"
    // This is the AMI for ubuntu in us-west-1 ami.  Note that images are region specific
    ami = "ami-03ba3948f6c37a4b0"
    // t3a.medium is fine for testing for production consider m5a.xlarge or m5.xlarge
    itype = "t3a.medium"
    subnet = "subnet-12345678909876543"
    publicip = true
    // Must exist
    keyname = "terraform-test-west-1"
    // Must not exist
    secgroupname = "terraform-test"
    instance_count = 3
  }
}

provider "aws" {
  region = lookup(var.awsprops, "region")
}

provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
}

resource "nirmata_host_group_direct_connect" "dc-host-group" {
  name = "sam-hg-1"
}

resource "aws_security_group" "nirmata-dc-sg" {
  name = lookup(var.awsprops, "secgroupname")
  description = lookup(var.awsprops, "secgroupname")
  vpc_id = lookup(var.awsprops, "vpc")

  // To Allow SSH Transport
  // Disable in production?
  ingress {
    from_port = 22
    protocol = "tcp"
    to_port = 22
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }
}


resource "aws_instance" "nirmata-dc" {
  count = lookup(var.awsprops, "instance_count")
  ami = lookup(var.awsprops, "ami")
  instance_type = lookup(var.awsprops, "itype")
  subnet_id = lookup(var.awsprops, "subnet") #FFXsubnet2
  associate_public_ip_address = lookup(var.awsprops, "publicip")
  key_name = lookup(var.awsprops, "keyname")
 
  // We are using remote-exec because we can't block on user data in AWS
  // This is for modern Ubuntu versions see Nirmata docs for various Linux distros
  provisioner "remote-exec" {
    inline = [ "sudo apt-get update", 
      "sudo apt-get install -y docker.io",
      "${nirmata_host_group_direct_connect.dc-host-group.curl_script}"]
  }
  
  connection {
    type        = "ssh"
    user        = "ubuntu"
    private_key = file("~/.ssh/terraform")
    host        = self.public_ip
  }

  vpc_security_group_ids = [
    aws_security_group.nirmata-dc-sg.id
  ]
  root_block_device {
    delete_on_termination = true
    iops = 150
    volume_size = 100
    volume_type = "gp2"
  }
  tags = {
    Name ="nirmata-dc"
    Environment = "DEV"
    OS = "UBUNTU"
  }

  depends_on = [ aws_security_group.nirmata-dc-sg ]
}

resource "nirmata_cluster_direct_connect" "dc-cluster-1" {
  name = "sam-cluster-1"
  policy = "default-v1.16.0"
  host_group = nirmata_host_group_direct_connect.dc-host-group.name
  depends_on = [ aws_instance.nirmata-dc ]
}


