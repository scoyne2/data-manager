variable "aws_emrserverless_application_id" {}
variable "layer_arn" {}
variable "aws_emrserverless_policy_arn" {}
variable "aws_emrserverless_role_arn" {}

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

resource "random_id" "unique_suffix" {
  byte_length = 2
}

data "archive_file" "python_lambda_package" {  
  type = "zip"  
  source_file  = "${path.module}/lambda_function.py"
  output_path  = "${path.module}/main.zip"
}

resource "aws_lambda_function" "data_manager_filter_lambda_func" {
  function_name    = local.app_id
  role             = aws_iam_role.iam_for_lambda.arn
  filename         = "${path.module}/main.zip"
  runtime          = "python3.9"
  handler          = "lambda_function.lambda_handler"
  timeout          = 60
  layers           = [var.layer_arn]
  environment {
    variables = {
      APPLICATION_ID     = var.aws_emrserverless_application_id
      OUTPUT_BUCKET      = aws_s3_bucket.data_manager_processed_s3.bucket
      LOG_BUCKET         = aws_s3_bucket.data_manager_resources_s3.bucket
      SPARK_SUBMIT_ARGS  = "--conf spark.executor.cores=1 --conf spark.executor.memory=4g --conf spark.driver.cores=1 --conf spark.driver.memory=4g --conf spark.executor.instances=1"
      SCRIPT_LOCATION    = "${aws_s3_bucket.data_manager_resources_s3.bucket}/scripts/spark.py"
      JOB_ROLE_ARN       = var.aws_emrserverless_role_arn
    }
  }
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
  bucket = local.envs["TRIGGER_BUCKET_NAME"]
}

resource "aws_s3_bucket" "data_manager_processed_s3" {
  bucket = local.envs["PROCESSED_BUCKET_NAME"]
}

resource "aws_s3_bucket" "data_manager_resources_s3" {
  bucket = local.envs["RESOURCES_BUCKET_NAME"]
}

resource "aws_s3_bucket_object" "spark_script" {
  bucket = aws_s3_bucket.data_manager_resources_s3.id
  key = "scripts/spark.py"
  source = "${path.module}/spark.py"
}

resource "aws_s3_bucket_acl" "data_manager_trigger_acl" {
  bucket = aws_s3_bucket.data_manager_trigger_s3.id
  acl    = "private"
}

resource "aws_s3_bucket_acl" "data_manager_processed_acl" {
  bucket = aws_s3_bucket.data_manager_processed_s3.id
  acl    = "private"
}

resource "aws_s3_bucket_acl" "data_manager_resources_acl" {
  bucket = aws_s3_bucket.data_manager_resources_s3.id
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
    filter_prefix       = "inbound"
  }

  depends_on = [aws_lambda_permission.allow_bucket]
}

resource "aws_iam_role_policy_attachment" "basic" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.iam_for_lambda.name
}

resource "aws_iam_role_policy_attachment" "emr_serverless" {
  policy_arn = var.aws_emrserverless_policy_arn
  role       = aws_iam_role.iam_for_lambda.name
}