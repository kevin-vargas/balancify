package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/grokify/go-awslambda"
)

type customStruct struct {
	Content       string
	FileName      string
	FileExtension string
}

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	res := events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
	}
	r, err := awslambda.NewReaderMultipart(req)
	if err != nil {
		return res, err
	}
	part, err := r.NextPart()
	if err != nil {
		return res, err
	}
	content, err := io.ReadAll(part)
	if err != nil {
		return res, err
	}
	awsRegion := os.Getenv("AWS_REGION")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to create AWS session",
		}, nil
	}
	svc := s3.New(sess)
	bucketName := os.Getenv("S3_BUCKET_NAME")
	input := &s3.PutObjectInput{
		Body:   bytes.NewReader(content),
		Bucket: aws.String(bucketName),
		Key:    aws.String(part.FileName()),
	}
	_, err = svc.PutObject(input)
	if err != nil {
		log.Printf("Failed to upload to S3: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to upload to S3",
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("File '%s' uploaded successfully to bucket '%s'", part.FileName(), bucketName),
	}, nil
}

func main() {
	lambda.Start(handleRequest)

}
