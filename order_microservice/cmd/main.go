package cmd

import (
	"Gym-Management-System/order_microservice/internal/repository"
	"log"
)

func main() {
	db, err := repository.ConnectionDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
