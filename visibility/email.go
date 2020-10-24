package visibility

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/email"
	"github.com/briggysmalls/detectordag/shared/iot"

	"log"
	"time"
)

type emailClient struct {
	email email.Client
	db    database.Client
}

type EmailClient interface {
	SendVisiblityStatus(device *iot.Device, timestamp time.Time, status bool) error
}

func New(sesh *session.Session, db database.Client) (EmailClient, error) {
	// Create a new email client
	email, err := email.New(sesh, htmlTemplateSource, textTemplateSource)
	if err != nil {
		return nil, err
	}
	return &emailClient{email: email, db: db}, nil
}

func (e *emailClient) SendVisiblityStatus(device *iot.Device, timestamp time.Time, status bool) error {
	log.Printf("Sending visibility email for device: %s with state %v", DeviceString(device), status)
	// Get the account
	account, err := e.db.GetAccountById(device.AccountId)
	if err != nil {
		return err
	}
	// Assemble the visibility status context
	context := struct {
		DeviceName string
		Timestamp  time.Time
		Status     bool
	}{
		DeviceName: device.Name,
		Timestamp:  timestamp,
		Status:     status,
	}
	// Determine the subject
	var subject string
	if context.Status {
		subject = "👋 We've found your dag again!"
	} else {
		subject = "💨 You're dag's gone missing!"
	}
	// Send the email.
	return e.email.SendEmail(account.Emails, Sender, subject, context)
}
