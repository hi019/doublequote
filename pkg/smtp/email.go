package smtp

import (
	"net/smtp"

	"doublequote/pkg/domain"
	"github.com/jordan-wright/email"
)

var _ domain.EmailService = (*EmailService)(nil)

type EmailService struct {
	addr string
	auth smtp.Auth
}

func NewEmailService(cfg domain.Config) *EmailService {
	e := EmailService{
		addr: cfg.SMTP.URL,
	}

	e.auth = smtp.PlainAuth(cfg.SMTP.Identity, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Host)

	return &e
}

func (s *EmailService) SendEmail(to []string, subject string, body string) error {
	e := &email.Email{
		To:      to,
		From:    "Test <test@example.com>",
		Subject: subject,
		Text:    []byte(body),
	}

	return e.Send(s.addr, s.auth)
}
