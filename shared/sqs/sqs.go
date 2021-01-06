package sqs

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/briggysmalls/detectordag/shared"
)

type client struct {
	sqs      sqsiface.SQSAPI
	queueUrl string
}

type ConnectionEventPayload struct {
	DeviceID string    `json:"deviceId" validate:"uuid"`
	Status   string    `json:"type" validate:"eq=connected|eq=disconnected"`
	Time     time.Time `json:"time" validate:"required"`
	ID       string    `json:"id" validate:"uuid"`
}

func (d *ConnectionEventPayload) Validate() error {
	return shared.Validate.Struct(d)
}

// Client is a client for sending status updates to the queue
type Client interface {
	QueueConnectionEvent(payload ConnectionEventPayload) error
}

// NewSender gets a new Client
func New(sesh *session.Session, queueUrl string) (Client, error) {
	// Create Amazon DynamoDB client
	sqs := sqs.New(sesh)
	if sqs == nil {
		return nil, errors.New("Failed to create SQS client")
	}
	// Create our client wrapper
	client := client{
		sqs:      sqs,
		queueUrl: queueUrl,
	}
	return &client, nil
}

func (c *client) QueueConnectionEvent(payload ConnectionEventPayload) error {
	// Ensure the struct is valid
	if err := shared.Validate.Struct(payload); err != nil {
		return err
	}
	// Marshal the payload to a string
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	// Send the message
	_, err = c.sqs.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl:    aws.String(c.queueUrl),
	})
	return err
}
