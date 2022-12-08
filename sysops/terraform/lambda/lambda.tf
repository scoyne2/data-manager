locals {
  envs = { for tuple in regexall("(.*)=(.*)", file("../../.env")) : tuple[0] => sensitive(tuple[1]) }
}

provider "aws" {
  shared_config_files      = ["~/.aws/config"]
  shared_credentials_files = ["~/.aws/credentials"]
  profile                  = local.envs["AWS_PROFILE"]
}

variable "app_name" {
  description = "Application name"
  default     = "dm-file-notifications"
}

variable "app_env" {
  description = "Application environment tag"
  default     = "dev"
}

locals {
  app_id = "${lower(var.app_name)}-${lower(var.app_env)}-${random_id.unique_suffix.hex}"
}

data "archive_file" "zip" {
  type        = "zip"
  source_file = "build/bin/app"
  output_path = "build/bin/app.zip"
}

resource "random_id" "unique_suffix" {
  byte_length = 2
}

resource "aws_lambda_function" "data_manager_filter_lambda_func" {
  filename         = data.archive_file.zip.output_path
  function_name    = local.app_id
  handler          = "app"
  source_code_hash = data.archive_file.zip.output_base64sha256
  runtime          = "go1.x"
  role             = aws_iam_role.iam_for_lambda.arn
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_s3_bucket" "data_manager_trigger_s3" {
  bucket = "data-manager-trigger"
}

resource "aws_s3_bucket_acl" "data_manager_trigger_acl" {
  bucket = aws_s3_bucket.data_manager_trigger_s3.id
  acl    = "private"
}

resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.data_manager_filter_lambda_func.arn
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.data_manager_trigger_s3.arn
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = aws_s3_bucket.data_manager_trigger_s3.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.data_manager_filter_lambda_func.arn
    events              = ["s3:ObjectCreated:*"]
  }

  depends_on = [aws_lambda_permission.allow_bucket]
}