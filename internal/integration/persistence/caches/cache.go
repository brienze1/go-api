package caches

import (
	"context"
	"github.com/brienze1/notes-api/internal/infra/properties"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheClientGet interface {
	Get(ctx context.Context, key string) *redis.StringCmd
}
type CacheClientSet interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type cache struct {
	clientSet CacheClientSet
	clientGet CacheClientGet
}

func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.clientGet.Get(ctx, properties.GetCachePrefix()+key).Bytes()
}
func (c *cache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return c.clientSet.Set(ctx, properties.GetCachePrefix()+key, value, expiration).Err()
}

func NewCache(clientSet CacheClientSet, clientGet CacheClientGet) *cache {
	return &cache{clientSet: clientSet, clientGet: clientGet}
}
