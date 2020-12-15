package app

import (
	"context"
	"fmt"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/sqs"
	"github.com/briggysmalls/detectordag/connection"
	"log"
	"time"
)

type app struct {
	sqs     sqs.Client
	iot     iot.Client
	updater connection.ConnectionUpdater
}

type App interface {
	RunJob(ctx context.Context, event DeviceLifecycleEvent) error
}

func New(
	updater connection.ConnectionUpdater,
	iot iot.Client,
	sqs sqs.Client,
) App {
	return &app{
		updater: updater,
		sqs:     sqs,
		iot:     iot,
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

// RunJob handles a lambda call
// The visibility status approach depends on the following invariants:
// - Connection events will always alternate (never consecutive "connected" or "disconnected")
// - "Connected" events are always geniune
// - "Disconnected" events may be spurious (quickly followed with a "connected" event)
func (a *app) RunJob(ctx context.Context, event DeviceLifecycleEvent) error {
	// Print the event
	log.Printf("%v\n", event)
	eventTime := time.Unix(event.Timestamp/1000, 0).UTC()
	// Handle a connected event
	if event.EventType == LifecycleEventTypeConnected {
		// Get the device
		device, err := a.iot.GetThing(event.DeviceID)
		if err != nil {
			return err
		}
		// "Connected" is always trustworthy, so update directly
		return a.updater.UpdateConnectionStatus(device, eventTime, true)
	}
	// Handle a disconnected event
	if event.EventType == LifecycleEventTypeDisconnected {
		// Delay dealing with disconnected events, to debounce
		return a.sqs.QueueDisconnectedEvent(sqs.DisconnectedPayload{
			DeviceID: event.DeviceID,
			Time:     eventTime,
		})
	}
	// Something went wrong
	return fmt.Errorf("Unexpected lifecycle event: %s", event.EventType)
}
