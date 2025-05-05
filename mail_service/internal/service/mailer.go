package service

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	SMTPHost string
	SMTPPort string
	Username string
	Password string
}

func NewMailer(host, port, username, password string) *Mailer {
	return &Mailer{
		SMTPHost: host,
		SMTPPort: port,
		Username: username,
		Password: password,
	}
}

func (m *Mailer) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.Username, m.Password, m.SMTPHost)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%s", m.SMTPHost, m.SMTPPort)
	return smtp.SendMail(addr, auth, m.Username, []string{to}, msg)
}
