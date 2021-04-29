---
page_title: "nirmata_cluster_type_eks Resource"
---

# nirmata_cluster_type_gke Resource

An Amazon Elastic Kubernetes Service (EKS) cluster type. A cluster type can be used to create multiple clusters.

## Example Usage

```hcl

// 1. Create a EKS cluster type
resource "nirmata_cluster_type_eks" "eks-cluster-type-1" {
  name                      = "tf-eks-cluster-type-1"
  version                   = "1.18"
  credentials               = "aws-xxxxx"
  region                    = "us-west-2"
  vpc_id                    = "vpc-xxxxxxxx"
  subnet_id                 = ["subnet-xxxxxxxx", "subnet-xxxxxxxx"]
  security_groups           = ["sg-xxxxxxxxxxxxxxxx"]
  cluster_role_arn          = "arn:aws:iam::xxxxxxxx:role/xxxxxxxx"
  enable_private_endpoint   = true
  enable_identity_provider  = true
  auto_sync_namespaces       = false

  nodepools {
    name                = "default"
    instance_type       = "t3.medium"
    disk_size           = 60
    ssh_key_name        = "xxxxxxxx"
    security_groups     = ["sg-xxxxxxxxxxxxxxxx"]
    iam_role            = "arn:aws:iam::xxxxxxxx:role/eks-xxxxxxxx"
  }
}

// 2. Create a nirmata_cluster using the cluster_type
resource "nirmata_cluster" "eks-cluster-1" {
  name                 = "eks-cluster-1"
  cluster_type         = nirmata_cluster_type_eks.eks-cluster-type-1.name
  node_count           = 1
}


```

## Argument Reference

* `name` - (Required) a unique name for the cluster.
* `version` - (Required) the EKS version (e.g. 1.18)
* `credentials` - (Required) the cloud credentials to use.
* `region` - (Required) the AWS region for the cluster.
* `vpc_id` - (Required) the AWS VPC ID for the cluster.
* `subnet_id` - (Required) a list of AWS VPC subnets to use for the cluster.
* `security_groups` - (Required) a list of AWS VPC security groups to use for the cluster.
* `cluster_role_arn` - (Required) the cluster role ARN.
* `enable_private_endpoint` - (Optional) specify if the cluster API seerver endpoint should be in a private network.
* `enable_identity_provider` - (Optional) enable IAM roles for service accounts.
* `auto_sync_namespaces` - (Optional) enable automatic synchronization of cluster namespaces to Nirmata.
* `enable_secrets_encryption` - (Optional) enable encryption at rest for secrets.
* `kms_key_arn` - (Optional) the KMS key ARN to use for secrets encryption.
* `log_types` - (Optional) the log types to collect.
* `enable_fargate` - (Optional) enable Fargate to provision nodes based on workload resource requests.
* `system_metadata` - (Optional) key-value pairs that will be provisioned as a ConfigMap called system-metadata-map in the cluster.
* `allow_override_credentials` - (Optional) Allow passing in cloud credentials when a cluster is created using this cluster type.
* `cluster_field_override` - (Optional) Allow override of cluster settings ('network' and 'subnetwork') when a cluster is created using this cluster type.
* `nodepool_field_override` - (Optional)  Allow override of node fields ('disk_size' and 'machine_type') when a cluster is created using this cluster type.
* `nodepools` - (Optional) A list of [nodepool](#nodepool) types.
* `addons` - (Optional) a list of add-on services.
* `vault_auth` - (Optional) vault authentication configuration.

## Nested Blocks

### nodepool

* `name` - (Required) a unique name for the node pool.
* `instance_type` - (Required) the EC2 instance type (e.g. "t3.medium").
* `disk_size` - (Required) the worker node disk size in GB (e.g. 60).
* `ssh_key_name` - (Required) the SSK key pair to access nodes.
* `security_groups` - (Required) the Node security groups.  
* `iam_role` - (Required) the IAM role to use for nodes.
* `ami_type` - (Required) the EKS-optimized Amazon Machine Image (AMI) type to use for node images.
* `image_id` - (Optional) an Amazon Machine Image (AMI) IS to use for node images.
* `node_annotations` -  (Optional) Annotations to set on each node in this pool. This setting is permanent.
* `node_labels` - (Optional) Labels to set on each Kubernetes node in this node pool. This setting is permanent.

### addons

* `name` - (Required) a unique name for the add-on service
* `addon_selector` - (Required) the catalog application name
* `catalog` - (Required) the catalog name
* `channel` - (Required) the release channel
* `sequence_number` - (Optional) a sequence number to control installation order

### vault_auth

* `name` - (Required) a unique name
* `path` - (Required) the vault authentication path. The variable $(cluster.name) is allowed for uniquenes.
* `addon_name` - (Required) the associated Vault Agent Injector add-on
* `credentials_name` - (Required) the Vault credentials to use 
* `roles` - (Required) a list of application roles to configure for add-on services

#### roles

* `name` - (Required) a unique name
* `service_account_name` - (Required) the allowed service account name
* `namespace` - (Required) the allowed namespace
* `policies` - (Required) the applied policies


