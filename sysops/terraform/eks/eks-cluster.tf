variable "cluster-name" {
  default = "data-manager-eks"
  type    = string
}

resource "aws_iam_role" "data_manager_eks_cluster_iam" {
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

resource "aws_iam_role_policy_attachment" "data_manager_eks_cluster_iam_AmazonEKSClusterPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.data_manager_eks_cluster_iam.name
}

resource "aws_iam_role_policy_attachment" "data_manager_eks_cluster_iam_AmazonEKSServicePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSServicePolicy"
  role       = aws_iam_role.data_manager_eks_cluster_iam.name
}

resource "aws_security_group" "data_manager_eks_cluster_sg" {
  name        = "data_manager_eks_cluster"
  description = "Cluster communication with worker nodes"
  vpc_id      = aws_vpc.data-manager-eks-vpc.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "data_manager_eks_cluster"
  }
}

resource "aws_eks_cluster" "data_manager_eks_cluster" {
  enabled_cluster_log_types = ["api", "audit", "authenticator", "controllerManager", "scheduler"]
  name     = var.cluster-name
  role_arn = aws_iam_role.data_manager_eks_cluster_iam.arn

  vpc_config {
    security_group_ids = [aws_security_group.data_manager_eks_cluster_sg.id]
    subnet_ids         = aws_subnet.data-manager-eks-subnet.*.id
  }

  depends_on = [
    aws_iam_role_policy_attachment.data_manager_eks_cluster_iam_AmazonEKSClusterPolicy,
    aws_iam_role_policy_attachment.data_manager_eks_cluster_iam_AmazonEKSServicePolicy,
  ]
}

data "tls_certificate" "cluster_tls_certificate" {
  url = aws_eks_cluster.data_manager_eks_cluster.identity[0].oidc[0].issuer
}

resource "aws_iam_openid_connect_provider" "cluster_aws_iam_openid_connect_provider" {
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = [data.tls_certificate.cluster_tls_certificate.certificates.[*].sha1_fingerprint]
  url             = data.tls_certificate.cluster_tls_certificate.url

}

data "aws_iam_policy_document" "oidc_role_policy" {
  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]
    effect  = "Allow"

    condition {
      test     = "StringEquals"
      variable = "${replace(aws_iam_openid_connect_provider.cluster_aws_iam_openid_connect_provider.url, "https://", "")}:sub"
      values   = ["system:serviceaccount:kube-system:aws-node"]
    }

    principals {
      identifiers = ["${aws_iam_openid_connect_provider.cluster_aws_iam_openid_connect_provider.arn}"]
      type        = "Federated"
    }
  }
}

resource "aws_iam_role" "aws_node" {
  name = "data-manager-aws-node"
  assume_role_policy = data.aws_iam_policy_document.oidc_role_policy.json
  tags = {
    "ServiceAccountName"      = "aws-node"
    "ServiceAccountNameSpace" = "kube-system"
  }
  depends_on = [aws_iam_openid_connect_provider.cluster_aws_iam_openid_connect_provider]
}

resource "aws_iam_role_policy_attachment" "aws_node" {
  role       = aws_iam_role.aws_node.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  depends_on = [aws_iam_role.aws_node]
}

locals {
  kubeconfig = <<-KUBECONFIG
    ---
    apiVersion: v1
    clusters:
    - cluster:
        server: ${aws_eks_cluster.data_manager_eks_cluster.endpoint}
        certificate-authority-data: ${aws_eks_cluster.data_manager_eks_cluster.certificate_authority.0.data}
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

