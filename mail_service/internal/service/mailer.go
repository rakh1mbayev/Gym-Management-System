package service

import (
	"fmt"
	"log"
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

func (m *Mailer) SendConfirmationEmail(email, subject, body string) error {
	// Set up email content
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	// Send the email using SMTP
	err := smtp.SendMail(m.SMTPHost+":"+m.SMTPPort,
		smtp.PlainAuth("", m.Username, m.Password, m.SMTPHost),
		m.Username, []string{email}, msg)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Confirmation email sent to: %s", email)
	return nil
}
