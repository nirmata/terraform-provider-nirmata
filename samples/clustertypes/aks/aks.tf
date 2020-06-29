provider "nirmata" {
}
resource "nirmata_aks_clusterType" "aks-cluster-type" {
  name       = "tf-akstype"
  version  = "" //The version of Kubernetes that should be used for this cluster.
  credentials = "" //  cloud credentials that hosts this cluster
  region= "" //The Azure region into which the cluster should be deployed
  resourcegroup= ""   //A resource group is a collection of resources that share the same lifecycle, permissions, and policies.
  subnetid= ""
  vmsize= ""  //computing, memory, networking, or storage needs
  vmsettype= "" 
  httpsapplicationrouting= true
  monitoring= false
  workspaceid = "" // Log Analytics workspace to store monitoring data.
  disksize = 30
}

resource "nirmata_ProviderManaged_cluster" "cluster-1" {
  name       = "tf-akscluster"
  type_selector  =  nirmata_aks_clusterType.aks-cluster-type.name
  node_count = 1
}

output "name" {
  description = "Cluster type name"
  value       = nirmata_aks_clusterType.aks-cluster-type.name
}


