package main

import (
	"log"
	"net"
	"os"

	rpc "github.com/rakh1mbayev/Gym-Management-System/mail_service/internal/delivery/grpc"
	"github.com/rakh1mbayev/Gym-Management-System/mail_service/internal/service"
	"github.com/rakh1mbayev/Gym-Management-System/mail_service/proto/mailpb"
	"google.golang.org/grpc"
)

func main() {
	mailer := service.NewMailer(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)

	lis, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := rpc.NewMailServiceServer(mailer)
	server := grpc.NewServer()
	mailpb.RegisterMailServiceServer(server, grpcServer)

	log.Println("Mail Service running on port 8084...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
