package email

const textTemplateSource = `
There's been an update for your dag {{ .DeviceName }}
At {{ .ContextData.Time.Format "15:04 02-Jan-2006" }}
{{ .TransitionText }}
{{ .Title }}
{{ .Description }}`

//go:generate go run mail/include_data.go mail/

const htmlTemplateSource = `
<h1>There's been an update for your dag {{ .DeviceName }}</h1>
<p>{{ .TransitionText }}</p>
<table>
  <tr>
    <td>Time</td>
    <td>{{ .Time.Format "15:04 02-Jan-2006" }}</td>
  </tr>
  <tr>
    <td>Title</td>
    <td>{{ .Title }}</td>
  </tr>
  <tr>
    <td>Description</td>
    <td>{{ .Description }}</td>
  </tr>
</table>
`
