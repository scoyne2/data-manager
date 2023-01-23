locals {
  envs = { for tuple in regexall("(.*)=(.*)", file("../../.env")) : tuple[0] => sensitive(tuple[1]) }
}

provider "aws" {
  shared_config_files      = ["~/.aws/config"]
  shared_credentials_files = ["~/.aws/credentials"]
  profile                  = local.envs["AWS_PROFILE"]
}

resource "aws_emrserverless_application" "data_manager_emr_serverless" {
  name          = "data_manager_emr_serverless"
  release_label = "emr-6.6.0"
  type          = "spark"

  initial_capacity {
    initial_capacity_type = "Driver"

    initial_capacity_config {
      worker_count = 1
      worker_configuration {
        cpu    = "2 vCPU"
        memory = "10 GB"
      }
    }
  }

  initial_capacity {
    initial_capacity_type = "Executor"

    initial_capacity_config {
      worker_count = 1
      worker_configuration {
        cpu    = "2 vCPU"
        memory = "10 GB"
      }
    }
  }

  maximum_capacity {
    cpu    = "15 vCPU"
    memory = "100 GB"
  }
}

resource "aws_iam_role" "iam_for_emr" {
  name = "iam_for_emr"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "emr-serverless.amazonaws.com"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "aws_emr_serverless_policy_attachment" {
  role       = aws_iam_role.iam_for_emr.name
  policy_arn = aws_iam_policy.emr_serverless_iam_policy.arn
  depends_on = [aws_iam_role.iam_for_emr, aws_iam_policy.emr_serverless_iam_policy]
}

output "aws_emrserverless_application_id" {
  value = aws_emrserverless_application.data_manager_emr_serverless.id
}

output "aws_emrserverless_role_arn" {
  value = aws_iam_role.iam_for_emr.arn
}

output "aws_emrserverless_policy_arn" {
  value = aws_iam_policy.emr_serverless_iam_policy.arn
}