package main

//go:generate swagger-codegen generate -i api.yaml --lang go-server -Dmodels --output swagger

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/briggysmalls/detectordag/api/swagger"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/joho/godotenv"
	"log"
)

var adapter *gorillamux.GorillaMuxAdapter

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Get config from environment
	c, err := swagger.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new Db client
	db := database.New()
	// Create a new shadow client
	shadow := shadow.New()
	// Create the server
	server := swagger.NewRouter(c, db, shadow)
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
