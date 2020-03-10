package app

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"html/template"
	"log"
	"time"
)

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "sjbriggs14@gmail.com"
	// Specify a configuration set. To use a configuration
	// set, comment the next line and line 92.
	//ConfigurationSet = "ConfigSet"
	// The subject line for the email.
	Subject = "🚨 Detectordag power update"
	// The character encoding for the email.
	CharSet = "UTF-8"
)

const textTemplateSource = `
Power status update.

Your trusty detectordag has noticed a change in your power status.
Device: {{ .DeviceId }}
Time: {{ .Timestamp.Format("15:04 01-02-06") }}
Status: {{ if .Status }}⚡️On{{else}}❗️Off{{end}}

{{ if .Status }}
We're back in business!
{{ else }}
Maybe it's time to send a neighbour round?
{{ end }}
Sent with ❤ from a dag
`

const htmlTemplateSource = `
<h1>Power status update</h1>
<p>
	Your trusty detectordag has noticed a change in your power status.
</p>
<table>
  <tr>
    <td>Device</td>
    <td>{{ .DeviceId }}</td>
  </tr>
  <tr>
    <td>Time</td>
    <td>{{ .Timestamp.Format("15:04 01-02-06") }}</td>
  </tr>
  <tr>
    <td>Status</td>
    <td>{{ if .Status }}⚡️On{{else}}❗️Off{{end}}</td>
  </tr>
</table>
{{ if .Status }}
<p>
	We're back in business!
</p>
{{ else }}
<p>
	Maybe it's time to send a neighbour round?
</p>
{{ end }}
<p>
	Sent with ❤ from a dag
</p>
`

type PowerStatusChangedEmailConfig struct {
	DeviceId  string
	Timestamp time.Time
	Status    bool
}

//It is a best practice to instantiate AWS clients outside of the AWS Lambda function handler.
//https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Streams.Lambda.BestPracticesWithDynamoDB.html
var svc *ses.SES
var htmlTemplate *template.Template
var textTemplate *template.Template

// EmailInit initialises a client for AWS SES
func EmailInit(session *session.Session) error {
	// Create SES client
	svc = ses.New(session)
	// Create templates
	var err error
	htmlTemplate, err = template.New("htmlTemplate").Parse(htmlTemplateSource)
	if err != nil {
		return err
	}
	textTemplate, err = template.New("textTemplate").Parse(textTemplateSource)
	if err != nil {
		return err
	}
	return err
}

func SendEmail(recipient string, status PowerStatusChangedEmailConfig) error {
	// Execute the templates
	var err error
	var htmlBody bytes.Buffer
	err = htmlTemplate.Execute(&htmlBody, status)
	if err != nil {
		return err
	}
	var textBody bytes.Buffer
	err = textTemplate.Execute(&textBody, status)
	if err != nil {
		return err
	}
	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
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
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}
	// Attempt to send the email.
	result, err := svc.SendEmail(input)
	if err != nil {
		return err
	}
	// Log result
	log.Printf("Message sent with ID: %s", result.MessageId)
	return nil
}
