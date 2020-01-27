package main

import (
	"github.com/briggysmalls/detectordag/shared"
	"github.com/denisbrodbeck/machineid"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Messenger interface {
	Connect(address string) error
	PowerStatusChanged(status int) error
}

type messenger struct {
	connection  amqp.Connection
	channel     amqp.Channel
	statusQueue amqp.Queue
}

func NewMessenger() Messenger {
	return &messenger{}
}

func (m *messenger) Connect(address string) {
	// Create a connection
	conn, err := amqp.Dial(address)
	if err != nil {
		return shared.wrapError(err, "Failed to connect to RabbitMQ")
	}
	m.connection = conn
	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return shared.wrapError(err, "Failed to open a channel")
	}
	m.channel = ch
	// Declare a queue
	q, err := ch.QueueDeclare(
		"power_status", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return shared.wrapError(err, "Failed to declare a queue")
	}
	m.statusQueue = q
}

// PowerStatusChanged sends a message saying the power status has changed
func (m *messenger) PowerStatusChanged(status int) error {
	// Create the message data
	data := newStatusMessage(status)
	// Serialise the message
	b, err := json.Marshal(data)
	if err != nil {
		return shared.wrapError(err, "Failed to serialise message")
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
		return shared.wrapError(err, "Failed to publish a message")
	}
	return nil
}

func newStatusMessage(status int) error {
	// Get machine ID
	id, err := machineid.ID()
	if err != nil {
		shared.wrapError(err, "Failed to get machine ID")
	}
	// Create the message data
	return shared.StatusMessageV1{
		Version:   "0.1.0",
		Status:    status,
		Timestamp: time.Now(),
		DeviceID:  id,
	}
}
