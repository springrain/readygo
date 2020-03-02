package test

import (
	"fmt"
	"readygo/cache"
	"testing"
)

func TestMemeryCache(t *testing.T) {

	cache.NewMemeryCacheManager()
	cache.PutToCache("cacheName", "testKey", "testValue")
	a := ""
	cache.GetFromCache("cacheName", "testKey", &a)
	fmt.Println(a)

}
