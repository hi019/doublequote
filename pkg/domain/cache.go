package domain

import (
	"context"
	"time"
)

type CacheService interface {
	Set(ctx context.Context, key string, value interface{}, exp time.Duration) error
	GetInt(ctx context.Context, key string) (int, error)
	Delete(ctx context.Context, key string) error
}
