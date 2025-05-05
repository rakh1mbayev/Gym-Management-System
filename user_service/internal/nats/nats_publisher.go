package nats

import (
	"github.com/nats-io/nats.go"
	"log"
)

type NatsPublisher struct {
	conn *nats.Conn
}

func NewNatsPublisher(conn *nats.Conn) *NatsPublisher {
	return &NatsPublisher{conn: conn}
}

func (p *NatsPublisher) PublishUserRegistered(email string) error {
	err := p.conn.Publish("user.registered", []byte(email))
	if err != nil {
		log.Printf("Failed to publish user.registered event: %v", err)
		return err
	}
	log.Printf("Published user.registered event for: %s", email)
	return nil
}
