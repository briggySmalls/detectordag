package app

import (
	"context"
	"log"
	"time"

	"github.com/briggysmalls/detectordag/connection"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/briggysmalls/detectordag/shared/sqs"
	"github.com/google/uuid"
)

type app struct {
	sqs     sqs.Client
	iot     iot.Client
	shadow  shadow.Client
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
	// Always update the transient state
	id := uuid.New().String()
	a.shadow.UpdateConnectionTransientID(event.DeviceID, id)
	// Get the device shadow
	shdw, err := a.shadow.Get(event.DeviceID)
	if err != nil {
		return err
	}
	// Check if we need to enqueue a handler
	if event.EventType != shdw.Connection.Status {
		// The status has changed, enqueue a callback
		if err := a.sqs.QueueConnectionEvent(sqs.ConnectionEventPayload{
			DeviceID: event.DeviceID,
			Status:   event.EventType,
			Time:     eventTime,
			ID:       id,
		}); err != nil {
			return err
		}
	}
	// Something went wrong
	return nil
}
