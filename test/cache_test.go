package test

import (
	"context"
	"fmt"
	"readygo/cache"
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
