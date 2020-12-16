package connection

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/email"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"log"
	"time"
)

type connectionUpdater struct {
	email  email.Emailer
	db     database.Client
	shadow shadow.Client
}

type ConnectionUpdater interface {
	UpdateConnectionStatus(device *iot.Device, timestamp time.Time, status string) error
}

func NewConnectionUpdater(sesh *session.Session, db database.Client, shadow shadow.Client) (ConnectionUpdater, error) {
	// Create a new email client
	email, err := email.NewEmailer(sesh, htmlTemplateSource, textTemplateSource)
	if err != nil {
		return nil, err
	}
	return &connectionUpdater{email: email, db: db, shadow: shadow}, nil
}

func (e *connectionUpdater) UpdateConnectionStatus(device *iot.Device, timestamp time.Time, status string) error {
	log.Printf("Sending visibility email for device: %s with state '%s'", DeviceString(device), status)
	// Update the internal record of connection status
	if err := e.shadow.UpdateConnectionStatus(device.DeviceId, status); err != nil {
		return err
	}
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
		Status:     status == shadow.CONNECTION_STATUS_CONNECTED,
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

func DeviceString(device *iot.Device) string {
	return fmt.Sprintf("Device '%s' ('%s')", device.DeviceId, device.Name)
}
