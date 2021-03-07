package email

const textTemplateSource = `
There's been an update for your dag {{ .DeviceName }}
At {{ .ContextData.Time.Format "15:04 02-Jan-2006" }}
{{ .TransitionText }}
{{ .Title }}
{{ .Description }}`
