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
| credentials | the GCP cloud credentials name configured in Nirmata (e.g. gcp-credentials) Required | `any` | n/a | yes |
| disk_size | the worker node disk size in GB Required | `any` | n/a | yes |
| machine_type | the GCP machine type (e.g. "e2-standard-2") Required | `any` | n/a | yes |
| name | a unique name for the cluster type (e.g. az-cluster) Required | `any` | n/a | yes |
| node_count | number of worker nodes Required | `any` | n/a | yes |
| region | the GCP region into which the cluster should be deployed (e.g. "us-central1-b") Required | `any` | n/a | yes |
| unquoted | n/a | `any` | n/a | yes |
| version | a unique name for the cluster type (e.g. eks-cluster) Required | `any` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| cluster_name | Cluster name |
| cluster_type_name | ClusterType name |

