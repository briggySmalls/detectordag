package app

import (
	"context"
	"fmt"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/briggysmalls/detectordag/shared/sqs"
	"log"
	"time"
)

type app struct {
	shadow shadow.Client
	sqs    sqs.Client
}

type App interface {
	RunJob(ctx context.Context, event DeviceLifecycleEvent) error
}

func New(
	shadow shadow.Client,
	sqs sqs.Client,
) App {
	return &app{
		shadow: shadow,
		sqs:    sqs,
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
	// Handle a connected event
	if event.EventType == LifecycleEventTypeConnected {
		// "Connected" is always trustworthy, so update directly
		return a.shadow.UpdateConnectionStatus(event.DeviceID, true)
	}
	// Handle a disconnected event
	if event.EventType == LifecycleEventTypeDisconnected {
		// Delay dealing with disconnected events, to debounce
		return a.sqs.QueueDisconnectedEvent(sqs.DisconnectedPayload{
			DeviceID: event.DeviceID,
			Time:     time.Unix(event.Timestamp/1000, 0).UTC(),
		})
	}
	// Something went wrong
	return fmt.Errorf("Unexpected lifecycle event: %s", event.EventType)
}
