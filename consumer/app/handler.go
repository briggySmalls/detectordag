package app

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/briggysmalls/detectordag/shared/database"
	iotp "github.com/briggysmalls/detectordag/shared/iot"
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
var iot iotp.Client

func init() {
	// Create an AWS session
	// Good practice will share this session for all services
	sesh, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	// Create a database client
	db, err = database.New(sesh)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	// Create an IOT client
	iot, err = iotp.New(sesh)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
}

// HandleRequest handles a lambda call
func HandleRequest(ctx context.Context, event StatusUpdatedEvent) {
	// Get the device
	device, err := iot.GetThing(event.DeviceId)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	accountID := device.AccountId
	log.Printf("Device '%s' associated with account '%s'", event.DeviceId, accountID)
	// Get the account
	account, err := db.GetAccountById(accountID)
	if err != nil {
		shared.LogErrorAndExit(err)
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
		shared.LogErrorAndExit(err)
	}
}
