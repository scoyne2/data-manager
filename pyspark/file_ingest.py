import argparse
import logging
from datetime import datetime
import os
import tempfile

import boto3
import requests
import yaml

from pyspark.sql.functions import lit
from pyspark.sql.session import SparkSession
from great_expectations.core.batch import RuntimeBatchRequest

import great_expectations as ge
from great_expectations.core.expectation_suite import ExpectationConfiguration

EMR_CLUSTER_ID = os.environ.get('SERVERLESS_EMR_VIRTUAL_CLUSTER_ID', default='unknown')
EMR_STEP_ID = os.environ.get('SERVERLESS_EMR_JOB_ID', default='unknown')

def create_spark_session(input_file):
    spark = (
        SparkSession.builder.appName(input_file)
        .enableHiveSupport()
        .config("mapreduce.fileoutputcommitter.algorithm.version", "2")
        .getOrCreate()
    )
    return spark


def inspect_file(header):
    # Check if header has '|' or ',' if so set delimiter
    quote = ""
    delimiter = ""
    if header.count('"') >= 2:
        quote = '"'
    if header.count("|") >= 1:
        delimiter = "|"
    elif header.count(",") >= 1:
        delimiter = ","
    return delimiter, quote


def process_file(spark, input_df, output_path, vendor, feed_name, file_name) -> int:
    # read the input file, header is required
    input_count = input_df.count()

    # write to parquet, allow overwrite. partitioned by /dt=YYYY-MM-DD/
    today = datetime.today().strftime("%Y-%m-%d")
    output_df = input_df.withColumn("dt", lit(today)).withColumn("file_name", lit(file_name))

    database_name = f"data_manager_output_{vendor}"
    glue_table_name = f"{database_name}.{feed_name}"
    spark.sql(f"CREATE DATABASE IF NOT EXISTS {database_name}")
    output_df.write.partitionBy("dt").format("parquet").mode("overwrite").saveAsTable(
        glue_table_name, path=output_path
    )

    # Log output file path
    logging.info(f"Wrote data to {output_path}")

    # QA counts
    output_count = spark.read.parquet(output_path).count()
    logging.info(f"Read {input_count} rows, wrote {output_count} rows.")
    return output_count

def get_ge_context():
    with tempfile.TemporaryDirectory() as temp_dir:
        # download data_context_config from S3
        destination = os.path.join(temp_dir, "great_expectations.yml")
        s3 = boto3.client("s3")
        key = "great_expectations/great_expectations.yml"
        s3.download_file(resources_bucket, key, destination)
        context = ge.DataContext(temp_dir)
        return context


def perform_quality_checks(pandas_df, qa_suite) -> int:
    context_gx = get_ge_context()
    suite = context_gx.create_expectation_suite(qa_suite, overwrite_existing=True)
    exp_config = ExpectationConfiguration(
        expectation_type="expect_column_values_to_not_be_null",
            kwargs={
                "auto": True,
                "column": "file_name",
            },
    )
    suite.add_expectation(exp_config)

    # Create Runtime Batch
    pandasDataset = ge.dataset.PandasDataset(pandas_df, expectation_suite=suite)
    batch_suite = pandasDataset.get_expectation_suite()

    # Create expectations
    context_gx.save_expectation_suite(batch_suite)

    # Build Data Docs
    context_gx.build_data_docs()
    error_count = 0
    return error_count


def update_feed_status(
    graphql_url: str,
    vendor: str,
    feed_name: str,
    file_name: str,
    record_count: int,
    error_count: int,
    status: str,
    emr_application_id: str,
    emr_step_id: str
):
    process_date = datetime.today().strftime("%Y-%m-%d %H:%M:%S")
    vendor_clean = vendor.replace("_", " ").title()
    feed_name_clean = feed_name.replace("_", " ").title()
    query = (
        "mutation UpdateFeedStatus {"
        "  updateFeedStatus("
        f"  recordCount: {record_count}"
        f"  errorCount: {error_count}"
        f'  status: "{status}"'
        f'  fileName: "{file_name}"'
        f'  vendor: "{vendor_clean}"'
        f'  feedName: "{feed_name_clean}"'
        f'  processDate: "{process_date}"'
        f'  emrApplicationID: "{emr_application_id}"'
        f'  emrStepID: "{emr_step_id}"'
        f"  )"
        "}"
    )
    r = requests.post(graphql_url, json={"query": query})
    return r.status_code

if __name__ == "__main__":
    # Read args
    parser = argparse.ArgumentParser(description="Spark Job Arguments")
    parser.add_argument("--vendor", required=True, help="Vendor Name")
    parser.add_argument("--feed_name", required=True, help="Feed Name")
    parser.add_argument("--file_name", required=True, help="File Name")
    parser.add_argument("--feed_method", required=True, help="Feed Method")
    parser.add_argument("--graphql_url", required=True, help="GraphQL URL")
    parser.add_argument("--input_file", required=True, help="Input S3 Path")
    parser.add_argument("--output_path", required=True, help="Output S3 Path")
    parser.add_argument("--file_extension", required=True, help="File Extension")
    parser.add_argument("--resources_bucket", required=True, help="Resources S3 Bucket")
    args = parser.parse_args()

    # Parse args
    vendor = args.vendor
    feed_name = args.feed_name
    file_name = args.file_name
    feed_method = args.feed_method
    graphql_url = args.graphql_url
    input_file = args.input_file
    output_path = args.output_path
    file_extension = args.file_extension
    resources_bucket = args.resources_bucket
    logging.info(f"Processing file: {input_file}")
    qa_suite = f"{vendor}_{feed_name}_{file_name.split('.')[0]}"

    # Create spark session
    spark = create_spark_session(input_file)

    # Gather data about file format
    header = spark.read.text(input_file).first()[0]
    delimiter, quote = inspect_file(header)

    # Read file
    input_df = spark.read.csv(input_file, sep=delimiter, quote=quote, header=True)

    # Update Feed Status
    record_count = input_df.count()
    error_count = 0
    response_code = update_feed_status(
        graphql_url,
        vendor,
        feed_name,
        file_name,
        record_count,
        error_count,
        "Processing",
        EMR_CLUSTER_ID,
        EMR_STEP_ID
    )
    logging.info(f"Update feed status response code {response_code}")

    # Convert file to parquet
    processed_count = process_file(spark, input_df, output_path, vendor, feed_name, file_name)
    error_count = record_count - processed_count
    response_code = update_feed_status(
        graphql_url,
        vendor,
        feed_name,
        file_name,
        record_count,
        error_count,
        "Validating",
        EMR_CLUSTER_ID,
        EMR_STEP_ID
    )
    logging.info(f"Convert status response code {response_code}")

    # Run Quality Checks and log results
    output_df = spark.read.parquet(output_path)
    pandas_df = output_df.toPandas()
    error_count_from_qa = perform_quality_checks(pandas_df, qa_suite)
    error_count = error_count_from_qa

    # Determine status
    # TODO allow user to configure threshold peer feed on what is considered a failure
    if record_count == error_count:
        status = "Failed"
    elif error_count_from_qa > 0:
        status = "Errors"
    else:
        status = "Success"

    # Update Final Feed Status
    response_code = update_feed_status(
        graphql_url,
        vendor,
        feed_name,
        file_name,
        record_count,
        error_count_from_qa,
        status,
        EMR_CLUSTER_ID,
        EMR_STEP_ID
    )
    logging.info(f"Final status response code {response_code}")
