package internal

import (
	"database/sql"
	"encoding/json"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/streadway/amqp"
	"log"
	"strings"
)

var QUERY string = `
SELECT email
FROM emails
INNER JOIN accounts
ON accounts.id = emails.account_id
WHERE accounts.id = (
    SELECT accounts.id
    FROM accounts
    INNER JOIN devices
    ON accounts.id = devices.account_id
    WHERE devices.id = UUID_TO_BIN(?)
)
`

func Run(address string, params DbParams) error {
	// Configure AMQP
	receiver, consumer, err := setupAMQP(address)
	if err != nil {
		return err
	}
	defer receiver.Close()

	// Configure database connection and query
	database, stmt, err := setupDB(params)
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

func setupAMQP(address string) (shared.Receiver, <-chan amqp.Delivery, error) {
	// Create a receiver
	r := shared.NewReceiver()

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

func setupDB(params DbParams) (Database, *sql.Stmt, error) {
	// Connect to the database
	database := NewDatabase()
	if err := database.Connect(params); err != nil {
		return nil, nil, err
	}

	// Prepare our query
	stmt, err := database.DB().Prepare(QUERY)
	if err != nil {
		return nil, nil, shared.WrapError(err, "Failed to prepare query")
	}

	return database, stmt, nil
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
	var emails []string
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
