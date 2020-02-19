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
    WHERE devices.id = UUID_TO_BIN(?, true)
)
`

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
