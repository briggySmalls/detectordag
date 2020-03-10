package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/briggysmalls/detectordag/consumer/app"
	"log"
)

func init() {
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	if session == nil {
		log.Fatal("Failed to start session")
	}
	// Initialise various AWS clients here in case the container is reused
	var err error
	err = app.DbInit(session)
	if err != nil {
		log.Fatal(err)
	}
	err = app.EmailInit(session)
	if err != nil {
		log.Fatal(err)
	}
}

// main is the entrypoint to the lambda function
func main() {
	lambda.Start(app.HandleRequest)
}
