package handshakejstransport_test

import (
	"../handshakejstransport"
	"testing"
)

func TestSetup(t *testing.T) {
	options := handshakejstransport.Options{SmtpAddress: "smtp.sendgrid.net", SmtpPort: "587", SmtpUsername: "username", SmtpPassword: "password"}
	handshakejstransport.Setup(options)
}
