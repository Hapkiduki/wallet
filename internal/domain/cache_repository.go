package domain

import (
	"context"
	"time"
)

// CacheRepository defines the contract for a cache.
type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
