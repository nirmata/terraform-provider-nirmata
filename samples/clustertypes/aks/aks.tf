provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
}
resource "nirmata_aks_clusterType" "cluster-type-aks" {
  name       = "aks-cluster-type"
  version  = "1.14.7" //The version of Kubernetes that should be used for this cluster.
  credentials = "" //  cloud credentials that hosts this cluster
  region= ""  //The Azure region into which the cluster should be deployed
  resourcegroup= ""   //A resource group is a collection of resources that share the same lifecycle, permissions, and policies.
  subnetid= ""
  vmsize= "Standard_DS2_v2"  //computing, memory, networking, or storage needs
  vmsettype= "AvailabilitySet" 
  httpsapplicationrouting= true
  monitoring= false
  workspaceid = "" // Log Analytics workspace to store monitoring data.
  disksize = 30
}