import boto3
import os

client = boto3.client('emr-serverless')

EMR_SERVERLESS_APPLICATION_ID = os.environ['EMR_SERVERLESS_APPLICATION_ID']
JOB_ROLE_ARN = os.environ['JOB_ROLE_ARN']
SPARK_SUBMIT_ARGS = os.environ['SPARK_SUBMIT_ARGS']
SCRIPT_LOCATION = os.environ['SCRIPT_LOCATION']

def lambda_handler(event, context):  
    # TODO get from event 
    vendor = 'test'
    file_id = 'test'
    input_file = "todo"
    output_path = "todo"
    spark_args = ["-input_file", input_file, "-output_path", output_path]

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
                        "logUri": f"s3://{S3_BUCKET_NAME}/logs/{vendor}/{file_id}"
                    }
                }
            },
        )