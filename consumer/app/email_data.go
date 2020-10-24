package app

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
