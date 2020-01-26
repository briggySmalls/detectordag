package main

import (
	"log"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/streadway/amqp"
)

type Messenger interface {
	Connect(address string) error
	StatusUpdate(status int) error
}

type messenger struct {
	connection  amqp.Connection
	channel     amqp.Channel
	statusQueue amqp.Queue
}

type statusMessageData struct {
	version   string    `json:""`
	status    boolean   `json:""`
	timestamp time.Time `json:""`
	deviceID  string    `json:""`
}

func NewMessenger() Messenger {
	return &messenger{}
}

func (m *messenger) Connect(address string) {
	// Create a connection
	conn, err := amqp.Dial(address)
	if err != nil {
		return wrapError(err, "Failed to connect to RabbitMQ")
	}
	m.connection = conn
	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return wrapError(err, "Failed to open a channel")
	}
	m.channel = ch
	// Declare a queue
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return wrapError(err, "Failed to declare a queue")
	}
	m.statusQueue = q
}

// SendStatusUpdate sends a message informing of status change
func (m *messenger) SendStatusUpdate(status int) error {
	// Create the message data
	data := newStatusMessage(status)
	// Serialise the message
	b, err := json.Marshal(data)
	if err != nil {
		return wrapError(err, "Failed to serialise message")
	}
	// Send the message
	return sendMessage(b)
}

func (m *messenger) sendMessage(queue amqp.Queue, body []byte) error {
	// Send the message
	err = m.channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		return wrapError(err, "Failed to publish a message")
	}
	return nil
}

func newStatusMessage(status int) error {
	// Get machine ID
	id, err := machineid.ID()
	if err != nil {
		wrapError(err, "Failed to get machine ID")
	}
	//
	return statusMessageData{
		version:   "0.1.0",
		status:    status,
		timestamp: time.Now(),
		deviceID:  id,
	}
}

func wrapError(err error, msg string) {
	return fmt.Errorf("%s: %w", msg, err)
}
