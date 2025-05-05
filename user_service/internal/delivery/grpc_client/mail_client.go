package grpc_client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rakh1mbayev/Gym-Management-System/mail_service/proto/mailpb"
	"google.golang.org/grpc"
)

type MailClient struct {
	client mailpb.MailServiceClient
}

func NewMailClient(cc *grpc.ClientConn) *MailClient {
	return &MailClient{
		client: mailpb.NewMailServiceClient(cc),
	}
}

func (m *MailClient) SendConfirmationEmail(email, subject, body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	req := &mailpb.ConfirmationRequest{
		Email:   email,
		Subject: subject,
		Body:    body,
	}

	resp, err := m.client.SendConfirmationEmail(ctx, req)
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("failed to send email: %s", resp.Message)
	}

	log.Printf("Email sent: %s", resp.Message)
	return nil
}
