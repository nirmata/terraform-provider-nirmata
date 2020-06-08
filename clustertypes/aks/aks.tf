provider "nirmata" {
  // Set NIRMATA_TOKEN with your API Key
  // You can also set NIRMATA_URL with the Nirmata URL address
}
resource "nirmata_aks_clusterType" "cluster-type-aks" {
  name       = "aks-cluster-type"
  version  = "1.14.7"
  credentials = "take-azure"
  region= "centralus"
  resourcegroup= "nirmata-test"
  subnetid= "/subscriptions/baf89069-e8f3-46f8-b74e-c146931ce7a4/resourceGroups/nirmata-test/providers/Microsoft.Network/virtualNetworks/azure-test/subnets/default"
  vmsize= "Standard_DS2_v2"
  vmsettype= "AvailabilitySet"
  httpsapplicationrouting= true
  monitoring= false
  workspaceid = ""
  disksize = 30
}