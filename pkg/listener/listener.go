package listener

import (
	"context"

	"doublequote/pkg/domain"
	"go.uber.org/fx"
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
	lc fx.Lifecycle,
	eventService domain.EventService,
	emailService domain.EmailService,
	cryptoService domain.CryptoService,
	reaperService domain.ReaperService,
	ingestService domain.IngestService,
	cfg domain.Config,
) *Service {
	s := &Service{
		eventService:  eventService,
		emailService:  emailService,
		cryptoService: cryptoService,
		reaperService: reaperService,
		cfg:           cfg,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return s.start()
		},
	})

	return s
}

func (s *Service) start() error {
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

	//s.eventService.Subscribe("ingest", func(_ domain.Event) error {
	//	s.ingestService.Start()
	//	return nil
	//})
	//if err := s.eventService.PublishPeriodic("ingest", "", nil); err != nil {
	//	return err
	//}

	return nil
}
