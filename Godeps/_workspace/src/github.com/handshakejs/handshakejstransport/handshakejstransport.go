package handshakejstransport

import (
	"github.com/handshakejs/handshakejserrors"
	"github.com/jordan-wright/email"
	"net/smtp"
)

var (
	SMTP_ADDRESS  string
	SMTP_PORT     string
	SMTP_USERNAME string
	SMTP_PASSWORD string
)

type Options struct {
	SmtpAddress  string
	SmtpPort     string
	SmtpUsername string
	SmtpPassword string
}

func Setup(options Options) {
	if options.SmtpAddress == "" {
		SMTP_ADDRESS = ""
	} else {
		SMTP_ADDRESS = options.SmtpAddress
	}
	if options.SmtpPort == "" {
		SMTP_PORT = ""
	} else {
		SMTP_PORT = options.SmtpPort
	}
	if options.SmtpUsername == "" {
		SMTP_USERNAME = ""
	} else {
		SMTP_USERNAME = options.SmtpUsername
	}
	if options.SmtpPassword == "" {
		SMTP_PASSWORD = ""
	} else {
		SMTP_PASSWORD = options.SmtpPassword
	}
}

func ViaEmail(to, from, subject, text, html string) *handshakejserrors.LogicError {
	e := email.NewEmail()
	e.From = from
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(text)
	e.HTML = []byte(html)

	err := e.Send(SMTP_ADDRESS+":"+SMTP_PORT, smtp.PlainAuth("", SMTP_USERNAME, SMTP_PASSWORD, SMTP_ADDRESS))
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return logic_error
	}

	return nil
}
