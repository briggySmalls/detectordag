package app

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/briggysmalls/detectordag/shared/database"
	"log"
	"time"
)

type state struct {
	Status bool
}

type timestamp struct {
	Timestamp int64
}

type updated struct {
	Status timestamp
}

type StatusUpdatedEvent struct {
	DeviceId  string  `json:""`
	Timestamp int     `json:""`
	State     state   `json:""`
	Updated   updated `json:""`
}

var db database.Client

func init() {
	// Create an AWS session
	// Good practice will share this session for all services
	sesh, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a database client
	db, err = database.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// HandleRequest handles a lambda call
func HandleRequest(ctx context.Context, event StatusUpdatedEvent) {
	// Update the device status in the database
	device, err := db.GetDeviceById(event.DeviceId)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Account ID: %d", device.AccountId)
	// Get the account
	account, err := db.GetAccountById(device.AccountId)
	if err != nil {
		log.Fatal(err)
	}
	// Construct an event to pass to the emailer
	update := PowerStatusChangedEmailConfig{
		DeviceName: device.Name,
		Timestamp:  time.Unix(event.Updated.Status.Timestamp, 0),
		Status:     event.State.Status,
	}
	// Send 'power status updated' emails
	log.Printf("Send emails to: %s", account.Emails)
	err = SendEmail(account.Emails, update)
	if err != nil {
		log.Fatal(err)
	}
}
