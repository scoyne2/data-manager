[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://github.com/codespaces/new?hide_repo_select=true&ref=main&repo=540243984)

[![codecov](https://codecov.io/gh/scoyne2/data-manager/branch/main/graph/badge.svg?token=GWRFE77EFG)](https://codecov.io/gh/scoyne2/data-manager)
# Data Manager
Data manager is a platform that allows ingestion of flat files into a raw layer of a datalake. It requires use of AWS and uses the following technology: `terraform`, `helm`, `go`, `nextjs`, `react`, `typescript`.
<img width="1728" alt="Data Manager Demo" src="docs/images/demo.png">

## Codespaces
If you are running this in codespaces you can skip the below homebrew install. Create [codespace secrets](https://docs.github.com/en/codespaces/managing-your-codespaces/managing-encrypted-secrets-for-your-codespaces) for `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` You will still need to complete the `Env Setup` steps.

## S3 Setup
Prior to starting you must create 3 S3 buckets in the account you plan to deploy to. These buckets will be used for input of raw files, output of processed data and metadata used during processing. Update the `.env` file with your bucket names. Only the `TRIGGER_BUCKET_NAME` requires setup. The Trigger Bucket is the bucket is where internal/external parties will drop flat files. When dropping files, you must follow the below S3 path `s3://{trigger_bucket}/inbound/{vendor}/{feed}/`. For example if your vendor is named `Coyne Enterprises` and they are going to send files related to `Orders` you should have them drop the files in `s3://data-manager-trigger/inbound/coyne_enterprises/orders/`. You would then [grant cross-account access](https://repost.aws/knowledge-center/cross-account-access-s3) on `s3://data-manager-trigger/inbound/coyne_enterprises/` to the Coyne Enterprises AWS account so that they could drop files as needed.

## Getting Started
Installation of dependencies and how they should be setup/configured. 

### Prerequisites
1. Homebrew
   * ``/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"``
2. AWS CLI
   * ``brew install awscli``
3. Terraform
   * ``brew install terraform``
4. Configure AWS
   * Setup ``~/.aws/config`` and ``~/.aws/credentials`` with the expected profile. You can run ``aws configure`` if you need to set this up.

### Env Setup
1. Make a copy of ``.env_template``, rename it to ``.env`` and update all necessary fields. In particular DOMAIN_NAME must align with the domain name you own via Route 53, AWS_ACCOUNT must be the account ID youwill be setting up resources in. AWS_PROFILE must be the name of the profile in your local ~/.aws/credentials that you will use to create resources (It will need access to all of the systems that Terraform will build)
2. To start you must had a domain name with AWS Route 53, and have an active Hosted Zone.
Take the Hosted Zone Id and update ``HOSTED_ZONE_ID`` in your .env file.
3. From project root run command ``make init`` then plan the creation of infra  ``make plan`` and finally ``make apply`` to build the infra. This will also output an ACM ARN to your ``.env`` file.
4. You can push images to ECR, update the kubeconfig and deploy the helm chart all with one command. From the project root run ``make deploy``

### Front End
The front end uses HTTPS. In your browser navigate to https://app.YOURDOMAIN.com/

### GraphQL API
The api end uses HTTPS. In your browser navigate to https://api.YOURDOMAIN/sandbox

### Database Admin
The db admin end uses HTTP. In your browser navigate to http://pgadmin4.YOURDOMAIN.com
Use the credentials you set in your ``.env`` file.


## Testing
* To run tests and upload coverage run `make test`

## Additional Documentation
### GraphQL API
#### Queries
The `feedstatusesdetailed` endpoint returns detailed results about all feeds. It does not accept any parameters.
This endpoint returns a list of type `FeedStatusDetails` which has nine fields which can be found [here](https://github.com/scoyne2/data-manager/blob/d7b1434acfeb3ff281aad73846f7d5d2aae0b998/api/schemas/factory.go#L55).


The `feedstatuseaggregates` endpoint returns aggregated statistics about all feeds. This query requires the following parameters.
- `startDate`: The start of the date range to query aggregate results for.
- `endDate`: The end of the date range to query aggregate results for.
This endpoint returns type `FeedStatusAggregate` which has three fields which can be found [here](https://github.com/scoyne2/data-manager/blob/d7b1434acfeb3ff281aad73846f7d5d2aae0b998/api/schemas/factory.go#L98).

#### Mutations
The `addFeed` endpoint is used to add a new feed to the Postgres database. This query requires the following parameters.
- `vendor`: The name of the Vendor who dropped the file.
- `feedName`: The name of the feed for the file which was dropped.
- `feedMethod`: The way the file was dropped. Either `SFTP` or `S3`.
This endpoint returns either the string `Feed Added` or an error.


The `updateFeedStatus` endpoint is used to update the status of a feed. This query requires the following parameters.
- `processDate`: The timestamp the feed status changed.
- `recordCount`: The total number of records from the file.
- `errorCount`: The number of records from the file which could not be loaded.
- `status`: The status of the file. Options include: `Received`, `Processing`, `Validating`, `Success`, `Failed`, `Errors`
- `fileName`: The name of the file that was dropped.
- `vendor`: The name of the Vendor who dropped the file.
- `feedName`: The name of the feed for the file which was dropped.
- `emrApplicationID`: The ID of the EMR application which processed the file. Used for generating the EMR log URL.
- `emrStepID`: The ID of the EMR Step which processed the file. Used for generating the EMR log URL.
This endpoint returns either the string `Feed status updated for id: ID` or an error.

### REST API
The data preview functionality is accessed via REST API. There is one HTTP GET endpoint which accepts values via query parameters. You must pass: `vendor`, `feedname` and `filename`. Example:

```
https://api.datamanagertool.com/preview?vendor=Coyne+Enterprises&feedname=Hello&filename=hello.csv
```

### Ingestion Process
The DataManagerTool processed flat files and outputs data in S3 as parquet + a glue table for easy access. The ingestion process is as follows.

- A customer drops a flat file (.csv or .txt) in the appropriate S3 path. The path must follow convention `s3://{input_bucket}/inbound/vendor/feed`. The S3 PUT event triggers an AWS Lambda.
- The AWS Lambda is triggered by the S3 PUT event. It first updates the FeedStatus in the postgres database, and if necessary creates a new feed in the `feeds` table. It then kisks off a Spark Job using EMR Serverless.
- The Spark job parses the flat file, automatically detecting the delimiter adn if the fields are quoted, and converts it to a spark dataframe, and then writes it in the parquet format to the processed bucket. The output path will follow convention `s3://{OUTPUT_BUCKET}/{vendor}/{feed}/` partitioned by todays date as `/dt=YYYY-MM-DD/` .
- The Spark job also creates an AWS Glue Database following naming convention `data_manager_output_{vendor}.{feed_name}`.
The Spark job then runs QA checks, and updates the final status of the feed in the postgres database.