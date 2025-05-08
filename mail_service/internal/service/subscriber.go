package service

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

type Order struct {
	ID         string  `json:"id"`
	UserID     string  `json:"userId"`
	TotalPrice float64 `json:"totalPrice"`
	// You can add other fields like Items if needed
}

type NatsSubscriber struct {
	conn    *nats.Conn
	mailer  *Mailer
	subject string
}

func NewNatsSubscriber(conn *nats.Conn, mailer *Mailer) *NatsSubscriber {
	return &NatsSubscriber{
		conn:    conn,
		mailer:  mailer,
		subject: "order.created",
	}
}

func (ns *NatsSubscriber) Subscribe() error {
	_, err := ns.conn.Subscribe(ns.subject, func(m *nats.Msg) {
		var order Order
		if err := json.Unmarshal(m.Data, &order); err != nil {
			log.Printf("Failed to unmarshal order data: %v", err)
			return
		}

		// Log or use the order info
		log.Printf("Received order.created event: %+v", order)

		// Construct and send confirmation email
		subject := "Order Confirmation"
		body := fmt.Sprintf("Thank you for your order!\n\nOrder ID: %s\nTotal Price: $%.2f", order.ID, order.TotalPrice)

		err := ns.mailer.SendEmail(order.UserID, subject, body) // Assuming UserID is an email for now
		if err != nil {
			log.Printf("Failed to send order confirmation email: %v", err)
		}
	})
	return err
}
