provider "kubernetes" {
  host = aws_eks_cluster.data_manager_eks_cluster.endpoint
}
