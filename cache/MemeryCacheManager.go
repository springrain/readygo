package cache

import "sync"

// 用于缓存反射的信息,sync.Map内部处理了并发锁.
var memeryCacheMap sync.Map

//内存的缓存管理器
type MemeryCacheManager struct {
}

//创建内存管理器
func NewMemeryCacheManager() *MemeryCacheManager {

	return nil
}

//从cache中获取key的值
func (cacheManager *MemeryCacheManager) GetFromCache(cacheName string, key string) (interface{}, error) {
	return nil, nil
}

//设置指定cache中的key值
func (cacheManager *MemeryCacheManager) PutToCache(cacheName string, key string, value interface{}) error {
	return nil
}

//清理cache
func (cacheManager *MemeryCacheManager) ClearCache(cacheName string) error {
	return nil
}

//失效一个cache中的key
func (cacheManager *MemeryCacheManager) EvictKey(cacheName string, key string) error {
	return nil
}
