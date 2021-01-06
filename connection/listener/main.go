package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/briggysmalls/detectordag/connection"
	"github.com/briggysmalls/detectordag/connection/listener/app"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/briggysmalls/detectordag/shared/sqs"
)

const (
	senderEnvVar = "SENDER_EMAIL"
)

// Prepare an application to reuse across lambda runs
var listener app.App

func init() {
	// Add file/line number to the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	// Create an AWS session
	// Good practice will share this session for all services
	sesh := shared.CreateSession(aws.Config{})
	// Create a new shadow client
	shadowClient, err := shadow.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new SQS queue client
	sqsQueue, err := sqs.New(sesh, os.Getenv("DELAY_QUEUE_URL"))
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new database client
	dbClient, err := database.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Get the email sender
	sender := os.Getenv(senderEnvVar)
	if sender == "" {
		shared.LogErrorAndReturn(fmt.Errorf("Env var '%s' unset", senderEnvVar))
	}
	// Create a new session just for emailing (there is no emailing service in eu-west-2)
	emailSesh := shared.CreateSession(aws.Config{Region: aws.String("eu-west-1")})
	connectionUpdater, err := connection.NewConnectionUpdater(emailSesh, dbClient, shadowClient, sender)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create the application
	listener = app.New(connectionUpdater, shadowClient, sqsQueue)
}

// main is the entrypoint to the lambda function
func main() {
	lambda.Start(listener.RunJob)
}
