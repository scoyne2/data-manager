# Trigger Lambda
The data-manager Lambda is triggered when a file is dropepd in the associated trigger S3 bucket. The lambda is built
using the python file `sysops/terraform/lambda/lambda_function.py` and tested via `sysops/terraform/lambda/test_lambda_function.py`.
Terraform will build `sysops/terraform/lambda/main.zip` which is deployed to AWS. Python libraries are made available to teh Lambda via
LambdaLayer which is built by terraform from `sysops/terraform/lambda_layer/requirements.txt`. If you want the lambda to have additional libraries
you must add them to the `requirements.txt` file.

The lambda is triggered by the S3 trigger which is built in terraform here https://github.com/scoyne2/data-manager/blob/f23515314ff8fd66ceded930d7b1e9f7c4896534/sysops/terraform/lambda/lambda.tf#L143

