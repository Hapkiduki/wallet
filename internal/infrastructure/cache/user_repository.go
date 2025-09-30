package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"wallet/internal/domain"
)

type cachedUserRepository struct {
	cacheRepo domain.CacheRepository
	nextRepo  domain.UserRepository // The "next" repository in the chain (Postgres)
}

func NewCachedUserRepository(cache domain.CacheRepository, next domain.UserRepository) domain.UserRepository {
	return &cachedUserRepository{
		cacheRepo: cache,
		nextRepo:  next,
	}
}

func (c *cachedUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	// 1. First, try to get the user from the cache.
	cacheKey := fmt.Sprintf("user:%s", id)
	cachedUserJSON, err := c.cacheRepo.Get(ctx, cacheKey)

	// 2. Cache Hit: If found, deserialize and return it.
	if err == nil && cachedUserJSON != "" {
		var user domain.User
		if err := json.Unmarshal([]byte(cachedUserJSON), &user); err == nil {
			return &user, nil
		}
	}

	// 3. Cache Miss: If not in cache, get it from the database.
	user, err := c.nextRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 4. Store the result in the cache for next time.
	c.cacheRepo.Set(ctx, cacheKey, user, 5*time.Minute) // Cache for 5 minutes

	return user, nil
}

// For methods that change data, we just pass them through and could optionally invalidate the cache.
func (c *cachedUserRepository) Save(ctx context.Context, user *domain.User) error {
	return c.nextRepo.Save(ctx, user)
}

func (c *cachedUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	// This could also be cached, but we'll leave it for simplicity.
	return c.nextRepo.FindByUsername(ctx, username)
}
