import boto3
import os
import requests
import urllib.parse
from datetime import datetime

client = boto3.client("emr-serverless")

EMR_SERVERLESS_APPLICATION_ID = os.environ["APPLICATION_ID"]
JOB_ROLE_ARN = os.environ["JOB_ROLE_ARN"]
SCRIPT_LOCATION = os.environ["SCRIPT_LOCATION"]
OUTPUT_BUCKET = os.environ["OUTPUT_BUCKET"]
RESOURCE_BUCKET = os.environ["RESOURCE_BUCKET"]
DOMAIN_NAME = os.environ["DOMAIN_NAME"]

GRAPHQL_URL = f"https://api.{DOMAIN_NAME}/graphql"

def add_feed(vendor: str, feed_name: str, feed_method: str):
    query = """
        mutation AddFeed(
            $vendor: String!
            $feedName: String!
            $feedMethod: String!){
                addFeed(vendor: $vendor, feedName: $feedName, feedMethod: $feedMethod)
            }
    """
    r = requests.post(GRAPHQL_URL, json={'query': query, 'vendor': vendor,
                                  'feedName': feed_name, 'feedMethod': feed_method})
    return r.status_code, r.json()

def file_received(vendor: str, feed_name: str, file_name: str, feed_method: str):
    record_count = 0
    process_date = datetime.today().strftime("%Y-%m-%d")
    error_count =0
    status = "Received"
    query = """
        mutation UpdateFeedStatus(
            $vendor: String!
            $feedName: String!
            $fileName: String!
            $recordCount: String!
            $processDate: String!
            $errorCount: String!
            $status: String!
            ){
                updateFeedStatus(vendor: $vendor, feedName: $feedName, fileName: $fileName, $status: status, recordCount: $recordCount, processDate: $processDate, errorCount: $errorCount)
            }
    """
    r = requests.post(GRAPHQL_URL, json={'query': query, 'vendor': vendor,
                                  'feedName': feed_name, 'fileName': file_name,
                                  'recordCount': record_count, 'processDate': process_date,
                                  'errorCount': error_count, 'status': status})
    return r.status_code, r.json()

def lambda_handler(event, context):
    input_bucket = event["Records"][0]["s3"]["bucket"]["name"]
    key = urllib.parse.unquote_plus(
        event["Records"][0]["s3"]["object"]["key"], encoding="utf-8"
    )
    key_parts = key.split("/")
    vendor = key_parts[1]
    feed = key_parts[2]
    file_id = key_parts[3]
    file_parts = file_id.split(".")
    file_name = file_parts[0]
    file_extension = file_parts[1]

    feed_method = "S3"

    input_file = f"s3://{input_bucket}/{key}"
    output_path = f"s3://{OUTPUT_BUCKET}/{vendor}/{feed}/"
    spark_args = [
        "--vendor",
        vendor,
        "--feed_name",
        feed,
        "--file_name",
        file_name,
        "--feed_method",
        feed_method,
        "--graphql_url",
        GRAPHQL_URL,
        "--input_file",
        input_file,
        "--file_extension",
        file_extension,
        "--output_path",
        output_path,
        "--resources_bucket",
        RESOURCE_BUCKET,
    ]

    # Make Call to GraphQL API to add feed
    add_feed(vendor.title(), feed.title(), feed_method)
    response_status_code, response_body = file_received(vendor.title(), feed.title(), file_name, feed_method)
    return {"statusCode": response_status_code, "body": f" resp: {response_body}"}

    python_zip_path = f"s3://{RESOURCE_BUCKET}/pyspark_requirements/pyspark_requirements.tar.gz#environment"
    spark_submit_args = (
        f"--conf spark.archives={python_zip_path} "
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
                "entryPoint": f"s3://{SCRIPT_LOCATION}",
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
