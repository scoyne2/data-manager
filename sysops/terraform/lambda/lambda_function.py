import boto3
import os
import json
import urllib.parse

client = boto3.client('emr-serverless')

EMR_SERVERLESS_APPLICATION_ID = os.environ['EMR_SERVERLESS_APPLICATION_ID']
JOB_ROLE_ARN = os.environ['JOB_ROLE_ARN']
SPARK_SUBMIT_ARGS = os.environ['SPARK_SUBMIT_ARGS']
SCRIPT_LOCATION = os.environ['SCRIPT_LOCATION']
OUTPUT_BUCKET = os.environ['OUTPUT_BUCKET']
LOG_BUCKET = os.environ['LOG_BUCKET']

def lambda_handler(event, context):  
    input_bucket = event['Records'][0]['s3']['bucket']['name']
    key = urllib.parse.unquote_plus(event['Records'][0]['s3']['object']['key'], encoding='utf-8')
    key_parts = key.split('/')
    inbound = key_parts[0]
    vendor = key_parts[1]
    feed = key_parts[2]
    file_id = key_parts[3]
    file_parts = file_id.split(".")
    file_name = file_parts[0]
    file_extension = file_parts[1]
    
    input_file = f"s3://{input_bucket}/{key}"
    output_path = f"s3://{OUTPUT_BUCKET}/{vendor}/{feed}/"
    spark_args = ["--input_file", input_file, "--file_extension", file_extension, "--output_path", output_path]

    client.start_job_run(
        applicationId=EMR_SERVERLESS_APPLICATION_ID,
        executionRoleArn=JOB_ROLE_ARN,
        jobDriver={
            "sparkSubmit": {
                "entryPoint": SCRIPT_LOCATION,
                "entryPointArguments": spark_args,
                "sparkSubmitParameters": SPARK_SUBMIT_ARGS,
            }
        },
        configurationOverrides={
            "monitoringConfiguration": {
                "s3MonitoringConfiguration": {
                    "logUri": f"s3://{LOG_BUCKET}/logs/{vendor}/{file_id}"
                }
            }
        },
    )