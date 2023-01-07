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

output "aws_emrserverless_application_id" {
  value = aws_emrserverless_application.data_manager_emr_serverless.id
}