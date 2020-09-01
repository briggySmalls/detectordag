package app

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/briggysmalls/detectordag/shared/email"
	"html/template"
	"time"
)

const (
	// Address from which emails will be sent
	Sender = "briggySmalls90@gmail.com"
	// The subject line for the email.
	Subject    = "🚨 Detectordag power update"
	DateFormat = "15:04 2/1/06"
)

const textTemplateSource = `
Power status update.

Your trusty detectordag has noticed a change in your power status.
Device: {{ .DeviceName }}
Time: {{ .Timestamp.Format "15:04 02-Jan-2006" }}
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
    <td>{{ .DeviceName }}</td>
  </tr>
  <tr>
    <td>Time</td>
    <td>{{ .Timestamp.Format "15:04 02-Jan-2006" }}</td>
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
	DeviceName string
	Timestamp  time.Time
	Status     bool
}

var emailClient email.Client
var htmlTemplate *template.Template
var textTemplate *template.Template

// init initialises a client for AWS SES
func init() {
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	var err error
	// Create a new session just for emailing (there is no emailing service in eu-west-2)
	sesh := shared.CreateSession(aws.Config{Region: aws.String("eu-west-1")})
	// Create a new email client
	emailClient, err = email.New(sesh)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	// Create templates
	htmlTemplate, err = template.New("htmlTemplate").Parse(htmlTemplateSource)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
	textTemplate, err = template.New("textTemplate").Parse(textTemplateSource)
	if err != nil {
		shared.LogErrorAndExit(err)
	}
}

func SendEmail(recipients []string, status PowerStatusChangedEmailConfig) error {
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
	// Send the email.
	err = emailClient.SendEmail(recipients, Sender, Subject, htmlBody.String(), textBody.String())
	if err != nil {
		return err
	}
	return nil
}
