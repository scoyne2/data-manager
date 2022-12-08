locals {
  envs = { for tuple in regexall("(.*)=(.*)", file("../../.env")) : tuple[0] => sensitive(tuple[1]) }
}

provider "aws" {
  shared_config_files      = ["~/.aws/config"]
  shared_credentials_files = ["~/.aws/credentials"]
  profile                  = local.envs["AWS_PROFILE"]
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