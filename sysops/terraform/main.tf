provider "aws" {
  shared_config_files      = ["~/.aws/config"]
  shared_credentials_files = ["~/.aws/credentials"]
  profile                  = "personal"
}

data "aws_availability_zones" "available" {
  state = "available"
}

module "eks" {
  source = "./eks"
}

module "lambda" {
  source = "./lambda"
}