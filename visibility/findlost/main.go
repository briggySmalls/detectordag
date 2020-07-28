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
		return err
	}
	// Iterate through devices
	for _, device := range devices {
		// Fetch the shadow
		shdw, err := shadowClient.Get(device.DeviceId)
		if err != nil {
			return err
		}
		_, ok := shdw.State.Reported["status"].(bool)
		if !ok {
			return fmt.Errorf("%s doesn't have status", deviceString(device))
		}
		lastSeen := shdw.Metadata.Reported["status"].Timestamp.Time
		// Device hasn't been seen for a while
		if time.Now().Before(lastSeen.Add(lastSeenDuration)) {
			// This device was seen recently enough
			continue
		}
		if !device.Visibility {
			// We searched for visible devices, something weird has happened
			log.Printf("%s already marked lost despite searching for visible devices", deviceString(device))
			continue
		}
		// The device is lost
		err = handleLostDevice(device, lastSeen)
		if err != nil {
			return err
		}
	}
	// Return the response
	return err
}

// main is the entrypoint to the lambda function
func main() {
	lambda.Start(runJob)
}

func deviceString(device *iot.Device) string {
	return fmt.Sprintf("Device '%s' ('%s')", device.DeviceId, device.Name)
}

func handleLostDevice(device *iot.Device, lastSeen time.Time) error {
	log.Printf("%s not seen since %s", deviceString(device), lastSeen.Format(time.RFC3339))
	// Get the account
	account, err := dbClient.GetAccountById(device.AccountId)
	if err != nil {
		return err
	}
	// Notify the account owner their device is missing
	err = visibility.SendEmail(
		account.Emails,
		visibility.VisibilityStatusChangedEmailConfig{
			DeviceName: device.Name,
			Timestamp:  lastSeen,
			Status:     false,
		},
	)
	if err != nil {
		return err
	}
	// Mark as lost
	return iotClient.SetVisibiltyState(device.DeviceId, false)
}
