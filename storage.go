package dq

import (
	"context"
	"io"
)

type StorageService interface {
	Get(ctx context.Context, key string) (io.Reader, error)
	Set(ctx context.Context, key string, value []byte) error
}
