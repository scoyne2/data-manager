import boto3
import os
import json
import urllib.parse

client = boto3.client('emr-serverless')

EMR_SERVERLESS_APPLICATION_ID = os.environ['APPLICATION_ID']
JOB_ROLE_ARN = os.environ['JOB_ROLE_ARN']
SCRIPT_LOCATION = os.environ['SCRIPT_LOCATION']
OUTPUT_BUCKET = os.environ['OUTPUT_BUCKET']
RESOURCE_BUCKET = os.environ['RESOURCE_BUCKET']

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
    
    python_zip_path = f"s3://{RESOURCE_BUCKET}/python/pyspark_ge.tar.gz#environment"
    spark_submit_args = (
       f"--conf spark.archives={python_zip_path}"
        "--conf spark.emr-serverless.driverEnv.PYSPARK_DRIVER_PYTHON=./environment/bin/python "
        "--conf spark.emr-serverless.driverEnv.PYSPARK_PYTHON=./environment/bin/python "
        "--conf spark.emr-serverless.executorEnv.PYSPARK_PYTHON=./environment/bin/python " 
        "--conf spark.hadoop.hive.metastore.client.factory.class=com.amazonaws.glue.catalog.metastore.AWSGlueDataCatalogHiveClientFactory "
        "--conf spark.executor.cores=1 "
        "--conf spark.executor.memory=4g "
        "--conf spark.driver.cores=1 "
        "--conf spark.driver.memory=4g "
        "--conf spark.executor.instances=1 "
    )

    client.start_job_run(
        applicationId=EMR_SERVERLESS_APPLICATION_ID,
        executionRoleArn=JOB_ROLE_ARN,
        jobDriver={
            "sparkSubmit": {
                "entryPoint": f's3://{SCRIPT_LOCATION}',
                "entryPointArguments": spark_args,
                "sparkSubmitParameters": spark_submit_args,
            }
        },
        configurationOverrides={
            "monitoringConfiguration": {
                "s3MonitoringConfiguration": {
                    "logUri": f"s3://{RESOURCE_BUCKET}/logs/{vendor}/{file_id}"
                }
            }
        },
    )
    return {"statusCode": 200, "body": f"Processing file: {input_file}"}
