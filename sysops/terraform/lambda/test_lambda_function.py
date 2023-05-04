import unittest
from unittest.mock import MagicMock, patch

from lambda_function import lambda_handler


class TestLambdaHandler(unittest.TestCase):
    def setUp(self):
        self.event = {
            "Records": [
                {
                    "s3": {
                        "bucket": {"name": "input-bucket"},
                        "object": {"key": "inbound/vendor/feed/file.txt"},
                    }
                }
            ]
        }
        self.context = MagicMock()

    @patch("botocore.client.BaseClient._make_api_call")
    def test_lambda_handler(self, mock_make_api_call):
        lambda_handler(self.event, self.context)

        mock_make_api_call.assert_called_with(
            "StartJobRun",
            {
                "applicationId": "999999999999",
                "executionRoleArn": "arn:aws:iam::999999999999:role/fake",
                "jobDriver": {
                    "sparkSubmit": {
                        "entryPoint": "s3://test_script_location",
                        "entryPointArguments": [
                            "--vendor",
                            "vendor",
                            "--feed_name",
                            "feed",
                            "--file_name",
                            "file.txt",
                            "--feed_method",
                            "S3",
                            "--graphql_url",
                            "https://api.datamanagertool.com/graphql",
                            "--input_file",
                            "s3://input-bucket/inbound/vendor/feed/file.txt",
                            "--file_extension",
                            "txt",
                            "--output_path",
                            "s3://test_outputbucket/vendor/feed/",
                            "--resources_bucket",
                            "test_resource_bucket",
                        ],
                        "sparkSubmitParameters": "--conf spark.archives=s3://test_resource_bucket/pyspark_requirements/pyspark_requirements.tar.gz#environment --conf spark.emr-serverless.driverEnv.PYSPARK_DRIVER_PYTHON=./environment/bin/python --conf spark.emr-serverless.driverEnv.PYSPARK_PYTHON=./environment/bin/python --conf spark.emr-serverless.executorEnv.PYSPARK_PYTHON=./environment/bin/python --conf spark.hadoop.hive.metastore.client.factory.class=com.amazonaws.glue.catalog.metastore.AWSGlueDataCatalogHiveClientFactory --conf spark.executor.cores=1 --conf spark.executor.memory=4g --conf spark.driver.cores=1 --conf spark.driver.memory=4g --conf spark.executor.instances=1 ",
                    }
                },
                "configurationOverrides": {
                    "monitoringConfiguration": {
                        "s3MonitoringConfiguration": {
                            "logUri": "s3://test_resource_bucket/logs/vendor/file.txt"
                        }
                    }
                },
            },
        )


if __name__ == "__main__":
    unittest.main()
