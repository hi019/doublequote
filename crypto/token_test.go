package crypto

import (
	"testing"
	"time"

	dq "doublequote"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateToken(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewService("ABC")

		data := map[string]interface{}{"email": "test@example.com"}
		token, err := s.CreateToken(data, time.Hour)
		assert.Nil(t, err)
		assert.NotNil(t, token)

		retrieved, err := s.VerifyToken(token)
		assert.Nil(t, err)
		assert.Equal(t, "test@example.com", retrieved["email"])
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		s := NewService("ABC")

		data := map[string]interface{}{"email": "test@example.com"}
		token, err := s.CreateToken(data, time.Nanosecond)
		assert.Nil(t, err)
		assert.NotNil(t, token)

		s.now = func() time.Time {
			return time.Now().Add(time.Hour)
		}

		retrieved, err := s.VerifyToken(token)
		assert.Equal(t, "Token expired.", dq.ErrorMessage(err))
		assert.Equal(t, dq.EINVALID, dq.ErrorCode(err))
		assert.Nil(t, retrieved)
	})
}

func TestService_VerifyToken(t *testing.T) {
	t.Run("InvalidToken", func(t *testing.T) {
		s := NewService("ABC")

		retrieved, err := s.VerifyToken("NOTATOKEN")
		assert.Equal(t, "Internal error.", dq.ErrorMessage(err))
		assert.Equal(t, dq.EINTERNAL, dq.ErrorCode(err))
		assert.Nil(t, retrieved)
	})
}
