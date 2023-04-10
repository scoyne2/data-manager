import unittest
from unittest.mock import patch

from pyspark.sql.session import SparkSession

from my_script import create_spark_session, inspect_file, process_file, perform_quality_checks


class TestMyScript(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.spark = SparkSession.builder.appName("test").enableHiveSupport().getOrCreate()

    def test_create_spark_session(self):
        spark = create_spark_session("test")
        self.assertIsInstance(spark, SparkSession)

    def test_inspect_file(self):
        # Test file with delimiter "|"
        delimiter, quote = inspect_file(self.spark, "tests/fixtures/test_file1.csv")
        self.assertEqual(delimiter, "|")
        self.assertEqual(quote, "")

        # Test file with delimiter ","
        delimiter, quote = inspect_file(self.spark, "tests/fixtures/test_file2.csv")
        self.assertEqual(delimiter, ",")
        self.assertEqual(quote, "")

        # Test file with quoted values
        delimiter, quote = inspect_file(self.spark, "tests/fixtures/test_file3.csv")
        self.assertEqual(delimiter, ",")
        self.assertEqual(quote, '"')

    def test_process_file(self):
        with patch("my_script.logging.info") as mock_logging_info:
            # Test reading and writing to parquet
            process_file(self.spark, ",", "", "tests/fixtures/test_file2.csv", "tests/fixtures/output/")
            mock_logging_info.assert_called_with("Wrote data to tests/fixtures/output/dt=")

            # Test file with quoted values
            process_file(self.spark, ",", '"', "tests/fixtures/test_file3.csv", "tests/fixtures/output/")
            mock_logging_info.assert_called_with("Wrote data to tests/fixtures/output/dt=")

    @patch("my_script.boto3.Session")
    def test_perform_quality_checks(self, mock_boto_session):
        mock_client = mock_boto_session.return_value.client.return_value
        mock_client.get_object.return_value = {"Body": "config file contents"}
        perform_quality_checks("tests/fixtures/output/", "test-bucket")
        mock_client.get_object.assert_called_with(
            Bucket="test-bucket",
            Key="great_expectations/great_expectations.yml",
        )
