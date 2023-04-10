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
    @classmethod
    def setUpClass(cls):
        cls.spark = (
            SparkSession.builder.appName("test").enableHiveSupport().getOrCreate()
        )

    def test_create_spark_session(self):
        spark = create_spark_session("test")
        self.assertIsInstance(spark, SparkSession)
