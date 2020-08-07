package main

//go:generate swagger-codegen generate -i ../shared/api.yaml --lang go-server -Dmodels --output swagger

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/briggysmalls/detectordag/api/swagger"
	"github.com/briggysmalls/detectordag/api/swagger/server"
	"github.com/briggysmalls/detectordag/api/swagger/tokens"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/email"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"log"
)

var adapter *gorillamux.GorillaMuxAdapter

func init() {
	// Add file/line number to the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Create an AWS session
	// Good practice will share this session for all services
	sesh := shared.CreateSession(aws.Config{})
	// Create a new Db client
	db, err := database.New(sesh)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	// Create a new shadow client
	shadow, err := shadow.New(sesh)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	// Create a new iot client
	iot, err := iot.New(sesh)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	// Create a new session just for emailing (there is no emailing service in eu-west-2)
	emailSesh := shared.CreateSession(aws.Config{Region: aws.String("eu-west-1")})
	// Create a new email client
	email, err := email.New(emailSesh)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	// Create the tokens
	tokens := createTokens()
	// Create the server
	s := server.New(db, shadow, email, iot, tokens)
	// Create the router
	r := swagger.NewRouter(iot, s, createTokens())
	// Create an adapter for aws lambda
	adapter = gorillamux.New(r)
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

func createTokens() tokens.Tokens {
	// Load config from environment
	c, err := loadConfig()
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	// Get the token duration
	tokenDuration, _ := c.ParseDuration()
	// Create a tokens
	return tokens.New(c.JwtSecret, tokenDuration)
}
