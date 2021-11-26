package redis

import (
	"context"
	"strconv"

	dq "doublequote"
)

var _ dq.Session = (*Session)(nil)

type Session struct {
	c      dq.CacheService
	sessId string
	uid    int
}

func (s *Session) Delete() error {
	return s.c.Delete(context.Background(), strconv.Itoa(s.uid))
}

func (s *Session) UserID() int {
	return s.uid
}
