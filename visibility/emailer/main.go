package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/visibility/listener/app"
	"log"
	"os"
)

// Prepare an application to reuse across lambda runs
var emailer app.App

func init() {
	// Add file/line number to the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	// Create an AWS session
	// Good practice will share this session for all services
	sesh := shared.CreateSession(aws.Config{})
	// Create a new iot client
	iotClient, err := iot.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new iot shadow client
	shadowClient, err := shadow.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new database client
	dbClient, err := database.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new session just for emailing (there is no emailing service in eu-west-2)
	emailSesh := shared.CreateSession(aws.Config{Region: aws.String("eu-west-1")})
	// Create a new visibility email client
	visibilityEmailClient, err := app.NewVisibilityEmailer(emailSesh, dbClient)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create the application
	emailer = app.New(iotClient, shadowClient, visibilityEmailClient)
}

// main is the entrypoint to the lambda function
func main() {
	lambda.Start(emailer.RunJob)
}
