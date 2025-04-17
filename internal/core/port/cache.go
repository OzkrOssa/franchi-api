package port

import (
	"context"
	"time"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	Get(ctx context.Context, key string) ([]byte, error)

	Delete(ctx context.Context, key string) error

	Close() error
}
