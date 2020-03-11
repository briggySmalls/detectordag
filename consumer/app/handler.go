package app

import (
	"context"
	"log"
	"time"
)

type PowerStatusChangedEvent struct {
	DeviceId  string    `json:""`
	Timestamp time.Time `json:""`
	Version   string    `json:""`
	Status    bool      `json:""`
}

// HandleRequest handles a lambda call
func HandleRequest(ctx context.Context, event PowerStatusChangedEvent) {
	// // Parse the time in the event
	// eventTime, err := time.Parse(time.RFC3339, event.Timestamp)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// Get the device ID
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
		Timestamp: event.Timestamp,
		Status:    event.Status,
	}
	// Send 'power status updated' emails
	log.Printf("Send emails to: %s", account.Emails)
	err = SendEmail(account.Emails, update)
	if err != nil {
		log.Fatal(err)
	}
}
