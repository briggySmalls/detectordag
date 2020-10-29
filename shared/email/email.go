package email

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"html/template"
)

const (
	// The character encoding for the email.
	CharSet = "UTF-8"
)

type Verifier interface {
	VerifyEmail(email string) error
	GetVerificationStatuses(emails []string) (map[string]string, error)
	VerifyEmailsIfNecessary(emails []string) error
}

type Emailer interface {
	SendEmail(toAddresses []string, sender, subject string, context interface{}) error
}

type client struct {
	ses          *ses.SES
	htmlTemplate *template.Template
	textTemplate *template.Template
}

// NewVerifier gets a new Verifier
func NewVerifier(sesh *session.Session) (Verifier, error) {
	return createClient(sesh)
}

// NewEmailer gets a new Emailer
func NewEmailer(sesh *session.Session, htmlTemplateSource, textTemplateSource string) (Emailer, error) {
	// Create a basic client
	client, err := createClient(sesh)
	// Create templates
	htmlTemplate, err := template.New("htmlTemplate").Parse(htmlTemplateSource)
	if err != nil {
		return nil, err
	}
	client.htmlTemplate = htmlTemplate
	textTemplate, err := template.New("textTemplate").Parse(textTemplateSource)
	if err != nil {
		return nil, err
	}
	client.textTemplate = textTemplate
	// Create our client wrapper
	return client, nil
}

func (c *client) SendEmail(recipients []string, sender, subject string, context interface{}) error {
	// Execute the templates
	var err error
	var htmlBody bytes.Buffer
	err = c.htmlTemplate.Execute(&htmlBody, context)
	if err != nil {
		return err
	}
	var textBody bytes.Buffer
	err = c.textTemplate.Execute(&textBody, context)
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
	_, err = c.ses.SendEmail(input)
	if err != nil {
		return fmt.Errorf("Failed to send email: %w", err)
	}
	return nil
}

// GetVerificationStatuses gets verification status of the provided emails
// Note: If the email has never been seen before, it will be omitted from the result
func (c *client) GetVerificationStatuses(emails []string) (map[string]string, error) {
	// Convert to AWS type
	emailPtrs := make([]*string, len(emails))
	for i, email := range emails {
		emailPtrs[i] = aws.String(email)
	}
	// Make the request
	input := &ses.GetIdentityVerificationAttributesInput{Identities: emailPtrs}
	result, err := c.ses.GetIdentityVerificationAttributes(input)
	if err != nil {
		return nil, err
	}
	// Pull out the relevant stuff
	statuses := make(map[string]string, len(result.VerificationAttributes))
	for email, data := range result.VerificationAttributes {
		statuses[email] = *data.VerificationStatus
	}
	return statuses, nil
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
	statuses, err := c.GetVerificationStatuses(emails)
	if err != nil {
		return err
	}
	// Send verification for all those that need it
	for _, email := range emails {
		status, ok := statuses[email]
		// Skip emails that are already verified
		if ok && status == ses.VerificationStatusSuccess {
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

func createClient(sesh *session.Session) (*client, error) {
	// Create Amazon DynamoDB client
	svc := ses.New(sesh)
	if svc == nil {
		return nil, errors.New("Failed to create email client")
	}
	return &client{
		ses: svc,
	}, nil
}
