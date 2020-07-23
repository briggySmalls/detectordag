package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/email"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"log"
	"time"
)

const lastSeenDurationHours = 24

// Prepare some clients to reuse across lambda runs
var db database.Client
var iot iot.Client
var shadow shadow.Client
var email email.Client
var lastSeenDuration time.Duration

func init() {
	// Create an AWS session
	// Good practice will share this session for all services
	sesh := createSession(aws.Config{})
	var err error
	// Create a new Db client
	db, err = database.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new iot client
	iot, err := iot.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new shadow client
	shadow, err = shadow.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new session just for emailing (there is no emailing service in eu-west-2)
	emailSesh := createSession(aws.Config{Region: aws.String("eu-west-1")})
	// Create a new email client
	email, err = email.New(emailSesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create the duration
	lastSeenDuration, err = time.ParseDuration(fmt.Sprintf("%sh", lastSeenDurationHours))
	if err != nil {
		log.Fata(err.Error())
	}
}

// runJob starts the device checks
func runJob(ctx context.Context) error {
	// Print out handler parameters
	log.Print("Context: ", ctx)
	// Request all devices
	devices, err := iot.GetThings()
	if err != nil {
		return err
	}
	// Iterate through devices
	for _, device := range devices {
		// Fetch the shadow
		shdw, err := shadow.Get(device.DeviceId)
		if err != nil {
			return err
		}
		status, ok := shdw.State.Reported["status"].(bool)
		if !ok {
			return fmt.Error("Device '%s' doesn't have status", device.DeviceId)
		}
		lastSeen := shdw.State.Reported["status"].Timestamp.Time
		// Device hasn't been seen for a while
		if time.Now().Before(time.Add(lastSeen, lastSeenDuration)) {
			// This device was seen recently enough
			continue
		}
		// Notify the account owner their device is missing
		log.Printf("Device '%s' ('%s') not seen since %s", device.ID, device.Name, lastSeen.Format(time.RFC3339))
	}
	// Return the response
	return err
}

// main is the entrypoint to the lambda function
func main() {
	lambda.Start(runJob)
}
