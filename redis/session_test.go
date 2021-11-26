package redis

import (
	"context"
	"testing"

	"doublequote/mock"
	"github.com/stretchr/testify/assert"
)

func TestSession_Delete(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		cache := mock.CacheService{}

		sess := Session{
			c:      &cache,
			sessId: "abc",
			uid:    1,
		}

		cache.
			On(
				"Delete",
				context.Background(),
				"1",
			).
			Return(nil)

		err := sess.Delete()

		assert.Nil(t, err)
		cache.AssertExpectations(t)
	})
}
