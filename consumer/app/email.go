package app

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"html/template"
	"log"
)

const (
	// Address from which emails will be sent
	Sender = "briggySmalls90@gmail.com"
	// The subject line for the email.
	Subject = "🚨 Detectordag power update"
	// The character encoding for the email.
	CharSet = "UTF-8"
)

const textTemplateSource = `
Power status update.

Your trusty detectordag has noticed a change in your power status.
Device: {{ .DeviceId }}
Time: {{ .Timestamp }}
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
    <td>{{ .Timestamp }}</td>
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
	Timestamp string
	Status    bool
}

//It is a best practice to instantiate AWS clients outside of the AWS Lambda function handler.
//https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Streams.Lambda.BestPracticesWithDynamoDB.html
var svc *ses.SES
var htmlTemplate *template.Template
var textTemplate *template.Template

// init initialises a client for AWS SES
func init() {
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	var err error
	sesh, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			// There is no emailing service in eu-west-2
			Region: aws.String("eu-west-1"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	// Create SES client
	svc = ses.New(sesh)
	// Create templates
	htmlTemplate, err = template.New("htmlTemplate").Parse(htmlTemplateSource)
	if err != nil {
		log.Fatal(err)
	}
	textTemplate, err = template.New("textTemplate").Parse(textTemplateSource)
	if err != nil {
		log.Fatal(err)
	}
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
	}
	// Attempt to send the email.
	log.Printf("Sending email")
	result, err := svc.SendEmail(input)
	if err != nil {
		return err
	}
	// Log result
	log.Printf("Message sent with ID: %s", result.MessageId)
	return nil
}
