package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.S3Event) (events.APIGatewayProxyResponse, error) {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("File from bucket %t\n", request.Records[0].S3.Bucket)
	log.Printf("File from key %t\n", request.Records[0].S3.Object.Key)
	return events.APIGatewayProxyResponse{
		Body:       "Hello world",
		StatusCode: 200,
	}, nil
}

func main() {
	log.Printf("Start lambda")
	lambda.Start(Handler)
}