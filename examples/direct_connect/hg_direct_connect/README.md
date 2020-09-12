## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| nirmata | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| direct_connect_cluster_name | The name of this cluster Required | `any` | n/a | yes |
| machine_type | The type of VM for worker nodes Required | `any` | n/a | yes |
| name | The name of this cluster Required | `any` | n/a | yes |
| policy | policy name Required | `any` | n/a | yes |
| region | The region into which the cluster should be deployed Required | `any` | n/a | yes |
| unquoted | n/a | `any` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| agent_script | Nirmata agent install command |

