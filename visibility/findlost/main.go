package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/briggysmalls/detectordag/visibility"
	"github.com/briggysmalls/detectordag/visibility/findlost/app"
	"log"
	"time"
)

const lastSeenDurationHours = 24

// Prepare an application to reuse across lambda runs
var findLost app.App

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
	// Create a new shadow client
	shadowClient, err := shadow.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new Db client
	dbClient, err := database.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new session just for emailing (there is no emailing service in eu-west-2)
	emailSesh := shared.CreateSession(aws.Config{Region: aws.String("eu-west-1")})
	// Create a new visibility email client
	visibilityEmailClient, err := visibility.New(emailSesh, dbClient)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create the duration
	lastSeenDuration, err := time.ParseDuration(fmt.Sprintf("%dh", lastSeenDurationHours))
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create the application
	findLost = app.New(iotClient, visibilityEmailClient, shadowClient, lastSeenDuration)
}

// main is the entrypoint to the lambda function
func main() {
	lambda.Start(findLost.RunJob)
}
