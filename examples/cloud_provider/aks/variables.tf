variable unquoted {}

// a unique name for the cluster type (e.g. az-cluster)
// Required
variable "name" {
  name = ""
}

// the Kubernetes version (e.g. 1.17.7)
// Required
variable "version" {
  version = "1.17.7"
}

// the Azure cloud credentials name configured in Nirmata (e.g. azure-credentials)
// Required
variable "credentials" {
  credentials = ""
}

// the Azure region into which the cluster should be deployed
// Required
variable "region" {
  region = ""
}

//the Azure resource group
// Required
variable "resource_group" {
  resource_group = ""
}

// the Azure subnet ID to use for the NodePool. The ID is a long path like this:
// "/subscriptions/{uuid}/resourceGroups/{name}/providers/Microsoft.Network/virtualNetworks/{name}/subnets/default"
// Required
variable "subnet_id" {
  subnet_id = ""
}

// the Azure VM size to use (e.g. Standard_D2_v3)
// Required
variable "vm_size" {
  vm_size = "Standard_D2_v3"
}

// the VM set type (VirtualMachineScaleSets or AvailabilitySets)
// Required
variable "vm_set_type" {
  vm_set_type = "VirtualMachineScaleSets"
}

// the worker node disk size in GB
// Required
variable "disk_size" {
  disk_size = 60
}

// enable HTTPS Application Routing
// Optional
variable "https_application_routing" {
  https_application_routing = false
}

// enable container monitoring
// Optional
variable "monitoring" {
  monitoring = false
}

// the workspace ID to store monitoring data
// Optional
variable "workspace_id" {
  workspace_id = ""
}

// a unique name for the Cluster
// Required
variable "cluster_name" {
  cluster_name = "tf-akscluster"
}


// number of worker nodes
// Required
variable "node_count" {
  node_count = 1
}
