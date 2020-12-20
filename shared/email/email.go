package email

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"html/template"
	"time"
)

const (
	// The character encoding for the email.
	CharSet = "UTF-8"
)

// StateType is an 'enum' indicating different states
type StateType int

const (
	StateTypeOn StateType = iota
	StateTypeOff StateType = iota
	StateTypeWasOn StateType = iota
	StateTypeWasOff StateType = iota
)

// TransitionType is an 'enum' on the different state transitions
type TransitionType int

const (
	TransitionTypeOn TransitionType = iota
	TransitionTypeOff TransitionType = iota
	TransitionTypeConnected TransitionType = iota
	TransitionTypeDisconnected TransitionType = iota
)

const (
	emailStatusUpdate = "There's been a change in your dag's status"
)

type emailer struct {
	ses          sesiface.SESAPI
	htmlTemplate *template.Template
	textTemplate *template.Template
	sender string
}

type Emailer interface {
	SendUpdate(toAddresses []string, state StateType, transition TransitionType, context ContextData) error
}

type ContextData struct {
	DeviceName string
	Time time.Time
}

type stateData struct {
	Graphic []byte
	Title string
	Description string
}

type transitionData struct {
	TransitionText string
}

type updateData struct {
	stateData
	transitionData
	ContextData
}

var stateDataLookup = map[StateType]stateData {
	StateTypeOn: {Title: "On", Description: "All good here!"},
	StateTypeOff: {Title: "Off", Description: "Power is gone :("},
	StateTypeWasOn: {Title: "Was On", Description: "We've lost it...but things were OK last we heard"},
	StateTypeWasOff: {Title: "Was Off", Description: "It's dead. We've lost it"},
}

var transitionDataLookup = map[TransitionType]transitionData {
	TransitionTypeOn: {TransitionText: "Power's back!"},
	TransitionTypeOff: {TransitionText: "You've lost power!"},
	TransitionTypeConnected: {TransitionText: "We've lost contact with your dag"},
	TransitionTypeDisconnected: {TransitionText: "You're dag is back"},
}

// NewEmailer gets a new Emailer
func NewEmailer(ses sesiface.SESAPI, sender string) (Emailer, error) {
	// Create templates
	htmlTemplate, err := template.New("htmlTemplate").Parse(htmlTemplateSource)
	if err != nil {
		return nil, err
	}
	textTemplate, err := template.New("textTemplate").Parse(textTemplateSource)
	if err != nil {
		return nil, err
	}
	// Create our client wrapper
	return &emailer{
		ses: ses,
		htmlTemplate: htmlTemplate,
		textTemplate: textTemplate,
		sender: sender,
	}, nil
}

func (e *emailer) SendUpdate(toAddresses []string, state StateType, transition TransitionType, context ContextData) error {
	// Get context
	c := updateData{
		ContextData: context,
		transitionData: transitionDataLookup[transition],
		stateData: stateDataLookup[state],
	}
	// Send mail
	return e.SendEmail(toAddresses, e.sender, emailStatusUpdate, c)
}

func (e *emailer) SendEmail(recipients []string, sender, subject string, context interface{}) error {
	// Execute the templates
	var err error
	var htmlBody bytes.Buffer
	err = e.htmlTemplate.Execute(&htmlBody, context)
	if err != nil {
		return err
	}
	var textBody bytes.Buffer
	err = e.textTemplate.Execute(&textBody, context)
	if err != nil {
		return err
	}
	// Convert the address into an AWS format
	toAddresses := make([]*string, len(recipients))
	for i, recipient := range recipients {
		toAddresses[i] = aws.String(recipient)
	}
	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: toAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(htmlBody.String()),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(textBody.String()),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}
	// Attempt to send the email.
	_, err = e.ses.SendEmail(input)
	if err != nil {
		return fmt.Errorf("Failed to send email: %w", err)
	}
	return nil
}
