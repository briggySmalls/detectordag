package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	swagger "github.com/briggysmalls/detectordag/api/go"
	"log"
)

var adapter *gorillamux.GorillaMuxAdapter

func init() {
	// Create the server
	server := swagger.NewRouter()
	// Create an adapter for aws lambda
	adapter = gorillamux.New(server)
}

// HandleRequest handles a lambda call
func handleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Print out hander paramters out details
	log.Print("Request body: ", event)
	log.Print("Context: ", ctx)
	// Pass the request to the adapter
	response, err := adapter.ProxyWithContext(ctx, event)
	// Return the response
	return response, err
}

// main is the entrypoint to the lambda function
func main() {
	lambda.Start(handleRequest)
}
