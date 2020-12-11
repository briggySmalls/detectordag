package app

import (
	"/briggysmalls/detectordag/shared/iot"
	awsSqs "github.com/aws/aws-sdk-go/aws/sqs"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/briggysmalls/detectordag/shared/sqs"
)

type app struct {
	iot    iot.Client
	shadow shadow.Client
	sqs    sqs.Client
}

type App interface {
	RunJob(ctx context.Context, event DeviceLifecycleEvent) error
}

func New(
	iot iot.Client,
	shadow shadow.Client,
	email EmailClient,
) App {
	return &app{
		iot:    iot,
		shadow: shadow,
		sqs:    sqs,
	}
}

// hander handles SQS events
// The messages all indicate a disconnected event, which we are debouncing
func (a *app) handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Handle SQS events
	for _, message := range sqsEvent.Records {
		a.processMessage(message)
	}
	return nil
}

func (a *app) processMessage(message awsSqs.Message) error {
	// Deserialise the data
	var payload sqs.ConnectionStatusPayload
	err := json.Unmarshal([]byte(message.String()), &payload)
	if err != nil {
		return err
	}
	// Get the current device shadow
	thing, err = a.shadow.Get(payload.DeviceId)
	if err != nil {
		return err
	}
	// Check if the two values agree (successful debouncing)
	if thing.Visibility != payload.Connected {
		continue
	}
	// Send emails to indicate a status change
	err := a.email.SendVisibilityStatus(device.DeviceId, status)
	if err != nil {
		return err
	}
}
