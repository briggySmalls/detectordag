package internal

import (
	"encoding/json"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/streadway/amqp"
	"log"
)

func Run(address string) error {
	// Create a receiver
	r := shared.NewReceiver()

	// Connect
	if err := r.Connect(address); err != nil {
		return shared.WrapError(err, "Failed to create receiver")
	}
	defer r.Close()

	// Obtain the consumer
	c, err := r.PowerStatusConsumer()
	if err != nil {
		return shared.WrapError(err, "Failed to create consumer")
	}

	forever := make(chan bool)
	// Listen for messages until we're told to stop
	go func() {
		for d := range c {
			handleMessage(d)
		}
	}()

	// Wait for user to indicate we should quit
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
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
