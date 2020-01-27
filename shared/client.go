package shared

import (
	"encoding/json"
	"github.com/denisbrodbeck/machineid"
	"github.com/streadway/amqp"
	"time"
)

type Messenger interface {
	Connect(address string) error
	PowerStatusChanged(status bool) error
	Close()
}

type Receiver interface {
	Connect(address string) error
	PowerStatusConsumer() (<-chan amqp.Delivery, error)
	Close()
}

type client struct {
	connection  *amqp.Connection
	channel     *amqp.Channel
	statusQueue amqp.Queue
}

func NewMessenger() Messenger {
	return &client{}
}

func NewReceiver() Receiver {
	return &client{}
}

func (c *client) Connect(address string) error {
	// Create a connection
	conn, err := amqp.Dial(address)
	if err != nil {
		return WrapError(err, "Failed to connect to RabbitMQ")
	}
	c.connection = conn
	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return WrapError(err, "Failed to open a channel")
	}
	c.channel = ch
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
		return WrapError(err, "Failed to declare a queue")
	}
	c.statusQueue = q
	return nil
}

// PowerStatusChanged sends a message saying the power status has changed
func (c *client) PowerStatusChanged(status bool) error {
	// Create the message data
	data, err := newStatusMessage(status)
	if err != nil {
		return WrapError(err, "Failed to create message")
	}
	// Serialise the message
	b, err := json.Marshal(data)
	if err != nil {
		return WrapError(err, "Failed to serialise message")
	}
	// Send the message
	return c.sendMessage(c.statusQueue, b)
}

func (c *client) PowerStatusConsumer() (<-chan amqp.Delivery, error) {
	msgs, err := c.channel.Consume(
		c.statusQueue.Name, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		return nil, WrapError(err, "Failed to register consumer")
	}
	return msgs, nil
}

func (c *client) Close() {
	c.channel.Close()
	c.connection.Close()
}

func (c *client) sendMessage(queue amqp.Queue, body []byte) error {
	// Send the message
	err := c.channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		return WrapError(err, "Failed to publish a message")
	}
	return nil
}

func newStatusMessage(status bool) (*PowerStatusChangedV1, error) {
	// Get machine ID
	id, err := machineid.ID()
	if err != nil {
		return nil, WrapError(err, "Failed to get machine ID")
	}
	// Create the message data
	return &PowerStatusChangedV1{
		Version:   "0.1.0",
		Status:    status,
		Timestamp: time.Now(),
		DeviceID:  id,
	}, nil
}
