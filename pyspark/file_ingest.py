import argparse
from datetime import datetime
import logging
import re


from pyspark.sql.session import SparkSession

def create_spark_session(input_file):
   spark = (
        SparkSession.builder.appName(input_file)
        .enableHiveSupport()
        .config("mapreduce.fileoutputcommitter.algorithm.version", "2")
    )
    return spark

def inspect_file(spark,  input_file):
    # Read header of file
    # check if header has '|' if so set delimiter
    # check if header has ',' if so set delimiter
    # check if header has '\t' if so set delimiter
    rdd = spark.textFile(input_file).take(1)
    header = rdd.first()
    if header.count('|') > 1:
        delimiter = '|'
    elif header.count('|') > 1:
        delimiter = ','
    # check if header has at least 2 " if so set quotes on
    is_quoted = header.count('"') > 2
    
    # check if file is compressed
    compression = None
    if 'gzip' in file_extension:
        compression = 'gzip
        
    return delimiter, is_quoted, compression

def process_file(spark, delimiter, is_quoted, compression, input_file, output_path):
    # read the input file, header is required
    df = spark.read.csv(input_file, sep=delimiter, header=True, quoteAll=is_quoted, inferSchema=True, compression=compression)
    input_count = df.count()
    logging.info(f"Reading file {input_file} with {input_count} rows")
    
    # write to parquet, allow overwrite. partitioned by /FILENAME/dt=YYYY-MM-DD/
    today = datetime.today().strftime('%Y-%m-%d')
    output_path = f"{output_path}/{}/dt={today}/"
    df.write.parquet(path=output_path, mode='overwrite')

    # Log output file path
    logging.info(f"Wrote data to {output_path}")
    
    # QA counts
    output_count = spark.read.parquet(output_path).count()
    logging.info(f"Read {input_count} rows, wrote {output_count} rows.")
    
    
def perform_quality_checks():
    pass

def notify_downstream():
    pass    

if __name__ == "__main__":
    # Read args
    parser = argparse.ArgumentParser(description="Spark Job Arguments")
    parser.add_argument("--input_file", required=True, help="Input S3 Path")
    parser.add_argument("--output_path", required=True, help="Output S3 Path")
    parser.add_argument("--file_extension", required=True, help="File Extension")
    args = parser.parse_args()

    # Parse args
    input_file = args.input_file
    output_path = args.output_path
    file_extension = args.file_extension
    logging.info(f"Processing file: {input_file}")
    
    # Create spark session
    spark = create_spark_session(input_file)
    
    # Gather data about file format
    delimiter, is_quoted, compression = inspect_file(spark, input_file)
    
    # Convert file to parquet
    process_file(spark, delimiter, is_quoted, compression, input_file, output_path)

    # Run Quality Checks and log results
    perform_quality_checks(output_path)

    # mock out as comments: make call to API to move on to next steps (update postgres table, trigger glue crawlers, etc)
    notify_downstream()
