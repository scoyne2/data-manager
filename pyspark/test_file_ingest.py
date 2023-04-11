import unittest
import boto3
import yaml
from unittest.mock import patch

from pyspark.sql.session import SparkSession

from file_ingest import (
    create_spark_session,
    inspect_file,
    process_file,
    perform_quality_checks,
)


class TestMyScript(unittest.TestCase):
    def setUp(self):
        self.spark = SparkSession.builder.appName("TestInspectFile").getOrCreate()
        self.output_path = "path/to/output/folder"

    def tearDown(self):
        self.spark.stop()
        
    def test_create_spark_session(self):
        spark = create_spark_session("test")
        self.assertIsInstance(spark, SparkSession)
        
    def test_inspect_file(self):
        header1 = 'Name, Age, Gender'
        header2 = 'Name|Age|Gender'
        header3 = '"Name","Age","Gender"'
        delimiter1, quote1 = inspect_file(self.spark, header1)
        delimiter2, quote2 = inspect_file(self.spark, header2)
        delimiter3, quote3 = inspect_file(self.spark, header3)
        self.assertEqual(delimiter1, ",")
        self.assertEqual(delimiter2, "|")
        self.assertEqual(delimiter3, ",")
        self.assertEqual(quote1, "")
        self.assertEqual(quote2, "")
        self.assertEqual(quote3, '"')
        
    def test_process_file(self):
        input_df = self.spark.createDataFrame([
            ("John", 25, "Male"),
            ("Mary", 30, "Female")
        ], ["Name", "Age", "Gender"])

        process_file(self.spark, input_df, self.output_path)
        today = datetime.today().strftime("%Y-%m-%d")
        output_file = f"{self.output_path}/dt={today}/"
        self.assertTrue(os.path.isdir(output_file))
        self.assertEqual(len(os.listdir(output_file)), 1)
        self.assertEqual(os.listdir(output_file)[0], "_SUCCESS")
        output_count = self.spark.read.parquet(output_file).count()
        self.assertEqual(output_count, 2)

    @patch('boto3.Session.client')
    def test_perform_quality_checks(self, mock_client):
        output_path = 'path/to/output'
        resources_bucket = 'my-bucket'
        # Prepare mock response
        expected_config = {
            'datasources': {
                'mydatasource': {
                    'data_asset_type': 'SparkDFDataset',
                    'data_asset_options': {
                        'path': output_path,
                        'file_format': 'parquet'
                    }
                }
            }
        }
        mock_body = yaml.safe_dump(expected_config)
        mock_response = {'Body': mock_body}
        mock_client.return_value.get_object.return_value = mock_response

        # Call the function
        perform_quality_checks(output_path, resources_bucket)

        # Check if mock client was called with the correct parameters
        mock_client.assert_called_with('s3')
