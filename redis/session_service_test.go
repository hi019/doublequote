package redis

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"doublequote/mock"
	"doublequote/utils"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
)

func NewTestSessionService() (*SessionService, *mock.CacheService) {
	cache := mock.CacheService{}
	s := NewSessionService(&cache)

	return s, &cache
}

func TestNewSessionService(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		cache := mock.CacheService{}

		_ = NewSessionService(&cache)
	})
}

func TestSessionService_Create(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		svc, store := NewTestSessionService()

		req, err := http.NewRequest("POST", "/", nil)
		rr := httptest.NewRecorder()

		store.On(
			"Set",
			req.Context(),
			tmock.IsType(""),
			1,
			utils.TimeYear,
		).Return(nil)

		sess, err := svc.Create(rr, req, 1)

		assert.Nil(t, err)
		assert.Equal(t, 1, sess.UserID())
	})

	t.Run("StoreErr", func(t *testing.T) {
		svc, cache := NewTestSessionService()

		req, err := http.NewRequest("POST", "/", nil)
		rr := httptest.NewRecorder()

		cache.On(
			"Set",
			req.Context(),
			tmock.IsType(""),
			1,
			utils.TimeYear,
		).Return(fmt.Errorf("what is life? for I am just a humble cache mock"))

		sess, err := svc.Create(rr, req, 1)

		assert.Nil(t, sess)
		assert.Equal(t, "what is life? for I am just a humble cache mock", err.Error())
		cache.AssertExpectations(t)
	})
}

func TestSessionService_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		svc, cache := NewTestSessionService()

		req, err := http.NewRequest("POST", "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session-id",
			Value: "abc",
		})

		cache.On(
			"GetInt",
			req.Context(),
			"abc",
		).Return(1, nil)

		sess, err := svc.Get(req)

		assert.Nil(t, err)
		assert.Equal(t, 1, sess.UserID())
		cache.AssertExpectations(t)
	})

	t.Run("StoreErr", func(t *testing.T) {
		svc, cache := NewTestSessionService()

		req, err := http.NewRequest("POST", "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session-id",
			Value: "abc",
		})

		cache.On(
			"GetInt",
			req.Context(),
			"abc",
		).Return(1, fmt.Errorf("error"))

		sess, err := svc.Get(req)

		assert.Nil(t, sess)
		assert.Equal(t, "error", err.Error())
		cache.AssertExpectations(t)
	})
}
