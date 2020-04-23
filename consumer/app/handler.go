package app

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/briggysmalls/detectordag/shared/database"
	iotm "github.com/briggysmalls/detectordag/shared/iot"
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
var iot iotm.Client

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
	// Create an IOT client
	iot, err = iotm.New(sesh)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// HandleRequest handles a lambda call
func HandleRequest(ctx context.Context, event StatusUpdatedEvent) {
	// Get the device
	device, err := iot.GetThing(event.DeviceId)
	if err != nil {
		log.Fatal(err)
	}
	// Get the account ID from the attributes
	accountID, err := device.AccountID()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Device '%s' associated with account '%d'", event.DeviceId, accountID)
	// Get the account
	account, err := db.GetAccountById(accountID)
	if err != nil {
		log.Fatal(err)
	}
	// Construct an event to pass to the emailer
	name, err := device.Name()
	if err != nil {
		log.Fatal(err)
	}
	update := PowerStatusChangedEmailConfig{
		DeviceName: name,
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
