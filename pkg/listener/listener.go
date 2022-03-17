package listener

import (
	dq "doublequote/pkg/config"
	"doublequote/pkg/domain"
)

// Ensure type implements interface.
var _ domain.ListenerService = (*Service)(nil)

type Service struct {
	eventService  domain.EventService
	emailService  domain.EmailService
	cryptoService domain.CryptoService
	cfg           dq.Config
}

func NewService(eventService domain.EventService, emailService domain.EmailService, cryptoService domain.CryptoService, cfg dq.Config) *Service {
	return &Service{
		eventService:  eventService,
		emailService:  emailService,
		cryptoService: cryptoService,
		cfg:           cfg,
	}
}

func (s *Service) Start() {
	if s.cfg.App.RequireEmailVerification {
		s.eventService.Subscribe(domain.EventTopicUserCreated, func(e domain.Event) error {
			return sendVerificationEmail(e, s.emailService, s.cryptoService, s.cfg)
		})
	}
}
