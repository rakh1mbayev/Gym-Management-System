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
	From     string
}

func NewMailer(host, port, user, pass, from string) *Mailer {
	return &Mailer{
		SMTPHost: host,
		SMTPPort: port,
		Username: user,
		Password: pass,
		From:     from,
	}
}

func (m *Mailer) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.Username, m.Password, m.SMTPHost)

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\n\r\n%s",
		to, subject, body,
	))

	return smtp.SendMail(
		m.SMTPHost+":"+m.SMTPPort,
		auth,
		m.From,
		[]string{to},
		msg,
	)
}
