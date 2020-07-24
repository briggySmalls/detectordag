package main

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/briggysmalls/detectordag/shared/email"
	"html/template"
	"log"
	"time"
)

const (
	// Address from which emails will be sent
	Sender     = "briggySmalls90@gmail.com"
	DateFormat = "15:04 2/1/06"
)

const textTemplateSource = `
Visibility status update.

Whilst keeping an eye on your dag, we noticed something change.
Device: {{ .DeviceName }}
Time: {{ .Timestamp.Format "15:04 02-Jan-2006" }}
Status: {{ if .Status }}👋Found{{else}}💨Lost{{end}}

{{ if .Status }}
We're back in business!
{{ else }}
Maybe it's time to send a neighbour round?
{{ end }}
Sent with ❤ from a dag
`

const htmlTemplateSource = `
<h1>Visibility status update</h1>
<p>
	Whilst keeping an eye on your dag, we noticed something change.
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
    <td>{{ if .Status }}👋Found{{else}}💨Lost{{end}}</td>
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

type VisibilityStatusChangedEmailConfig struct {
	DeviceName string
	Timestamp  time.Time
	Status     bool
}

var emailClient email.Client
var htmlTemplate *template.Template
var textTemplate *template.Template

// init initialises a client for AWS SES
func init() {
	var err error
	// Create a new session just for emailing (there is no emailing service in eu-west-2)
	emailSesh := shared.CreateSession(aws.Config{Region: aws.String("eu-west-1")})
	// Create a new email client
	emailClient, err = email.New(emailSesh)
	if err != nil {
		log.Fatal(err.Error())
	}
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

func SendEmail(recipients []string, status VisibilityStatusChangedEmailConfig) error {
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
	var subject string
	if status.Status {
		subject = "👋 We've found your dag again!"
	} else {
		subject = "💨 You're dag's gone missing!"
	}
	// Send the email.
	err = emailClient.SendEmail(recipients, Sender, subject, htmlBody.String(), textBody.String())
	if err != nil {
		return err
	}
	return nil
}
