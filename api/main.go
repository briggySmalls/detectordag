package main

//go:generate swagger-codegen generate -i api.yaml --lang go-server -Dmodels --output swagger

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/briggysmalls/detectordag/api/swagger"
	"github.com/briggysmalls/detectordag/api/swagger/server"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/email"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"log"
)

var adapter *gorillamux.GorillaMuxAdapter

func init() {
	// Get config from environment
	c, err := server.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create an AWS session
	// Good practice will share this session for all services
	sesh, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new Db client
	db, err := database.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new shadow client
	shadow, err := shadow.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new session just for emailing (we have to use a different region)
	emailSesh, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			// There is no emailing service in eu-west-2
			Region: aws.String("eu-west-1"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	// Create a new email client
	email, err := email.New(emailSesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create the server
	serverParams := server.Params{
		Db:     db,
		Shadow: shadow,
		Email:  email,
		Config: *c,
	}
	s := server.New(serverParams)
	// Create the router
	router := swagger.NewRouter(s)
	// Create a router
	// Create an adapter for aws lambda
	adapter = gorillamux.New(router)
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
