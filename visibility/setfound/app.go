package main

import (
	"context"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/visibility"
	"log"
	"time"
)

type app struct {
	iot   iot.Client
	email visibility.EmailClient
}

type timestamp struct {
	Timestamp int64
}

type updated struct {
	Status timestamp
}

// DeviceSeenEvent tells us when a device has last been seen
type DeviceSeenEvent struct {
	DeviceId string  `json:""`
	Updated  updated `json:""`
}

// handleRequest handles a lambda call
func (a *app) handleRequest(ctx context.Context, event DeviceSeenEvent) error {
	var err error
	// Get the current device state
	device, err := a.iot.GetThing(event.DeviceId)
	if err != nil {
		return err
	}
	// Check if it is marked as lost
	if device.Visibility {
		// Short-circuit (it's already marked as visible)
		return nil
	}
	// Update status
	status := true
	err = a.iot.SetVisibiltyState(device.DeviceId, status)
	if err != nil {
		return err
	}
	// Send emails to indicate visibility status was updated
	return a.email.SendVisibilityStatus(
		device,
		time.Unix(event.Updated.Status.Timestamp, 0),
		status,
	)
}
