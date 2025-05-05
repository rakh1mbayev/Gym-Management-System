package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	rpc "github.com/rakh1mbayev/Gym-Management-System/mail_service/internal/delivery/grpc"
	"github.com/rakh1mbayev/Gym-Management-System/mail_service/internal/service"
	"github.com/rakh1mbayev/Gym-Management-System/mail_service/proto/mailpb"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load() // before os.Getenv(...)

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" {
		log.Fatal("Missing required SMTP environment variables")
	}

	mailer := service.NewMailer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	lis, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	mailpb.RegisterMailServiceServer(server, rpc.NewMailServiceServer(mailer))

	log.Println("Mail Service running on port 8084...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
