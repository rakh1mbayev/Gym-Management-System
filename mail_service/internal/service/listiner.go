package service

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

type MailListener struct {
	natsConn   *nats.Conn
	mailClient *Mailer
}

func NewMailListener(nc *nats.Conn, mailClient *Mailer) *MailListener {
	return &MailListener{
		natsConn:   nc,
		mailClient: mailClient,
	}
}

func (ml *MailListener) ListenForUserRegistration() {
	_, err := ml.natsConn.Subscribe("user.registered", func(m *nats.Msg) {
		log.Println("NATS Consumer: generating mail to new user")
		userEmail := string(m.Data)
		subject := "Please confirm your email"
		body := fmt.Sprintf("Hello %s! Please confirm your email by clicking the link below.", userEmail)

		if err := ml.mailClient.SendConfirmationEmail(userEmail, subject, body); err != nil {
			log.Println("Failed to send confirmation email:", err)
		} else {
			log.Println("Confirmation email sent to:", userEmail)
		}
	})
	if err != nil {
		log.Println("Failed to subscribe to user.registered:", err)
	}
}
