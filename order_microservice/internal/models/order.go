package models

import "time"

type Order struct {
	OrderId  int       `json:"order_id"`
	ClientId int       `json:"client_id"`
	Orders   []Orders  `json:"orders"`
	Price    float64   `json:"price"`
	Status   string    `json:"status"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type Orders struct {
	OrderId       int     `json:"order_id"`
	ProductId     int     `json:"product_id"`
	ProductName   string  `json:"product_name"`
	Quantity      int     `json:"quantity"`
	PricerPerItem float64 `json:"pricer_per_item"`
	Total         float64 `json:"total"`
}
