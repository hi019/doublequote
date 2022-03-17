package redis

import (
	"context"
	"strconv"

	"doublequote/pkg/domain"
)

var _ domain.Session = (*Session)(nil)

type Session struct {
	c      domain.CacheService
	sessId string
	uid    int
}

func (s *Session) Delete() error {
	return s.c.Delete(context.Background(), strconv.Itoa(s.uid))
}

func (s *Session) UserID() int {
	return s.uid
}
