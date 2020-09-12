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
| credentials | the Azure cloud credentials name configured in Nirmata (e.g. azure-credentials) Required | `any` | n/a | yes |
| disk_size | the worker node disk size in GB Required | `any` | n/a | yes |
| https_application_routing | enable HTTPS Application Routing Optional | `any` | n/a | yes |
| monitoring | enable container monitoring Optional | `any` | n/a | yes |
| name | a unique name for the cluster type (e.g. az-cluster) Required | `any` | n/a | yes |
| node_count | number of worker nodes Required | `any` | n/a | yes |
| region | the Azure region into which the cluster should be deployed Required | `any` | n/a | yes |
| resource_group | the Azure resource group Required | `any` | n/a | yes |
| subnet_id | the Azure subnet ID to use for the NodePool. The ID is a long path like this: "/subscriptions/{uuid}/resourceGroups/{name}/providers/Microsoft.Network/virtualNetworks/{name}/subnets/default" Required | `any` | n/a | yes |
| unquoted | n/a | `any` | n/a | yes |
| version | the Kubernetes version (e.g. 1.17.7) Required | `any` | n/a | yes |
| vm_set_type | the VM set type (VirtualMachineScaleSets or AvailabilitySets) Required | `any` | n/a | yes |
| vm_size | the Azure VM size to use (e.g. Standard_D2_v3) Required | `any` | n/a | yes |
| workspace_id | the workspace ID to store monitoring data Optional | `any` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| cluster_name | Cluster name |
| cluster_type_name | ClusterType name |

