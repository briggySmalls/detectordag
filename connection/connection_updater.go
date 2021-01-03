package connection

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/email"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
)

type connectionUpdater struct {
	email  email.Emailer
	db     database.Client
	shadow shadow.Client
}

type ConnectionUpdater interface {
	UpdateConnectionStatus(device *iot.Device, timestamp time.Time, status string) error
}

func NewConnectionUpdater(sesh *session.Session, db database.Client, shadow shadow.Client, sender string) (ConnectionUpdater, error) {
	// Create a new email client
	email, err := email.NewEmailer(ses.New(sesh), sender)
	if err != nil {
		return nil, err
	}
	return &connectionUpdater{email: email, db: db, shadow: shadow}, nil
}

func (e *connectionUpdater) UpdateConnectionStatus(device *iot.Device, timestamp time.Time, status string) error {
	log.Printf("Sending visibility email for device: %s with state '%s'", DeviceString(device), status)
	// Update the internal record of connection status
	shdw, err := e.shadow.UpdateConnectionStatus(device.DeviceId, status)
	if err != nil {
		return err
	}
	// Get the account
	account, err := e.db.GetAccountById(device.AccountId)
	if err != nil {
		return err
	}
	// Assemble the visibility status context
	context := email.ContextData{
		DeviceName: shdw.Name,
		Time:       timestamp,
	}
	state, transition, err := connectionStatusToEnums(shdw)
	if err != nil {
		return err
	}
	// Send the email.
	return e.email.SendUpdate(account.Emails, state, transition, context)
}

func DeviceString(device *iot.Device) string {
	return fmt.Sprintf("Device '%s'", device.DeviceId)
}

func connectionStatusToEnums(shdw *shadow.Shadow) (email.StateType, email.TransitionType, error) {
	connection := shdw.Connection.Value
	power := shdw.Power.Value
	// Lookup the state
	state, err := email.ToStateType(connection, power)
	if err != nil {
		return 0, 0, err
	}
	// Set the transition type
	switch connection {
	case shadow.CONNECTION_STATUS_CONNECTED:
		return state, email.TransitionTypeConnected, nil
	case shadow.CONNECTION_STATUS_DISCONNECTED:
		return state, email.TransitionTypeDisconnected, nil
	default:
		return 0, 0, errors.New("Unexpected power status")
	}
}
