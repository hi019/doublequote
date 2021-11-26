package listener

import dq "doublequote"

// Ensure type implements interface.
var _ dq.ListenerService = (*Service)(nil)

type Service struct {
	eventService  dq.EventService
	emailService  dq.EmailService
	cryptoService dq.CryptoService
	cfg           dq.Config
}

// NewService creates a new Listener service.
// The caller is responsible for setting the dependencies.
func NewService(eventService dq.EventService, emailService dq.EmailService, cryptoService dq.CryptoService, cfg dq.Config) *Service {
	return &Service{
		eventService:  eventService,
		emailService:  emailService,
		cryptoService: cryptoService,
		cfg:           cfg,
	}
}

func (s *Service) Start() {
	if s.cfg.App.RequireEmailVerification {
		s.eventService.Subscribe(dq.EventTopicUserCreated, func(e dq.Event) error {
			return sendVerificationEmail(e, s.emailService, s.cryptoService, s.cfg)
		})
	}
}
