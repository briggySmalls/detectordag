package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
)

const (
	// The character encoding for the email.
	CharSet = "UTF-8"
)

// StateType is an 'enum' indicating different states
type StateType int

const (
	StateTypeOn     StateType = iota
	StateTypeOff    StateType = iota
	StateTypeWasOn  StateType = iota
	StateTypeWasOff StateType = iota
)

// TransitionType is an 'enum' on the different state transitions
type TransitionType int

const (
	TransitionTypeOn           TransitionType = iota
	TransitionTypeOff          TransitionType = iota
	TransitionTypeConnected    TransitionType = iota
	TransitionTypeDisconnected TransitionType = iota
)

const (
	emailStatusUpdate = "There's been a change in your dag's status"
)

type emailer struct {
	ses          sesiface.SESAPI
	htmlTemplate *template.Template
	textTemplate *template.Template
	sender       string
}

type Emailer interface {
	SendUpdate(toAddresses []string, state StateType, transition TransitionType, context ContextData) error
}

type ContextData struct {
	DeviceName string
	Time       time.Time
}

type stateData struct {
	ImageSrc    string
	Title       string
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

var stateDataLookup = map[StateType]stateData{
	StateTypeOn: {
		Title:       "On",
		Description: "The power is on!",
		ImageSrc:    "https://detectordag.tk/android-chrome-192x192.png"
	},
	StateTypeOff: {Title: "Off",
		Description: "Your dag says that the power is off",
		ImageSrc:    "https://detectordag.tk/android-chrome-192x192.png"
	},
	StateTypeWasOn: {
		Title:       "Was On",
		Description: "We've lost contact with your dag. The power was on the last we heard...",
		ImageSrc:    "https://detectordag.tk/android-chrome-192x192.png"
	},
	StateTypeWasOff: {
		Title:       "Was Off",
		Description: "Your dag noticed the power go, and then we lost contact. It may have run out of battery.",
		ImageSrc:    "https://detectordag.tk/android-chrome-192x192.png"
	},
}

var transitionDataLookup = map[TransitionType]transitionData{
	TransitionTypeOn:           {TransitionText: "Your power's back!"},
	TransitionTypeOff:          {TransitionText: "You've lost power!"},
	TransitionTypeConnected:    {TransitionText: "We've lost contact with your dag"},
	TransitionTypeDisconnected: {TransitionText: "You're dag is back"},
}

// NewEmailer gets a new Emailer
func NewEmailer(ses sesiface.SESAPI, sender string) (Emailer, error) {
	// Create templates
	htmlTemplate, err := template.New("htmlTemplate").Parse(updateHtmlTemplateSource)
	if err != nil {
		return nil, err
	}
	textTemplate, err := template.New("textTemplate").Parse(textTemplateSource)
	if err != nil {
		return nil, err
	}
	// Create our client wrapper
	return &emailer{
		ses:          ses,
		htmlTemplate: htmlTemplate,
		textTemplate: textTemplate,
		sender:       sender,
	}, nil
}

func (e *emailer) SendUpdate(toAddresses []string, state StateType, transition TransitionType, context ContextData) error {
	// Get context
	c := updateData{
		ContextData:    context,
		transitionData: transitionDataLookup[transition],
		stateData:      stateDataLookup[state],
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
	log.Print(htmlBody)
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
