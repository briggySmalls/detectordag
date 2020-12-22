package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/email"
	iotp "github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
)

const (
	senderEnvVar = "SENDER_EMAIL"
)

type StatusUpdatedEvent struct {
	DeviceId  string
	Timestamp int
	State     struct {
		Status string `validate:"required,eq=on|eq=off"`
	}
	Updated struct {
		Status struct {
			Timestamp int64 `validate:"required"`
		}
	}
}

var db database.Client
var iot iotp.Client
var emailClient email.Emailer

func init() {
	// Create an AWS session
	// Good practice will share this session for all services
	sesh, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	// Create a database client
	db, err = database.New(sesh)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	// Create an IOT client
	iot, err = iotp.New(sesh)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	// Get the email sender
	sender := os.Getenv(senderEnvVar)
	if sender == "" {
		shared.LogErrorAndReturn(fmt.Errorf("Env var '%s' unset", senderEnvVar))
	}
	// Create a new session just for emailing (there is no emailing service in eu-west-2)
	sesh = shared.CreateSession(aws.Config{Region: aws.String("eu-west-1")})
	// Create a new email client
	emailClient, err = email.NewEmailer(ses.New(sesh), sender)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
}

// HandleRequest handles a lambda call
func HandleRequest(ctx context.Context, event StatusUpdatedEvent) error {
	// Print the event
	log.Printf("%v\n", event)
	// Validate the event
	if err := shared.Validate.Struct(event); err != nil {
		return err
	}
	// Get the device
	device, err := iot.GetThing(event.DeviceId)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	accountID := device.AccountId
	log.Printf("Device '%s' associated with account '%s'", event.DeviceId, accountID)
	// Get the account
	account, err := db.GetAccountById(accountID)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	// Determine parameters for the email
	stateType, transitionType, err := powerStatusToEnums(event.State.Status)
	if err != nil {
		return err
	}
	// Construct an event to pass to the emailer
	update := email.ContextData{
		DeviceName: device.Name,
		Time:       time.Unix(event.Updated.Status.Timestamp, 0),
	}
	// Send 'power status updated' emails
	log.Printf("Send emails to: %s", account.Emails)
	err = emailClient.SendUpdate(account.Emails, stateType, transitionType, update)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	return nil
}

func powerStatusToEnums(status string) (email.StateType, email.TransitionType, error) {
	// We assume we are connected if we've been given a status update
	if status == shadow.POWER_STATUS_ON {
		return email.StateTypeOn, email.TransitionTypeOn, nil
	} else if status == shadow.POWER_STATUS_OFF {
		return email.StateTypeOn, email.TransitionTypeOff, nil
	}
	return 0, 0, errors.New("Unexpected power status")
}
