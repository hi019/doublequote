package crypto

import (
	"time"

	domain "doublequote/pkg/domain"
)

// Ensure type implements interface.
var _ domain.CryptoService = (*Service)(nil)

type Service struct {
	key []byte
	now func() time.Time
}

func NewService(cfg domain.Config) *Service {
	return &Service{
		key: []byte(cfg.App.Secret),
		now: time.Now,
	}
}
