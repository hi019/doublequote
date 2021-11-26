package http

import (
	"net/http"
	"testing"

	dq "doublequote"
	dqmock "doublequote/mock"
	"doublequote/utils"
	"github.com/stretchr/testify/assert"
)

func TestOpenClose(t *testing.T) {
	s := NewServer()

	err := s.Open()
	assert.Nil(t, err)

	err = s.Close()
	assert.Nil(t, err)
}

func TestServer_requireAuth(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestServer()

		req, err := http.NewRequest("", "", nil)
		assert.Nil(t, err)

		sess := dqmock.Session{}
		sess.On("UserID").Return(1)
		s.SessionService.On("Get", req).Return(&sess, nil)

		s.UserService.
			On(
				"FindUser",
				req.Context(),
				dq.UserFilter{ID: utils.IntPtr(1)},
				dq.UserInclude{},
			).
			Return(&dq.User{ID: 1}, nil)

		_, r, accepted := MakeMiddlewareRequest(req, s.requireAuth)

		assert.True(t, accepted)
		assert.Equal(t, &dq.User{ID: 1}, dq.UserFromContext(r.Context()))
		assert.Equal(t, 1, dq.UserIDFromContext(r.Context()))
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		s := NewTestServer()

		req, err := http.NewRequest("", "", nil)
		assert.Nil(t, err)

		sess := dqmock.Session{}
		sess.On("UserID").Return(1)
		s.SessionService.On("Get", req).Return(nil, nil)

		_, _, passed := MakeMiddlewareRequest(req, s.requireAuth)

		assert.False(t, passed)
	})
}
