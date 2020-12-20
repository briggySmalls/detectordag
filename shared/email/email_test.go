package email

//go:generate go run github.com/golang/mock/mockgen -destination mock_ses.go -package email github.com/aws/aws-sdk-go/service/ses/sesiface SESAPI

import (
	"testing"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/golang/mock/gomock"
)

func TestPowerOffEmail(t *testing.T) {
	// Create an email client
	const (
		sender = "admin@example.com"
	)
	emailer, ses := createUnitAndMocks(t, sender)
	// Set the environment variable
	// Request an update
	to := []string{
		"john@example.com",
		"jane@example.com",
	}
	// Configure expectation to send an email
	ses.EXPECT().SendEmail(createExpectedEmailInput(
		to,
		sender,
		`<h1>Visibility status update</h1>`,
		`Visibility status update.`,
	))
	context := ContextData {
		DeviceName: "Alderney",
		Time: createTime(t, "2020/12/29 23:48:00"),
	}
	emailer.SendUpdate(to, StateTypeOff, TransitionTypeOff, context)
}

func createExpectedEmailInput(recipients []string, sender, htmlBody, textBody string) *ses.SendEmailInput {
	// Define some email parameters that never change
	charSet := "UTF-8"
	subject := "There's been a change in your dag's status"
	// Convert string array, to array of string pointers
	toAddresses := make([]*string, len(recipients))
	for i := range recipients {
		toAddresses[i] = &recipients[i]
	}
	// Construct the expected intput
	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{}, // Never CC addresses
			ToAddresses: toAddresses, // Assert the senders we've been given
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: &charSet, // Always use the same charset
					Data:    &htmlBody,
				},
				Text: &ses.Content{
					Charset: &charSet,
					Data:    &textBody,
				},
			},
			Subject: &ses.Content{
				Charset: &charSet, // Always use the same charset
				Data:    &subject, // Always use the same subject
			},
		},
		Source: &sender,
	}
}

func createUnitAndMocks(t *testing.T, sender string) (Emailer, *MockSESAPI) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock SESAPI
	ses := NewMockSESAPI(ctrl)
	// Create the new iot client
	emailer, err := NewEmailer(ses, sender)
	assert.Nil(t, err)
	return emailer, ses
}

func createTime(t *testing.T, timeString string) time.Time {
	tme, err := time.Parse("2006/01/02 15:04:05", timeString)
	assert.NoError(t, err)
	return tme
}
