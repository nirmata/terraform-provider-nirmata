output "cluster_type_name" {
  description = "ClusterType name"
  value       = nirmata_eks_clusterType.eks-cluster-type.name
}

output "cluster_name" {
  description = "Cluster name"
  value       = nirmata_ProviderManaged_cluster.eks-cluster.name
}