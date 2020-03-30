package email

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"log"
)

const (
	// The character encoding for the email.
	CharSet = "UTF-8"
)

type Client interface {
	SendEmail(toAddresses []string, sender, subject, htmlBody, textBody string) error
	VerifyEmail(email string) error
}

type client struct {
	ses *ses.SES
}

// New gets a new Client
func New(sesh *session.Session) (Client, error) {
	// Create Amazon DynamoDB client
	svc := ses.New(sesh)
	if svc == nil {
		return nil, errors.New("Failed to create email client")
	}
	// Create our client wrapper
	client := client{
		ses: svc,
	}
	return &client, nil
}

func (c *client) SendEmail(recipients []string, sender, subject, htmlBody, textBody string) error {
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
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(textBody),
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
	log.Printf("Sending email")
	result, err := c.ses.SendEmail(input)
	if err != nil {
		return err
	}
	// Log result
	log.Printf("Message sent with ID: %s", *result.MessageId)
	return nil
}

func (c *client) VerifyEmail(email string) error {
	// Construct the input
	input := &ses.VerifyEmailIdentityInput{
		EmailAddress: aws.String(email),
	}
	// Ask to verify the email
	_, err := c.ses.VerifyEmailIdentity(input)
	return err
}
