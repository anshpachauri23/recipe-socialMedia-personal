package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Function started successfully")
	log.Printf("Request method: %s", request.HTTPMethod)
	log.Printf("Request path: %s", request.Path)
	
	response := map[string]interface{}{
		"message": "Test function working",
		"method":  request.HTTPMethod,
		"path":    request.Path,
	}
	
	responseBody, _ := json.Marshal(response)
	
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func main() {
	log.Printf("Starting test Lambda function")
	lambda.Start(handler)
}
