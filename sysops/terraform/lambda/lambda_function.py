import os
import urllib.parse
from datetime import datetime

import boto3
import requests

client = boto3.client("emr-serverless")

EMR_SERVERLESS_APPLICATION_ID = os.environ.get("APPLICATION_ID", "999999999999")
JOB_ROLE_ARN = os.environ.get("JOB_ROLE_ARN", "arn:aws:iam::999999999999:role/fake")
SCRIPT_LOCATION = os.environ.get("SCRIPT_LOCATION", "test_script_location")
OUTPUT_BUCKET = os.environ.get("OUTPUT_BUCKET", "test_outputbucket")
RESOURCE_BUCKET = os.environ.get("RESOURCE_BUCKET", "test_resource_bucket")
DOMAIN_NAME = os.environ.get("DOMAIN_NAME", "localhost")

GRAPHQL_URL = f"https://api.{DOMAIN_NAME}/graphql"


def add_feed(vendor: str, feed_name: str, feed_method: str):
    vendor_clean = vendor.replace("_", " ").title()
    feed_name_clean = feed_name.replace("_", " ").title()
    query = (
        "mutation AddFeed {"
        f'  addFeed(vendor: "{vendor_clean}", feedName: "{feed_name_clean}", feedMethod: "{feed_method}")'
        "}"
    )
    r = requests.post(GRAPHQL_URL, json={"query": query})
    return r.status_code, r.json()


def file_received(vendor: str, feed_name: str, file_name: str):
    process_date = datetime.today().strftime("%Y-%m-%d %H:%M:%S")
    vendor_clean = vendor.replace("_", " ").title()
    feed_name_clean = feed_name.replace("_", " ").title()
    query = (
        "mutation UpdateFeedStatus {"
        "  updateFeedStatus("
        "    recordCount: 0"
        "    errorCount: 0"
        '    status: "Received"'
        f'    fileName: "{file_name}"'
        f'    vendor: "{vendor_clean}"'
        f'    feedName: "{feed_name_clean}"'
        f'    processDate: "{process_date}"'
        f'    EMRApplicationID: ""'
        f'    EMRStepID: ""'
        f"  )"
        "}"
    )
    r = requests.post(GRAPHQL_URL, json={"query": query})
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
        file_id,
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
    response_status_code, response_body = file_received(
        vendor.title(), feed.title(), file_id
    )

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
    return {"statusCode": 200, "body": f"Processing file: {file_id}"}
