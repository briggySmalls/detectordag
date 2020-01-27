package main

import (
	"encoding/json"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"power_status", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	// Listen for messages until we're told to stop
	go func() {
		for d := range msgs {
			handleMessage(d)
		}
	}()

	// Wait for user to indicate we should quit
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func handleMessage(d amqp.Delivery) {
	// Get the message body
	body := d.Body
	log.Printf("Message received: %s", body)
	// Deserialise the JSON
	var data shared.PowerStatusChangedV1
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
}
