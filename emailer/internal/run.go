package internal

import (
	"encoding/json"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/streadway/amqp"
	"log"
	"strings"
	"net/smtp"
)

func Run(address string, params EmailParams) error {
	// Configure AMQP
	receiver, consumer, err := setupAMQP(address)
	if err != nil {
		return err
	}
	defer receiver.Close()

	// Configure email connection
	database, stmt, err := setupEmail(params)
	if err != nil {
		return err
	}
	defer database.Close()
	defer stmt.Close()

	forever := make(chan bool)
	// Listen for messages until we're told to stop
	go func() {
		for delivery := range consumer {
			handleMessage(delivery, stmt)
		}
	}()

	// Wait for user to indicate we should quit
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

func setupAMQP(address string) (shared.SensingReceiver, <-chan amqp.Delivery, error) {
	// Create a receiver
	r := shared.NewSensingReceiver()

	// Connect
	if err := r.Connect(address); err != nil {
		return nil, nil, err
	}

	// Obtain the consumer
	c, err := r.PowerStatusConsumer()
	if err != nil {
		return nil, nil, err
	}

	return r, c, nil
}

func setupEmail(params EmailParams) (, error) {

}

func handleMessage(m amqp.Delivery, stmt *sql.Stmt) {
	// Get the message body
	body := m.Body
	log.Printf("Message received: %s", body)
	// Deserialise the JSON
	var data shared.PowerStatusChangedV1
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	// Query the emails associated with the device
	rows, err := stmt.Query(data.DeviceID)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()
	// Iterate through the emails
	emails := make([]string, 0)
	for rows.Next() {
		// Obtain the email
		var email string
		err := rows.Scan(&email)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		emails = append(emails, email)
	}
	// Print the emails
	log.Printf("Emails to send: %s", strings.Join(emails, ", "))
}
