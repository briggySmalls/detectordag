package app

import (
	"context"
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

// HandleRequest handles a lambda call
func HandleRequest(ctx context.Context, event StatusUpdatedEvent) {
	// Update the device status in the database
	device, err := getDevice(event.DeviceId)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Account ID: %d", device.AccountId)
	// Get the account
	account, err := getAccount(device.AccountId)
	if err != nil {
		log.Fatal(err)
	}
	// Construct an event to pass to the emailer
	update := PowerStatusChangedEmailConfig{
		DeviceId:  event.DeviceId,
		Timestamp: time.Unix(event.Updated.Status.Timestamp, 0),
		Status:    event.State.Status,
	}
	// Send 'power status updated' emails
	log.Printf("Send emails to: %s", account.Emails)
	err = SendEmail(account.Emails, update)
	if err != nil {
		log.Fatal(err)
	}
}
