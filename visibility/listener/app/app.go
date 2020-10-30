package app

import (
	"context"
	"fmt"
	"github.com/briggysmalls/detectordag/shared/iot"
	"log"
	"time"
)

type app struct {
	iot   iot.Client
	email EmailClient
}

type App interface {
	RunJob(ctx context.Context, event DeviceLifecycleEvent) error
}

func New(
	iot iot.Client,
	email EmailClient,
) App {
	return &app{
		iot:   iot,
		email: email,
	}
}

const (
	LifecycleEventTypeConnected    = "connected"
	LifecycleEventTypeDisconnected = "disconnected"
)

// DeviceLifecycleEvent tells us when a device has last been seen
type DeviceLifecycleEvent struct {
	DeviceID  string `json:"clientId"`
	EventType string `json:""`
	Timestamp int64  `json:""`
}

// handleRequest handles a lambda call
func (a *app) RunJob(ctx context.Context, event DeviceLifecycleEvent) error {
	// Print the event
	log.Printf("%v\n", event)
	// Determine the event
	var visibility bool
	if event.EventType == LifecycleEventTypeConnected {
		visibility = true
	} else if event.EventType == LifecycleEventTypeDisconnected {
		visibility = false
	} else {
		return fmt.Errorf("Unexpected lifecycle event: %s", event.EventType)
	}
	// Parse the time
	lastSeen := time.Unix(event.Timestamp/1000, 0).UTC()
	// Get the device
	device, err := a.iot.GetThing(event.DeviceID)
	if err != nil {
		return err
	}
	// Record the new device status
	err = a.iot.SetVisibiltyState(device.DeviceId, visibility)
	if err != nil {
		return err
	}
	// Indicate the device status has changed
	err = a.email.SendVisibilityStatus(device, lastSeen, visibility)
	if err != nil {
		return err
	}
	return nil
}
