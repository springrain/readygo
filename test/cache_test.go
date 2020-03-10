package test

import (
	"context"
	"fmt"
	"readygo/cache"
	"readygo/permission/permstruct"
	"testing"
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
	cache.NewRedisClient(&redisConfig)
	cache.NewRedisCacheManager()
	user := permstruct.UserStruct{}
	user.Id = "abc"
	cache.PutToCache(ctx, "cacheName", "testKey", user)
	user.Id = ""
	err := cache.GetFromCache(ctx, "cacheName", "testKey", &user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user)

}
