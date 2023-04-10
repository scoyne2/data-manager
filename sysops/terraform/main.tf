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

module "lambda_layer" {
  source = "./lambda_layer"
}

module "lambda" {
  source = "./lambda"
  aws_emrserverless_application_id = module.emr.aws_emrserverless_application_id
  layer_arn = module.lambda_layer.layer_arn
  aws_emrserverless_role_arn = module.emr.aws_emrserverless_role_arn
  aws_emrserverless_policy_arn = module.emr.aws_emrserverless_policy_arn
}

module "emr" {
  source = "./emr"
}

output "aws_acm_certificate_arn" {
  description = "ARN of ACM Certificate"
  value = module.eks.aws_acm_certificate_arn
  sensitive = false
}

output "aws_wafv2_web_acl_arn" {
  description = "ARN of waf"
  value = module.eks.aws_wafv2_web_acl_arn
  sensitive = false
}