package database

import (
	"context"
	"fmt"
	"github.com/brienze1/notes-api/internal/integration/persistence/caches"

	"github.com/brienze1/notes-api/internal/infra/properties"
	"github.com/redis/go-redis/v9"
)

func cache(addr string) *redis.Client {
	opts := &redis.Options{
		Addr:       addr,
		MaxRetries: properties.GetCacheMaxRetries(),
		// TLSConfig: &tls.Config{
		// 	InsecureSkipVerify: false,
		// },
		DialTimeout: properties.GetCacheTimeout(),
	}

	if properties.GetEnv() == "local" {
		opts.TLSConfig = nil
	}

	client := redis.NewClient(opts)

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
	}

	return client
}

func NewCacheGet() caches.CacheClientSet {
	return cache(properties.GetCacheGetHost())
}

func NewCacheSet() caches.CacheClientGet {
	return cache(properties.GetCacheSetHost())
}
