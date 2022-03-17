package domain

type EmailService interface {
	SendEmail(to []string, subject string, body string) error
}
