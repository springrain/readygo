package test

import (
	"fmt"
	"readygo/cache"
	"testing"
)

func TestMemeryCache(t *testing.T) {

	cache.NewMemeryCacheManager()
	cache.PutToCache("cacheName", "testKey", "testValue")
	value, _ := cache.GetFromCache("cacheName", "testKey")
	a := value.(string)
	fmt.Println(a)

}
