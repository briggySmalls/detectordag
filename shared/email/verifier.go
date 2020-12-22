package email

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"errors"
)

type verifier struct {
	ses          *ses.SES
}

type Verifier interface {
	VerifyEmail(email string) error
	GetVerificationStatuses(emails []string) (map[string]string, error)
	VerifyEmailsIfNecessary(emails []string) error
}

// NewVerifier gets a new Verifier
func NewVerifier(sesh *session.Session) (Verifier, error) {
	svc := ses.New(sesh)
	if svc == nil {
		return nil, errors.New("Failed to create email client")
	}
	return &verifier{
		ses: svc,
	}, nil
}

// GetVerificationStatuses gets verification status of the provided emails
// Note: If the email has never been seen before, it will be omitted from the result
func (v *verifier) GetVerificationStatuses(emails []string) (map[string]string, error) {
	// Convert to AWS type
	emailPtrs := make([]*string, len(emails))
	for i, email := range emails {
		emailPtrs[i] = aws.String(email)
	}
	// Make the request
	input := &ses.GetIdentityVerificationAttributesInput{Identities: emailPtrs}
	result, err := v.ses.GetIdentityVerificationAttributes(input)
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

func (v *verifier) VerifyEmail(email string) error {
	// Construct the input
	input := &ses.VerifyEmailIdentityInput{
		EmailAddress: aws.String(email),
	}
	// Ask to verify the email
	_, err := v.ses.VerifyEmailIdentity(input)
	return err
}

func (v *verifier) VerifyEmailsIfNecessary(emails []string) error {
	// Get the verification statuses
	statuses, err := v.GetVerificationStatuses(emails)
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
		err := v.VerifyEmail(email)
		if err != nil {
			return err
		}
	}
	return nil
}
