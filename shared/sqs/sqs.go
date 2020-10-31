package sqs

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"time"
)

type client struct {
	sqs      sqsiface.SQSAPI
	queueUrl string
}

type ConnectionStatusPayload struct {
	Connected bool      `json:"connected"`
	Time      time.Time `json:"time"`
}

// Client is a client for sending status updates to the queue
type Client interface {
	SendMessage(payload ConnectionStatusPayload) error
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

func (c *client) SendMessage(payload ConnectionStatusPayload) error {
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
