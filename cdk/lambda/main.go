package main

import (
		"fmt"
        "context"
        "github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
    Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

func HandleRequest(ctx context.Context, event Event) (string, error) {
	bucketName := event.Records[0].S3.Bucket.Name
	key := event.Records[0].S3.Object.Key

	// Log job started in postgress
	// Fire off spark job
	// Spark ingest job then logs success and stats in postgres
	// Spark ingest job then triggers data quality checks
	// data quality checks then log stats in postgres

	return fmt.Sprintf("File dropped in bucket %s, path %s", bucketName, key), nil
}

func main() {
        lambda.Start(HandleRequest)
}