## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| nirmata | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| credentials | Cloud credentials to use for this cluster Required | `any` | n/a | yes |
| machine_type | The type of VM for worker nodes Required | `any` | n/a | yes |
| name | The name of this cluster Required | `any` | n/a | yes |
| region | The region into which the cluster should be deployed Required | `any` | n/a | yes |
| unquoted | n/a | `any` | n/a | yes |
| version | The version of Kubernetes that should be used for this cluster. Required | `any` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| cluster_name | Cluster name |
| cluster_type_name | ClusterType name |

