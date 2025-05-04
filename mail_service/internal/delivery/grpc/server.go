package grpc

import (
	"context"
	"fmt"
	"mail_service/internal/service"
	"mail_service/proto/mailpb"
)

type MailServiceServer struct {
	mailpb.UnimplementedMailServiceServer
	Mailer *service.Mailer
}

func NewMailServiceServer(mailer *service.Mailer) *MailServiceServer {
	return &MailServiceServer{Mailer: mailer}
}

func (s *MailServiceServer) SendConfirmationEmail(ctx context.Context, req *mailpb.EmailRequest) (*mailpb.EmailResponse, error) {
	subject := "Confirm Your Email Address"
	confirmationLink := fmt.Sprintf("http://localhost:3000/confirm?token=%s", req.GetToken())

	body := fmt.Sprintf("Hi %s,\n\nPlease confirm your email by clicking this link:\n%s", req.GetName(), confirmationLink)

	err := s.Mailer.SendEmail(req.GetEmail(), subject, body)
	if err != nil {
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	return &mailpb.EmailResponse{Message: "Email sent successfully"}, nil
}
