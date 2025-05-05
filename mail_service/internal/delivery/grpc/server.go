package grpc

import (
	"context"
	"github.com/rakh1mbayev/Gym-Management-System/mail_service/internal/service"
	"github.com/rakh1mbayev/Gym-Management-System/mail_service/proto/mailpb"
	"log"
)

type MailServiceServer struct {
	mailpb.UnimplementedMailServiceServer
	Mailer *service.Mailer
}

func NewMailServiceServer(mailer *service.Mailer) *MailServiceServer {
	return &MailServiceServer{Mailer: mailer}
}

func (s *MailServiceServer) SendConfirmationEmail(ctx context.Context, req *mailpb.ConfirmationRequest) (*mailpb.ConfirmationResponse, error) {
	err := s.Mailer.SendEmail(req.GetEmail(), req.GetSubject(), req.GetBody())
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return &mailpb.ConfirmationResponse{
			Success: false,
			Message: "Failed to send confirmation email",
		}, nil
	}
	return &mailpb.ConfirmationResponse{
		Success: true,
		Message: "Confirmation email sent successfully",
	}, nil
}
