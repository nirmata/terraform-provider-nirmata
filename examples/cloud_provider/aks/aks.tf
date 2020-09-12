provider "nirmata" {
  // Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  // token = ""

  // Nirmata address. Defaults to https://nirmata.io and can be configured as the environment variable NIRMATA_URL.
  // url = ""
}

// A ClusterType contains reusable configuration to create clusters.
resource "nirmata_aks_clusterType" "aks-cluster-type" {

  // a unique name for the cluster type (e.g. az-cluster)
  // Required
  name = var.name

  // the Kubernetes version (e.g. 1.17.7)
  // Required
  version = var.version

  // the Azure cloud credentials name configured in Nirmata (e.g. azure-credentials)
  // Required
  // credentials = var.credentials

  // the Azure region into which the cluster should be deployed
  // Required
  region = var.region

  // the Azure resource group
  // Required
  // resource_group = var.resource_group

  // the Azure subnet ID to use for the NodePool. The ID is a long path like this:
  // "/subscriptions/{uuid}/resourceGroups/{name}/providers/Microsoft.Network/virtualNetworks/{name}/subnets/default"
  // Required
  // subnet_id = var.subnet_id

  // the Azure VM size to use (e.g. Standard_D2_v3)
  // Required
  vm_size = var.vm_size

  // the VM set type (VirtualMachineScaleSets or AvailabilitySets)
  // Required
  vm_set_type = var.vm_set_type

  // the worker node disk size in GB
  // Required
  disk_size = var.disk_size

  // enable HTTPS Application Routing
  // Optional
  https_application_routing = var.https_application_routing

  // enable container monitoring
  // Optional
  monitoring = var.monitoring

  // the workspace ID to store monitoring data
  // Optional
  workspace_id = var.workspace_id
}

// A Cluster is created using a ClusterType
resource "nirmata_ProviderManaged_cluster" "aks-cluster" {

  // a unique name for the Cluster
  // Required
  name = var.cluster_name

  // the cluster type
  // Required
  type_selector = nirmata_aks_clusterType.aks-cluster-type.name

  // number of worker nodes
  // Required
  node_count = var.node_count
}

