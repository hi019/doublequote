package redis

import (
	"net/http"
	"time"

	dq "doublequote"
	"doublequote/utils"
	"github.com/go-redis/cache/v8"
)

const cookieName = "session-id"

// Ensure type implements interface.
var _ dq.SessionService = (*SessionService)(nil)

type SessionService struct {
	c dq.CacheService
}

func NewSessionService(c dq.CacheService) *SessionService {
	s := &SessionService{}
	s.c = c

	return s
}

func (s *SessionService) Get(r *http.Request) (dq.Session, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, nil
	}

	uid, err := s.c.GetInt(r.Context(), cookie.Value)
	switch err {
	case cache.ErrCacheMiss:
		return nil, nil
	case nil:
		break
	default:
		return nil, err
	}

	sess := &Session{
		c:      s.c,
		sessId: cookie.Value,
		uid:    uid,
	}

	return sess, nil
}

func (s *SessionService) Create(w http.ResponseWriter, r *http.Request, uid int) (dq.Session, error) {
	sessId := utils.RandomString(10)

	err := s.c.Set(r.Context(), sessId, uid, utils.TimeYear)
	if err != nil {
		return nil, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    cookieName,
		Value:   sessId,
		Path:    "/api",
		Expires: time.Now().Add((time.Hour * 24) * 7),
	})

	sess := Session{
		c:      s.c,
		sessId: sessId,
		uid:    uid,
	}

	return &sess, err
}
