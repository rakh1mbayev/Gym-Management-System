package grpc

import (
	"context"
	"github.com/rakh1mbayev/Gym-Management-System/mail_service/proto/mailpb"
	"log"
	_ "time"

	"google.golang.org/grpc"
)

type MailClient struct {
	client mailpb.MailServiceClient
}

func NewMailClient(cc *grpc.ClientConn) *MailClient {
	return &MailClient{client: mailpb.NewMailServiceClient(cc)}
}

func (m *MailClient) SendConfirmationEmail(ctx context.Context, email, name, token string) error {
	_, err := m.client.SendConfirmationEmail(ctx, &mailpb.EmailRequest{
		Email: email,
		Name:  name,
		Token: token,
	})
	if err != nil {
		log.Printf("Error sending confirmation email: %v", err)
	}
	return err
}
