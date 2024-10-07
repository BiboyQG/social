package mailer

import "embed"

const (
	maxRetries   = 3
	TemplatePath = "user_invitation.html"
)

//go:embed "templates"
var FS embed.FS

type Mailer interface {
	Send(templatePath, username, email, activationURL string) error
}
