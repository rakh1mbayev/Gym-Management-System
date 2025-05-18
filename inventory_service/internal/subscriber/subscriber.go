package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/internal/domain"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/internal/repository/postgres"
)

type NatsSubscriber struct {
	nc   *nats.Conn
	repo postgres.ProductRepository
}

func NewNatsSubscriber(nc *nats.Conn, repo postgres.ProductRepository) *NatsSubscriber {
	return &NatsSubscriber{nc: nc, repo: repo}
}

func (s *NatsSubscriber) Subscribe() error {
	_, err := s.nc.Subscribe("order.created", func(msg *nats.Msg) {
		var event domain.OrderCreatedEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			fmt.Printf("Failed to parse order.created event: %v\n", err)
			return
		}

		for _, item := range event.Items {
			fmt.Println("NATS Consumer: receive new order", item.ProductID)
			err := s.repo.DecreaseStock(context.Background(), item.ProductID, item.Quantity)
			if err != nil {
				fmt.Printf("Failed to decrease stock for product %d: %v\n", item.ProductID, err)
			}
		}
	})
	return err
}
