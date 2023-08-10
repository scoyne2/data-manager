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
  profile                  = local.envs[" "]
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
  handler          = "datadog_lambda.handler.handler"
  timeout          = 60
  layers           = [
    var.layer_arn,
    "arn:aws:lambda:us-west-2:464622532012:layer:Datadog-python39:78",
    "arn:aws:lambda:us-west-2:464622532012:layer:Datadog-Extension:45"
    ]
  source_code_hash = filebase64sha256("${path.module}/main.zip")
  environment {
    variables = {
      APPLICATION_ID        = var.aws_emrserverless_application_id
      OUTPUT_BUCKET         = aws_s3_bucket.data_manager_processed_s3.bucket
      RESOURCE_BUCKET       = aws_s3_bucket.data_manager_resources_s3.bucket
      SCRIPT_LOCATION       = "${aws_s3_bucket.data_manager_resources_s3.bucket}/scripts/file_ingest.py"
      JOB_ROLE_ARN          = var.aws_emrserverless_role_arn
      DOMAIN_NAME           = local.envs["DOMAIN_NAME"]
      DD_SITE               = "us5.datadoghq.com"
      DD_LAMBDA_HANDLER     = "lambda_function.lambda_handler"
      DD_API_KEY            = local.envs["DD_API_KEY"]
      DD_TRACE_ENABLED      = "true"
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
  key = "scripts/file_ingest.py"
  source = "${path.module}/../../../pyspark/file_ingest.py"
  etag = filemd5("${path.module}/../../../pyspark/file_ingest.py")
}

resource "aws_s3_bucket_object" "great_expectations" {
  bucket = aws_s3_bucket.data_manager_resources_s3.id
  key = "great_expectations/great_expectations.yml"
  source = "${path.module}/../../../great_expectations/great_expectations.yml"
  etag = filemd5("${path.module}/../../../great_expectations/great_expectations.yml")
}

resource "aws_s3_bucket_object" "python_files" {
  bucket = aws_s3_bucket.data_manager_resources_s3.id
  key = "pyspark_requirements/pyspark_requirements.tar.gz"
  source = "${path.module}/../../../pyspark/requirements/out/pyspark_requirements.tar.gz"
  etag = filemd5("${path.module}/../../../pyspark/requirements/out/pyspark_requirements.tar.gz")
}

resource "aws_s3_bucket_object" "test_file" {
  bucket = aws_s3_bucket.data_manager_trigger_s3.id
  key = "inbound/coyne_enterprises/hello/hello.csv"
  source = "${path.module}/../../../tests/lambda_manual_test/hello.csv"
  etag = filemd5("${path.module}/../../../tests/lambda_manual_test/hello.csv")
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
