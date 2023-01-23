# parse the args

# Read header of file

# check if header has '|' if so set delimiter
# check if header has ',' if so set delimiter
# check if header has '\t' if so set delimiter

# check if header has at least 2 " if so set quotes on

# Get a count of the raw file rows, log it and store as variable

# read the input file

# write to parquet, allow overwrite. partitioned by /FILENAME/dt=YYYY-MM-DD/

# Log output file path

# Get a count of the output file rows, log it along with the input count

# mock out as comments: Run great expectations check (how to get results data into postgres?)

# mock out as comments: make call to API to move on to next steps
#  (update posgres table, trigger glue crawlers, etc)
if __name__ == "__main__":
    print("hello spark")


