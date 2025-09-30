package redis

import (
	"context"
	"encoding/json"
	"time"
	"wallet/internal/domain"

	"github.com/go-redis/redis/v8"
)

type redisCacheRepository struct {
	client *redis.Client
}

func NewRedisCacheRepository(addr string) (domain.CacheRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &redisCacheRepository{client: client}, nil
}

func (r *redisCacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// We serialize the struct to JSON before storing
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, jsonValue, ttl).Err()
}

func (r *redisCacheRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
