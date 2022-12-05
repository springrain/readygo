package test

import (
	"context"
	"fmt"
	"testing"

	"readygo/cache"
)

func TestMemeryCache(t *testing.T) {
	ctx := context.Background()
	cache.NewMemeryCacheManager()

	cache.PutToCache(ctx, "cacheName", "testKey", "testValue")
	a := ""
	cache.GetFromCache(ctx, "cacheName", "testKey", &a)
	fmt.Println(a)
}

func TestRedis(t *testing.T) {
	ctx := context.Background()
	redisConfig := cache.RedisConfig{Addr: "127.0.0.1:6379"}
	cache.NewRedisClient(ctx, &redisConfig)
	cache.NewRedisCacheManager()
	cache.PutToCache(ctx, "cacheName", "testKey", "testValue")
	a := ""
	cache.GetFromCache(ctx, "cacheName", "testKey", &a)
	fmt.Println(a)
}
