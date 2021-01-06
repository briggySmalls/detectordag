package app

import (
	"context"
	"log"
	"time"

	"github.com/briggysmalls/detectordag/connection"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/briggysmalls/detectordag/shared/sqs"
	"github.com/google/uuid"
)

type app struct {
	sqs    sqs.Client
	shadow shadow.Client
}

type App interface {
	RunJob(ctx context.Context, event DeviceLifecycleEvent) error
}

func New(
	updater connection.ConnectionUpdater,
	shadow shadow.Client,
	sqs sqs.Client,
) App {
	return &app{
		sqs:    sqs,
		shadow: shadow,
	}
}

// DeviceLifecycleEvent tells us when a device has last been seen
type DeviceLifecycleEvent struct {
	DeviceID  string `json:"clientId"`
	EventType string `json:""`
	Timestamp int64  `json:""`
}

// RunJob handles a lambda call
// The visibility status approach depends on the following invariants:
// - Connection events will always alternate (never consecutive "connected" or "disconnected")
// - "Disconnected" events may be spurious (quickly followed with a "connected" event)
func (a *app) RunJob(ctx context.Context, event DeviceLifecycleEvent) error {
	// Print the event
	log.Printf("%v\n", event)
	// Prepare and validate the event
	eventTime := time.Unix(event.Timestamp/1000, 0).UTC()
	id := uuid.New().String()
	connectionEventPayload := sqs.ConnectionEventPayload{
		DeviceID: event.DeviceID,
		Status:   event.EventType,
		Time:     eventTime,
		ID:       id,
	}
	if err := connectionEventPayload.Validate(); err != nil {
		return err
	}
	// Always update the transient state
	a.shadow.UpdateConnectionTransientID(event.DeviceID, id)
	// Get the device shadow
	shdw, err := a.shadow.Get(event.DeviceID)
	if err != nil {
		return err
	}
	// Check if we need to enqueue a handler
	if event.EventType != shdw.Connection.Status {
		// The status has changed, enqueue a callback
		if err := a.sqs.QueueConnectionEvent(connectionEventPayload); err != nil {
			return err
		}
	}
	// Something went wrong
	return nil
}
