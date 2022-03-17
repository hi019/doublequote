package smtp

import (
	"net/smtp"

	dq "doublequote/pkg/domain"
	"github.com/jordan-wright/email"
)

var _ dq.EmailService = (*EmailService)(nil)

type EmailService struct {
	addr string
	auth smtp.Auth
}

func NewEmailService(addr, identity, username, password, host string) *EmailService {
	e := EmailService{
		addr: addr,
	}

	e.auth = smtp.PlainAuth(identity, username, password, host)

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
