variable "node-type" {
  default = "t2.medium"
  type    = string
}

resource "aws_iam_role" "data-manager-eks-nodegroup-iam" {
  name = "data-manager-eks-nodegroup"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "data-manager-eks-nodegroup-iam-AmazonEKSWorkerNodePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.data-manager-eks-nodegroup-iam.name
}

resource "aws_iam_role_policy_attachment" "data-manager-eks-nodegroup-iam-AmazonEKS_CNI_Policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.data-manager-eks-nodegroup-iam.name
}

resource "aws_iam_role_policy_attachment" "data-manager-eks-nodegroup-iam-AmazonEC2ContainerRegistryReadOnly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.data-manager-eks-nodegroup-iam.name
}

resource "aws_iam_role_policy_attachment" "data-manager-eks-nodegroup-iam-ElasticLoadBalancingFullAccess" {
  policy_arn = "arn:aws:iam::aws:policy/ElasticLoadBalancingFullAccess"
  role       = aws_iam_role.data-manager-eks-nodegroup-iam.name
}

resource "aws_security_group" "data-manager-eks-nodegroup" {
  name        = "data-manager-eks-node"
  description = "Security group for all nodes in the cluster"
  vpc_id      = aws_vpc.data-manager-vpc.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    "Name"                                      = "data-manager-eks-node"
    "kubernetes.io/cluster/${var.cluster-name}" = "owned"
  }
}

resource "aws_security_group_rule" "data-manager-eks-node-ingress-self" {
  description              = "Allow node to communicate with each other"
  from_port                = 0
  protocol                 = "-1"
  security_group_id        = aws_security_group.data-manager-eks-nodegroup.id
  source_security_group_id = aws_security_group.data_manager_eks_cluster_sg.id
  to_port                  = 65535
  type                     = "ingress"
}

resource "aws_security_group_rule" "data-manager-eks-node-ingress-cluster" {
  description              = "Allow worker Kubelets and pods to receive communication from the cluster control      plane"
  from_port                = 1025
  protocol                 = "tcp"
  security_group_id        = aws_security_group.data-manager-eks-nodegroup.id
  source_security_group_id = aws_security_group.data_manager_eks_cluster_sg.id
  to_port                  = 65535
  type                     = "ingress"
}

resource "aws_security_group_rule" "data_manager_eks_cluster_ingress_node_https" {
  description              = "Allow pods to communicate with the cluster API Server"
  from_port                = 443
  protocol                 = "tcp"
  security_group_id        = aws_security_group.data-manager-eks-nodegroup.id
  source_security_group_id = aws_security_group.data_manager_eks_cluster_sg.id
  to_port                  = 443
  type                     = "ingress"
}

resource "aws_eks_node_group" "data-manager-eks-nodegroup" {
  cluster_name    = aws_eks_cluster.data_manager_eks_cluster.name
  node_group_name = "data-manager-eks"
  node_role_arn   = aws_iam_role.data-manager-eks-nodegroup-iam.arn
  subnet_ids      = aws_subnet.data-manager-public-subnet[*].id
  instance_types  = [var.node-type]

  scaling_config {
    desired_size = 1
    max_size     = 4
    min_size     = 1
  }

  # Ensure that IAM Role permissions are created before and deleted after EKS Node Group handling.
  # Otherwise, EKS will not be able to properly delete EC2 Instances and Elastic Network Interfaces.
  depends_on = [
    aws_iam_role_policy_attachment.data-manager-eks-nodegroup-iam-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.data-manager-eks-nodegroup-iam-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.data-manager-eks-nodegroup-iam-AmazonEC2ContainerRegistryReadOnly,
  ]
}