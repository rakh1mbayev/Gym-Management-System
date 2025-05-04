package main

import (
	"log"
	"net"
	"os"

	rpc "google.golang.org/grpc"
	"mail_service/internal/delivery/grpc"
	"mail_service/internal/service"
	"mail_service/proto/mailpb"
)

func main() {
	listener, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Load from env or config
	mailer := service.NewMailer(
		"smtp.gmail.com",
		"587",
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
		os.Getenv("SMTP_USER"),
	)

	s := grpc.NewMailServiceServer(mailer)

	grpcServer := rpc.NewServer()
	mailpb.RegisterMailServiceServer(grpcServer, s)

	log.Println("Mail service running on :8084")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
