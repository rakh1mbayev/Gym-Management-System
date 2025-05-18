package domain

type MailService interface {
	SendConfirmationEmail(to, subject, body string) error
}

type EventPublisher interface {
	PublishUserRegistered(email string) error
}
