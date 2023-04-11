import unittest
from unittest.mock import patch

from pyspark.sql.session import SparkSession

from file_ingest import (
    create_spark_session,
    # inspect_file,
    # process_file,
    # perform_quality_checks,
)


class TestMyScript(unittest.TestCase):
    def setUp(self):
        self.spark = SparkSession.builder.appName("TestInspectFile").getOrCreate()

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
