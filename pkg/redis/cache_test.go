package redis

import (
	"context"
	"os"
	"testing"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewTestCache() (*CacheService, redismock.ClientMock) {
	db, mock := redismock.NewClientMock()

	cc := cache.New(&cache.Options{
		Redis: db,
	})

	c := CacheService{store: cc}

	return &c, mock
}

func TestNewCache(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		err := godotenv.Load("../.env.testing")
		require.Nil(t, err)

		NewCache(os.Getenv("REDIS_URL"))
	})
}

func TestCacheService_GetInt(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		c, mock := NewTestCache()

		// Value is the Redis key of an actual user ID, set by go-redis/cache,
		// plus \u0000 which means no compression.
		mock.ExpectGet("1").SetVal("\x05\x00\u0000")

		val, err := c.GetInt(context.Background(), "1")
		assert.Nil(t, err)
		assert.Equal(t, 5, val)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// TODO add tests for cache set.
// Mocking the Redis cache is annoying because we have to account for compression
// and such from go-redis/cache.
func TestCacheService_Set(t *testing.T) {
	t.Parallel()
}
