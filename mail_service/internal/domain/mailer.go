package domain

type MailSender interface {
	SendEmail(to, subject, body string) error
	SendConfirmationEmail(email, subject, body string) error
}
