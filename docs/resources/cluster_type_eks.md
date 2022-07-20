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
  addons {
    name            = "vault-agent-injector"
    addon_selector  = "vault-agent-injector"
    catalog         = "default-catalog"
    channel         = "Stable"
    sequence_number = 1
  }

  vault_auth {
    name             = "vault-auth"
    path             = "nirmata/$(cluster.name)"
    addon_name       = "vault-agent-injector"
    credentials_name = "vault_access"
    delete_auth_path = true

    roles {
      name                 = "sample-role"
      service_account_name = "application-sample-sa"
      namespace            = "application-sample-ns"
      policies             = "application-sample-policy"
    }
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

* `name` - (Required) Enter a unique name for the cluster.
* `version` - (Required) Enter the EKS version (example, 1.18).
* `credentials` - (Required) Enter the cloud credentials that is used.
* `region` - (Required) Enter the AWS region for the cluster.
* `vpc_id` - (Required) Enter the AWS VPC ID for the cluster.
* `subnet_id` - (Required) Enter a list of AWS VPC subnets to be used for the cluster.
* `security_groups` - (Required) Enter a list of AWS VPC security groups to be used for the cluster.
* `cluster_role_arn` - (Required) Enter the cluster role ARN.
* `enable_private_endpoint` - (Optional) This field indicates if the cluster API server endpoint should be in a private network.
* `enable_identity_provider` - (Optional) This value enables IAM roles for service accounts.
* `auto_sync_namespaces` - (Optional) This value enables automatic synchronization of cluster namespaces for Nirmata.
* `enable_secrets_encryption` - (Optional) This value enables encryption at REST for secrets.
* `kms_key_arn` - (Optional) This field indicates the KMS key ARN to be used for the secrets encryption.
* `log_types` - (Optional) This field indicates the log types to be collected.
* `enable_fargate` - (Optional) This value enables Fargate to provision nodes based on the workload resource requests.
* `system_metadata` - (Optional) This is the key-value pairs that will be provisioned as a ConfigMap called system-metadata-map in the cluster.
* `allow_override_credentials` - (Optional) This allows the passing of cloud credentials when a cluster is created using this cluster type.
* `cluster_field_override` - (Optional) This allows the override of cluster settings ('network' and 'subnetwork') when a cluster is created using this cluster type.
* `nodepool_field_override` - (Optional)  This allows the override of node fields ('disk_size' and 'machine_type') when a cluster is created using this cluster type.
* `nodepools` - (Optional) This field indicates a list of [nodepool](#nodepool) types.
* `addons` - (Optional) This field indicates a list of add-on services.
* `vault_auth` - (Optional) This field indicates the vault authentication configuration.

## Nested Blocks

### nodepool

* `name` - (Required) Enter a unique name for the node pool.
* `instance_type` - (Required) Enter the EC2 instance type (example, "t3.medium").
* `disk_size` - (Required) Enter the worker node disk size in GB (example, 60).
* `ssh_key_name` - (Required) Enter the SSK key pair to access nodes.
* `security_groups` - (Required) Enter the node security groups.  
* `iam_role` - (Required) Enter the IAM role to be used for nodes.
* `ami_type` - (Required) Enter the EKS-optimized Amazon Machine Image (AMI) type to be used for node images.
* `image_id` - (Optional) Enter the Amazon Machine Image (AMI) IS to be used for node images.
* `node_annotations` -  (Optional) This value indicates the annotations to be set on each node in this pool. This setting is permanent.
* `node_labels` - (Optional) This value indicates the labels to be set on each Kubernetes node in this node pool. This setting is permanent.

### addons

* `name` - (Required) Enter a unique name for the add-on service.
* `addon_selector` - (Required) Enter the catalog application name.
* `catalog` - (Required) Enter the catalog name.
* `channel` - (Required) Enter the channel from which the application should be deployed.
* `sequence_number` - (Optional) This feild indicates a sequence number that control installation order.

### vault_auth

* `name` - (Required) Enter a unique name for the vault authentication.
* `path` - (Required) Enter the vault authentication path. The variable $(cluster.name) is allowed for uniquenes.
* `addon_name` - (Required) Enter the associated Vault Agent Injector add-on.
* `credentials_name` - (Required) Enter the Vault credentials to be used for the authentication. 
* `roles` - (Required) Enter a list of application roles to be configured for the add-on services.
* `delete_auth_path` - (Optional) This field indicates the delete authentication path on cluster delete.

#### roles

* `name` - (Required) Enter a unique name for roles.
* `service_account_name` - (Required) Enter the allowed service account name.
* `namespace` - (Required) Enter the allowed namespace.
* `policies` - (Required) Enter the applied policies.


