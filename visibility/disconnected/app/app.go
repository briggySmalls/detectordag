package app

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/briggysmalls/detectordag/shared/sqs"
	"github.com/briggysmalls/detectordag/visibility"
)

type app struct {
	iot     iot.Client
	updater visibility.ConnectionUpdater
	shadow  shadow.Client
}

type App interface {
	Handler(ctx context.Context, sqsEvent events.SQSEvent) error
}

func New(
	updater visibility.ConnectionUpdater,
	iot iot.Client,
	shadow shadow.Client,
) App {
	return &app{
		iot:     iot,
		shadow:  shadow,
		updater: updater,
	}
}

// hander handles SQS events
// The messages all indicate a disconnected event, which we are debouncing
func (a *app) Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Handle SQS events
	for _, message := range sqsEvent.Records {
		if err := a.processMessage(message); err != nil {
			return err
		}
	}
	return nil
}

func (a *app) processMessage(message events.SQSMessage) error {
	// Deserialise the disconnection message
	var payload sqs.DisconnectedPayload
	err := json.Unmarshal([]byte(message.Body), &payload)
	if err != nil {
		return err
	}
	// Validate the parsed struct
	if err := payload.Validate(); err != nil {
		return err
	}
	// Get the current connection status
	connectionState, err := a.shadow.GetConnectionStatus(payload.DeviceID)
	if err != nil {
		return err
	}
	// Check if the connection status has changed in this time
	if connectionState.Updated.After(payload.Time) {
		// Some other status change got there first
		return nil
	}
	// Fetch the device
	device, err := a.iot.GetThing(payload.DeviceID)
	if err != nil {
		return err
	}
	// Send emails to indicate the disconnection
	return a.updater.UpdateConnectionStatus(device, payload.Time, false)
}
