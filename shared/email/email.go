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
	GetVerificationStatus(emails []string) map[string]string
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

func (c *client) GetVerificationStatus(emails []string) map[string]string {
	// Convert to AWS type
	emailPtrs := make([]*string, len(emails))
	for i, email := range emails {
		emailPtrs[i] = aws.String(email)
	}
	// Construct the input
	input := &ses.GetIdentityVerificationAttributesInput{Identities: emailPtrs}
	// Make the request
	result, err := svc.GetIdentityVerificationAttributes(input)
	if err != nil {
		return nil, err
	}
	// Pull out the relevant stuff
	statuses := make(map[string]string, len(result.VerificationAttributes))
	for email, data := range result.VerificationAttributes {
		statuses[email] = data.VerificationStatus
	}
	return statuses
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

func (c *client) VerifyEmailsIfNecessary(emails []string) error {
	// Get the verification statuses
	statuses, err := c.GetVerificationStatus(emails)
	if err != nil {
		return err
	}
	// Send verification for all those that need it
	for email := range emails {
		status, ok := statuses[email]
		if !ok {
			// We have somehow not got the right status
			return errors.New("Failed to get status for email: %s", email)
		}
		if status == ses.VerificationStatusSuccess {
			// Skip emails that are already verified
			continue
		}
		// Ask to verify email
		err := c.VerifyEmail(email)
		if err != nil {
			return err
		}
	}
	return nil
}
