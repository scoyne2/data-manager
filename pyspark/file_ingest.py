import argparse
from datetime import datetime
import logging
import re


from pyspark.sql.session import SparkSession

if __name__ == "__main__":
    print("hello spark")

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
    
    spark = (
        SparkSession.builder.appName(input_file)
        .enableHiveSupport()
        .config("mapreduce.fileoutputcommitter.algorithm.version", "2")
    )
    
    
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
    isQuoted = header.count('"') > 2
    
    # check if file is compressed
    compression = None
    if 'gzip' in file_extension:
        compression = 'gzip    

    # Get a count of the raw file rows, log it and store as variable

    # read the input file, header is required
    df = spark.read.csv(input_file, sep=delimiter, header=True, quoteAll=isQuoted, inferSchema=True, compression=compression)

    # write to parquet, allow overwrite. partitioned by /FILENAME/dt=YYYY-MM-DD/
    today = datetime.today().strftime('%Y-%m-%d')
    output_path = f"{output_path}/{}/dt={today}/"
    df.write.parquet(path=output_path, mode='overwrite')

    # Log output file path

    # Get a count of the output file rows, log it along with the input count

    # mock out as comments: Run great expectations check (how to get results data into postgres?)

    # mock out as comments: make call to API to move on to next steps
    #  (update posgres table, trigger glue crawlers, etc)
