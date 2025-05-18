package service_test

import (
	"github.com/rakh1mbayev/Gym-Management-System/mail_service/internal/service"
	"testing"
)

func TestMailer_SendEmail_Integration(t *testing.T) {
	mailer := service.NewMailer(
		"localhost",
		"1025",
		"test@example.com",
		"password",
	)

	err := mailer.SendEmail(
		"receiver@example.com",
		"Test Subject",
		"This is a test email from integration test",
	)

	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}
}
