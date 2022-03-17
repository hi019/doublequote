package crypto

import (
	"time"

	dq "doublequote/pkg/domain"
)

// Ensure type implements interface.
var _ dq.CryptoService = (*Service)(nil)

type Service struct {
	key []byte
	now func() time.Time
}

func NewService(key string) *Service {
	return &Service{
		key: []byte(key),
		now: time.Now,
	}
}
