variable "cluster-name" {
  default = "data-manager-eks"
  type    = string
}

resource "aws_iam_role" "data-manager-eks-cluster-iam" {
  name = "data-manager-eks-cluster-iam"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "eks.amazonaws.com"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "data-manager-eks-cluster-iam-AmazonEKSClusterPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.data-manager-eks-cluster-iam.name
}

resource "aws_iam_role_policy_attachment" "data-manager-eks-cluster-iam-AmazonEKSServicePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSServicePolicy"
  role       = aws_iam_role.data-manager-eks-cluster-iam.name
}

resource "aws_security_group" "data-manager-eks-cluster-sg" {
  name        = "data-manager-eks-cluster"
  description = "Cluster communication with worker nodes"
  vpc_id      = aws_vpc.data-manager-eks-vpc.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "data-manager-eks-cluster"
  }
}

resource "aws_eks_cluster" "data-manager-eks-cluster" {
  enabled_cluster_log_types = ["api", "audit", "authenticator", "controllerManager", "scheduler"]
  name     = var.cluster-name
  role_arn = aws_iam_role.data-manager-eks-cluster-iam.arn

  vpc_config {
    security_group_ids = [aws_security_group.data-manager-eks-cluster-sg.id]
    subnet_ids         = aws_subnet.data-manager-eks-subnet.*.id
  }

  depends_on = [
    aws_iam_role_policy_attachment.data-manager-eks-cluster-iam-AmazonEKSClusterPolicy,
    aws_iam_role_policy_attachment.data-manager-eks-cluster-iam-AmazonEKSServicePolicy,
  ]
}

locals {
  kubeconfig = <<-KUBECONFIG
    ---
    apiVersion: v1
    clusters:
    - cluster:
        server: ${aws_eks_cluster.data-manager-eks-cluster.endpoint}
        certificate-authority-data: ${aws_eks_cluster.data-manager-eks-cluster.certificate_authority.0.data}
      name: eks-test
    contexts:
    - context:
        cluster: eks-test
        user: aws
      name: eks-test
    current-context: eks-test
    kind: Config
    preferences: {}
    users:
    - name: aws
      user:
        exec:
          apiVersion: client.authentication.k8s.io/v1beta1
          command: aws-iam-authenticator
          args:
            - "token"
            - "-i"
            - "${var.cluster-name}"
  KUBECONFIG
}

output "kubeconfig" {
  value = "${local.kubeconfig}"
}

