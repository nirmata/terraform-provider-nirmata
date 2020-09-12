## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| nirmata | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| cluster_name | a unique name for the Cluster Required | `any` | n/a | yes |
| cluster_role_arn | the AWS cluster role ARN (e.g. "arn:aws:iam::000000007:role/sample") Required | `any` | n/a | yes |
| credentials | the AWS cloud credentials name configured in Nirmata (e.g. aws-credentials) Required | `any` | n/a | yes |
| disk_size | the worker node disk size Required | `any` | n/a | yes |
| instance_types | the AWS instance type for worker nodes (e.g. "t3.medium") Required | `any` | n/a | yes |
| key_name | the AWS SSH key name (e.g. ssh-keys) Optional | `any` | n/a | yes |
| name | a unique name for the cluster type (e.g. eks-cluster) Required | `any` | n/a | yes |
| node_count | number of worker nodes Required | `any` | n/a | yes |
| node_iam_role | the AWS IAM role for worker nodes (e.g. "arn:aws:iam::000000007:role/sample") Required | `any` | n/a | yes |
| node_security_groups | the AWS security group for worker node firewalling (e.g. ["sg-028208181hh110"]) Required | `any` | n/a | yes |
| region | the AWS region into which the cluster should be deployed Required | `any` | n/a | yes |
| security_groups | the AWS security group for firewalling (e.g. ["sg-028208181hh110"]) Required | `any` | n/a | yes |
| subnet_id | the AWS VPC subnet ID in whi  ch the cluster should be provisioned (e.g. ["subnet-e8b1a2k1j", "subnet-ey907f5v"]) Required | `any` | n/a | yes |
| unquoted | n/a | `any` | n/a | yes |
| version | the Kubernetes version (e.g. 1.16) Required | `any` | n/a | yes |
| vpc_id | the AWS VPC subnet ID in which the cluster should be provisioned Required | `any` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| cluster_name | Cluster name |
| cluster_type_name | ClusterType name |

