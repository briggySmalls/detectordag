package shared

import (
	"github.com/streadway/amqp"
	"log"
)

type queueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Arguments  amqp.Table
}

type clientInterface interface {
	Connect(address string) error
	Close()
}

type client struct {
	connection   *amqp.Connection
	channel      *amqp.Channel
	queues       map[string]amqp.Queue
	queueConfigs []queueConfig
}

func newClient(queueConfigs []queueConfig) *client {
	return &client{
		queueConfigs: queueConfigs,
	}
}

func (c *client) Connect(address string) error {
	// Create a connection
	log.Print("Connecting to AMQP server...")
	conn, err := amqp.Dial(address)
	if err != nil {
		return WrapError(err, "Failed to connect to RabbitMQ")
	}
	c.connection = conn
	// Create a channel
	log.Print("Creating AMQP channel...")
	ch, err := conn.Channel()
	if err != nil {
		return WrapError(err, "Failed to open a channel")
	}
	c.channel = ch
	// Declare a queue for each config
	c.queues = make(map[string]amqp.Queue, len(c.queueConfigs))
	for _, config := range c.queueConfigs {
		log.Printf("Declaring '%s' queue...", config.Name)
		q, err := ch.QueueDeclare(
			config.Name,
			config.Durable,
			config.AutoDelete,
			config.Exclusive,
			config.NoWait,
			config.Arguments,
		)
		if err != nil {
			return WrapError(err, "Failed to declare a queue")
		}
		c.queues[config.Name] = q
	}
	return nil
}

func (c *client) getConsumer(queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := c.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
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

// Send a message to the queue
func (c *client) send(queueName string, body []byte) error {
	// Get the queue out
	q := c.queues[queueName]
	// Send the message
	err := c.channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		return WrapError(err, "Failed to publish a message")
	}
	return nil
}
