package mailer

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/go-mail/mail/v2"
)

type Gomailer struct {
	dialer *mail.Dialer
	sender string
}

func NewGomailer(host string, port int, username string, password string, sender string) *Gomailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return &Gomailer{dialer: dialer, sender: sender}
}

func (g *Gomailer) Send(templatePath, username, email, activationURL string) error {
	// read template file
	t, err := template.ParseFS(FS, "templates/"+templatePath)
	if err != nil {
		return fmt.Errorf("mailer: failed to parse template file %s: %w", templatePath, err)
	}

	// execute template
	var body bytes.Buffer
	if err := t.Execute(&body, struct {
		Username      string
		ActivationURL string
	}{
		Username:      username,
		ActivationURL: activationURL,
	}); err != nil {
		return fmt.Errorf("mailer: failed to execute template file %s: %w", templatePath, err)
	}

	// create message
	msg := mail.NewMessage()
	msg.SetHeader("From", g.sender)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Finish Registration with the system")
	msg.SetBody("text/html", body.String())

	// send email
	var retryErr error
	for i := 0; i < maxRetries; i++ {
		if retryErr = g.dialer.DialAndSend(msg); retryErr != nil {
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		log.Println("mailer: email successfully sent to ", email)
		return nil
	}
	return fmt.Errorf("mailer: failed to send email to %s after %d attempts: %w", email, maxRetries, retryErr)
}
