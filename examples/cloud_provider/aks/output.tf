output "cluster_type_name" {
  description = "ClusterType name"
  value       = nirmata_aks_clusterType.aks-cluster-type.name
}

output "cluster_name" {
  description = "Cluster name"
  value       = nirmata_ProviderManaged_cluster.aks-cluster.name
}
