package app

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
