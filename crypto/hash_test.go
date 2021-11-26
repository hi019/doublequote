package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_HashPassword(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewService("ABC")

		hash, err := s.HashPassword("password")
		assert.Nil(t, err)
		assert.NotNil(t, hash)

		valid := s.VerifyPassword(hash, "password")
		assert.True(t, valid)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		s := NewService("ABC")

		hash, err := s.HashPassword("password")
		assert.Nil(t, err)
		assert.NotNil(t, hash)

		valid := s.VerifyPassword(hash, "password1")
		assert.False(t, valid)
	})
}
