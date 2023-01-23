module "lambda_layer" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 4.0"

  create_function = false
  create_layer    = true

  layer_name          = "data-manager-layer"
  compatible_runtimes = ["python3.9"]

  runtime = "python3.9"

  source_path = [
    {
      path             = "${path.module}/"
      pip_tmp_dir      = "${path.cwd}/"
      pip_requirements = "${path.module}/requirements.txt"
      prefix_in_zip    = "python" # required to get the path correct
    }
  ]
}

output "layer_arn" {
  value = module.lambda_layer.lambda_layer_arn
}