package mock

import (
	"time"

	dq "doublequote"
	"github.com/stretchr/testify/mock"
)

// Ensure type implements interface.
var _ dq.CryptoService = (*CryptoService)(nil)

type CryptoService struct {
	mock.Mock
}

func (s *CryptoService) CreateToken(data map[string]interface{}, expires time.Duration) (string, error) {
	ret := s.Called(data, expires)

	return ret.String(0), ret.Error(1)
}

func (s *CryptoService) VerifyToken(t string) (map[string]interface{}, error) {
	ret := s.Called(t)

	return ret.Get(0).(map[string]interface{}), ret.Error(1)
}

func (s *CryptoService) HashPassword(password string) (string, error) {
	ret := s.Called(password)

	return ret.String(0), ret.Error(1)
}

func (s *CryptoService) VerifyPassword(hash, password string) bool {
	ret := s.Called(hash, password)

	return ret.Bool(0)
}
