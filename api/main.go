package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/briggysmalls/detectordag/api/app"
	"log"
	"net/http"
)

// HandleRequest handles a lambda call
func handleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// For now just print out details
	log.Print("Request body: ", event)
	log.Print("Context: ", ctx)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

// main is the entrypoint to the lambda function
func main() {
	lambda.Start(handleRequest)
}
