package listener

import (
	"doublequote/pkg/domain"
)

// Ensure type implements interface.
var _ domain.ListenerService = (*Service)(nil)

type Service struct {
	eventService  domain.EventService
	emailService  domain.EmailService
	cryptoService domain.CryptoService
	reaperService domain.ReaperService
	ingestService domain.IngestService

	cfg domain.Config
}

func NewService(
	eventService domain.EventService,
	emailService domain.EmailService,
	cryptoService domain.CryptoService,
	reaperService domain.ReaperService,
	ingestService domain.IngestService,
	cfg domain.Config,
) *Service {
	return &Service{
		eventService:  eventService,
		emailService:  emailService,
		cryptoService: cryptoService,
		reaperService: reaperService,
		cfg:           cfg,
	}
}

func (s *Service) Start() error {
	if s.cfg.App.RequireEmailVerification {
		s.eventService.Subscribe(domain.EventTopicUserCreated, func(e domain.Event) error {
			return sendVerificationEmail(e, s.emailService, s.cryptoService, s.cfg)
		})
	}

	s.eventService.Subscribe("reaper", func(_ domain.Event) error {
		return s.reaperService.Run()
	})
	if err := s.eventService.PublishPeriodic("reaper", "", nil); err != nil {
		return err
	}

	s.eventService.Subscribe("ingest", func(_ domain.Event) error {
		s.ingestService.Start()
		return nil
	})
	if err := s.eventService.PublishPeriodic("ingest", "", nil); err != nil {
		return err
	}

	return nil
}
