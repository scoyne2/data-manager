# Data Ingestion Process

Below details how Data Manager ingests a file end-to-end.

## Trigger Lambda

When a file is dropped in the S3 trigger bucket, an AWS Lambda (`sysops/terraform/lambda/lambda_function.py`) is invoked. Terraform packages this handler into `sysops/terraform/lambda/main.zip` and deploys it. Python libraries are provided via a Lambda Layer built from `sysops/terraform/lambda_layer/requirements.txt`. To add extra libraries for the Lambda, update that `requirements.txt` and re-deploy.

The Lambda reads the S3 event, parses the S3 key to extract `vendor`, `feed`, and `file_name`, then:
1. Calls the GraphQL API to add or update feed metadata.
2. Starts an EMR Serverless Spark job to process the file.

## EMR Serverless

EMR Serverless jobs are defined in Terraform (`sysops/terraform/emr/emr-serverless.tf`) and run the PySpark script `pyspark/file_ingest.py`. Spark requirements are uploaded to S3 via Terraform as part of the resources bucket.

If additional Python dependencies are needed by the Spark job, update `pyspark/requirements/requirements.txt`.

The Spark job will:
- Create or connect to the Hive/Glue database `data_manager_output_{vendor}`
- Read and parse the flat file (automatic delimiter and quote detection)
- Write data to Parquet partitioned by date (`/dt=YYYY-MM-DD/`) and register as a Glue table
- Perform data quality checks using Great Expectations
- Update the feed status via GraphQL mutations at each stage

## Deployment

Use the provided Helm chart (`sysops/helm/data-manager`) to deploy the ingestion components:
```bash
bash scripts/deploy_helm.sh
```
This script applies the Kubernetes resources for the Lambda triggers (via respective controllers) and EMR Serverless configurations.