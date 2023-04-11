import argparse
from datetime import datetime
import logging

import boto3
import yaml

from pyspark.sql.session import SparkSession


def create_spark_session(input_file):
    spark = (
        SparkSession.builder.appName(input_file)
        .enableHiveSupport()
        .config("mapreduce.fileoutputcommitter.algorithm.version", "2")
        .getOrCreate()
    )
    return spark


def inspect_file(spark, header):
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


def process_file(spark, delimiter, quote, input_file, output_path):
    # read the input file, header is required
    df = spark.read.csv(input_file, sep=delimiter, quote=quote, header=True)
    input_count = df.count()
    logging.info(f"Reading file {input_file} with {input_count} rows")

    # write to parquet, allow overwrite. partitioned by /FILENAME/dt=YYYY-MM-DD/
    today = datetime.today().strftime("%Y-%m-%d")
    output_path = f"{output_path}/dt={today}/"
    df.write.parquet(path=output_path, mode="overwrite")

    # Log output file path
    logging.info(f"Wrote data to {output_path}")

    # QA counts
    output_count = spark.read.parquet(output_path).count()
    logging.info(f"Read {input_count} rows, wrote {output_count} rows.")


def perform_quality_checks(output_path, resources_bucket):
    session = boto3.Session()
    s3_client = session.client("s3")
    response = s3_client.get_object(
        Bucket=resources_bucket,
        Key="great_expectations/great_expectations.yml",
    )
    config_file = yaml.safe_load(response["Body"])


#    df = spark.read.parquet(output_path)

#    config = DataContextConfig(
#         config_version=config_file["config_version"],
#         datasources=config_file["datasources"],
#         expectations_store_name=config_file["expectations_store_name"],
#         validations_store_name=config_file["validations_store_name"],
#         evaluation_parameter_store_name=config_file["evaluation_parameter_store_name"],
#         plugins_directory="/great_expectations/plugins",
#         stores=config_file["stores"],
#         data_docs_sites=config_file["data_docs_sites"],
#         config_variables_file_path=config_file["config_variables_file_path"],
#         anonymous_usage_statistics=config_file["anonymous_usage_statistics"],
#         checkpoint_store_name=config_file["checkpoint_store_name"],
#         store_backend_defaults=S3StoreBackendDefaults(
#             default_bucket_name=config_file["data_docs_sites"]["s3_site"][
#                 "store_backend"
#             ]["bucket"]
#         ),
#     )

#    context_gx = get_context(project_config=config)

#    expectation_suite_name = suite_name
#    suite = context_gx.get_expectation_suite(suite_name)

#    batch_request = RuntimeBatchRequest(
#      datasource_name="spark_s3",
#      data_connector_name="default_inferred_data_connector_name",
#      data_asset_name="datafile_name",
#      batch_identifiers={"runtime_batch_identifier_name": "default_identifier"},
#      runtime_parameters={"path": output_path},
#    )
#    validator = context_gx.get_validator(
#      batch_request=batch_request,
#      expectation_suite_name=expectation_suite_name,
#    )
#    print(validator.head())
#    # TODO Add Tests
#    # validator.expect_column_values_to_not_be_null(
#    #  column="passenger_count"
#    #)
#    validator.save_expectation_suite(discard_failed_expectations=False)
#    my_checkpoint_name = "in_memory_checkpoint"
#    python_config = {
#      "name": my_checkpoint_name,
#      "class_name": "Checkpoint",
#      "config_version": 1,
#      "run_name_template": "%Y%m%d-%H%M%S-my-run-name-template",
#      "action_list": [
#          {
#              "name": "store_validation_result",
#              "action": {"class_name": "StoreValidationResultAction"},
#          },
#          {
#              "name": "store_evaluation_params",
#              "action": {"class_name": "StoreEvaluationParametersAction"},
#          },
#      ],
#      "validations": [
#          {
#              "batch_request": {
#                  "datasource_name": "spark_s3",
#                  "data_connector_name": "default_runtime_data_connector_name",
#                  "data_asset_name": "pyspark_df",
#              },
#              "expectation_suite_name": expectation_suite_name,
#          }
#      ],
#    }
#    context_gx.add_checkpoint(**python_config)

#    results = context_gx.run_checkpoint(
#      checkpoint_name=my_checkpoint_name,
#      run_name="run_name",
#      batch_request={
#          "runtime_parameters": {"batch_data": df},
#          "batch_identifiers": {
#              "runtime_batch_identifier_name": "default_identifier"
#          },
#      },
#    )

#    validation_result_identifier = results.list_validation_result_identifiers()[0]
#    context_gx.build_data_docs()


def notify_downstream():
    pass


if __name__ == "__main__":
    # Read args
    parser = argparse.ArgumentParser(description="Spark Job Arguments")
    parser.add_argument("--input_file", required=True, help="Input S3 Path")
    parser.add_argument("--output_path", required=True, help="Output S3 Path")
    parser.add_argument("--file_extension", required=True, help="File Extension")
    parser.add_argument("--resources_bucket", required=True, help="Resources S3 Bucket")
    args = parser.parse_args()

    # Parse args
    input_file = args.input_file
    output_path = args.output_path
    file_extension = args.file_extension
    resources_bucket = args.resources_bucket
    logging.info(f"Processing file: {input_file}")

    # Create spark session
    spark = create_spark_session(input_file)

    # Gather data about file format
    header = spark.read.text(input_file).first()[0]
    delimiter, quote = inspect_file(spark, header)

    # Convert file to parquet
    process_file(spark, delimiter, quote, input_file, output_path)

    # Run Quality Checks and log results
    perform_quality_checks(output_path, resources_bucket)

    # Make call to data-manager API to move on to next steps (update postgres table, trigger glue crawlers, etc)
    notify_downstream()
