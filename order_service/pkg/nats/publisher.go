package nats

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
)

type NatsPublisher struct {
	Conn *nats.Conn
}

func NewNatsPublisher(conn *nats.Conn) *NatsPublisher {
	return &NatsPublisher{Conn: conn}
}

func (p *NatsPublisher) PublishOrderCreated(data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return p.Conn.Publish("order.created", payload)
}
