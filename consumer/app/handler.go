package app

import (
	"context"
	"log"
)

type PowerStatusChangedEvent struct {
	DeviceId  string `json:""`
	Timestamp string `json:""`
	Version   string `json:""`
	Status    bool   `json:""`
}

// HandleRequest handles a lambda call
func HandleRequest(ctx context.Context, event PowerStatusChangedEvent) {
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
	// Send 'power status updated' emails
	for _, email := range account.Emails {
		log.Printf("Send email to: %s", email)
	}
}
