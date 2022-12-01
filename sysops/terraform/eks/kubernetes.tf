provider "kubernetes" {
  host = aws_eks_cluster.data-manager-eks-cluster.endpoint
}
