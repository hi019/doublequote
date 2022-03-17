package redis

import (
	"context"
	"time"

	dq "doublequote/pkg/domain"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

// Ensure type implements interface.
var _ dq.CacheService = (*CacheService)(nil)

type CacheService struct {
	store *cache.Cache
	addr  string
}

// NewCache initializes a new cache.
// TODO password support
func NewCache(addr string) *CacheService {
	c := CacheService{}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	cc := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	c.store = cc

	return &c
}

func (c *CacheService) GetInt(ctx context.Context, key string) (val int, err error) {
	err = c.store.Get(ctx, key, &val)
	return val, err
}

func (c *CacheService) Set(ctx context.Context, key string, val interface{}, exp time.Duration) error {
	return c.store.Set(&cache.Item{
		Key:   key,
		Value: val,
		TTL:   exp,
		Ctx:   ctx,
	})
}

func (c *CacheService) Delete(ctx context.Context, key string) error {
	return c.store.Delete(ctx, key)
}
