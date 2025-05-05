package nats

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/internal/domain"
	"log"
)

func SubscribeToOrderCreated(nc *nats.Conn) error {
	_, err := nc.Subscribe("order.created", func(msg *nats.Msg) {
		var order domain.Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("Failed to unmarshal order.created message: %v", err)
			return
		}

		log.Printf("Received order.created event:\nOrderID: %s\nUserID: %s\nTotal: %.2f", order.ID, order.UserID, order.TotalPrice)
	})

	if err != nil {
		return err
	}

	log.Println("Subscribed to 'order.created'")
	return nil
}
