package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	t.Run("Size10", func(t *testing.T) {
		r := RandomString(10)
		assert.Len(t, r, 10)
	})
}
