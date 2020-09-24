package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/visibility"
	"log"
	"time"
)

type timestamp struct {
	Timestamp int64
}

type updated struct {
	Status timestamp
}

type DeviceSeenEvent struct {
	DeviceId string  `json:""`
	Updated  updated `json:""`
}

// Prepare some clients to reuse across lambda runs
var iotClient iot.Client
var dbClient database.Client

func init() {
	// Add file/line number to the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	// Create an AWS session
	// Good practice will share this session for all services
	sesh := shared.CreateSession(aws.Config{})
	// Create a new iot client
	iotClient, err = iot.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new Db client
	dbClient, err = database.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// handleRequest handles a lambda call
func handleRequest(ctx context.Context, event DeviceSeenEvent) error {
	// Get the current device state
	device, err := iotClient.GetThing(event.DeviceId)
	if err != nil {
		return shared.LogErrorAndReturn(err)
	}
	// Check if it is marked as lost
	if !device.Visibility {
		err = visibility.EmailVisiblityStatus(dbClient, device, time.Unix(event.Updated.Status.Timestamp, 0), true)
		if err != nil {
			return shared.LogErrorAndReturn(err)
		}
	}
	// Indicate we've now seen it
	return iotClient.SetVisibiltyState(event.DeviceId, true)
}

// main is the entrypoint to the lambda function
func main() {
	lambda.Start(handleRequest)
}
