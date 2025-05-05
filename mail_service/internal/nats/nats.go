package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

var conn *nats.Conn

func ConnectNATS() (*nats.Conn, error) {
	// Connect to NATS server (default is "nats://localhost:4222")
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %v", err)
	}
	log.Println("Connected to NATS server")

	// Return the NATS connection
	conn = nc
	return conn, nil
}

func Close() {
	if conn != nil {
		conn.Close()
		log.Println("NATS connection closed")
	}
}
