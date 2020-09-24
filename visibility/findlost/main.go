package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/briggysmalls/detectordag/visibility"
	"log"
	"time"
)

const lastSeenDurationHours = 24

// Prepare some clients to reuse across lambda runs
var dbClient database.Client
var iotClient iot.Client
var shadowClient shadow.Client
var lastSeenDuration time.Duration

func init() {
	// Add file/line number to the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	// Create an AWS session
	// Good practice will share this session for all services
	sesh := shared.CreateSession(aws.Config{})
	// Create a new Db client
	dbClient, err = database.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new iot client
	iotClient, err = iot.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new shadow client
	shadowClient, err = shadow.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create the duration
	lastSeenDuration, err = time.ParseDuration(fmt.Sprintf("%dh", lastSeenDurationHours))
	if err != nil {
		log.Fatal(err.Error())
	}
}

// runJob starts the device checks
func runJob(ctx context.Context) error {
	// Print out handler parameters
	log.Print("Context: ", ctx)
	// Request all devices that are considered 'visible'
	devices, err := iotClient.GetThingsByVisibility(true)
	if err != nil {
		return shared.LogErrorAndReturn(err)
	}
	log.Printf("Checking %d devices for updated visibility status", len(devices))
	// Iterate through devices
	for _, device := range devices {
		// Fetch the shadow
		shdw, err := shadowClient.Get(device.DeviceId)
		if err != nil {
			return shared.LogErrorAndReturn(err)
		}
		// Check we have a reported status
		_, ok := shdw.State.Reported["status"].(bool)
		if !ok {
			return fmt.Errorf("%s doesn't have status", visibility.DeviceString(device))
		}
		// Check when the device was last seen
		lastSeen := shdw.Metadata.Reported["status"].Timestamp.Time
		if time.Now().Before(lastSeen.Add(lastSeenDuration)) {
			// This device was seen recently enough
			continue
		}
		if !device.Visibility {
			// We searched for visible devices, something weird has happened
			log.Printf("%s already marked lost despite searching for visible devices", visibility.DeviceString(device))
			continue
		}
		// Mark the device as lost
		err = iotClient.SetVisibiltyState(device.DeviceId, false)
		if err != nil {
			return shared.LogErrorAndReturn(err)
		}
		// Email to say so
		err = visibility.EmailVisiblityStatus(dbClient, device, lastSeen, false)
		if err != nil {
			return shared.LogErrorAndReturn(err)
		}
	}
	// Return the response
	return err
}

// main is the entrypoint to the lambda function
func main() {
	lambda.Start(runJob)
}
