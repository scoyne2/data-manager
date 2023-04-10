The below details how data-manager ingests a file.

# Trigger Lambda
The data-manager Lambda is triggered when a file is dropepd in the associated trigger S3 bucket. The lambda is built
using the python file `sysops/terraform/lambda/lambda_function.py` and tested via `sysops/terraform/lambda/test_lambda_function.py`.
Terraform will build `sysops/terraform/lambda/main.zip` which is deployed to AWS. Python libraries are made available to the Lambda via
LambdaLayer which is built by terraform from `sysops/terraform/lambda_layer/requirements.txt`. If you want the lambda to have additional libraries
you must add them to the `requirements.txt` file.

The lambda is triggered by the S3 trigger which is built in terraform here https://github.com/scoyne2/data-manager/blob/f23515314ff8fd66ceded930d7b1e9f7c4896534/sysops/terraform/lambda/lambda.tf#L143 when a file is dropped in the trigger bucket, it will
spin up the Lambda and the Lambda will start a spark job on EMRServerless which sill run pyspark file `pyspark/file_ingest.py`

# EMR Serverless
The spark jobs for file ingestion are executed using EMRServerless, the serverless cluster is defined using terraform at `sysops/terraform/emr/emr-serverless.tf`. Terraform will build and upload the spark requirements using the S3 resources bucket https://github.com/scoyne2/data-manager/blob/f1736065ff758d7f53ffae1bb2a299b009ba50d3/sysops/terraform/lambda/lambda.tf#L106

If any additional python requirements are needed by `pyspark/file_ingest.py` you should update `pyspark/requirements/requirements.txt`

The Spark job that the Lambda will run is `pyspark/file_ingest.py` which is tested by `pyspark/test_file_ingest.py`
The spark job will:
* create a spark session
* read in the csv file
* convert the data to parquet and write to the s3 output path
* perform quality checks
* notify downstream steps so that the data processing results can be seen via the data-manager UI
