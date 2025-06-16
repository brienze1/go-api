package mocks

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"sync"
)

var redisConnOnce sync.Once
var redisConn *redis.Client

func Redis() *redis.Client {
	if redisConn == nil {
		redisConnOnce.Do(
			func() {
				redisConn = openRedisConn()
			},
		)
	}

	return redisConn
}

func openRedisConn() *redis.Client {
	miniRedis, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	conn := redis.NewClient(
		&redis.Options{
			Addr: miniRedis.Addr(),
		},
	)

	return conn
}

func ClearRedis(redis *redis.Client) error {
	return redis.FlushAll(context.TODO()).Err()
}
