package db

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func RedisDb() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		// DB:           0,
		// PoolSize:     100,
		// MinIdleConns: 10,
		// MaxIdleConns: 50,
		// PoolTimeout:  30 * time.Second,
		PoolSize:     100,
		MinIdleConns: 10,
		DialTimeout:  500 * time.Millisecond,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
		PoolTimeout:  1 * time.Second,
	})

	return client
}
