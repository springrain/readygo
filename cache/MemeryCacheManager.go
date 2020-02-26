package cache

import "sync"

//缓存管理器接口.缓存的结构是map[cacheName string]map[key string]value interface{}
// 用于缓存反射的信息,sync.Map内部处理了并发锁.
var memeryCacheMap sync.Map

//内存的缓存管理器
type memeryCacheManager struct {
}

//创建内存管理器,需要给CacheManager中的cacheManager变量赋值
func NewMemeryCacheManager() error {
	return nil
}

//从cache中获取key的值
func (cacheManager *memeryCacheManager) getFromCache(cacheName string, key string) (interface{}, error) {
	return nil, nil
}

//设置指定cache中的key值
func (cacheManager *memeryCacheManager) putToCache(cacheName string, key string, value interface{}) error {
	return nil
}

//清理cache
func (cacheManager *memeryCacheManager) clearCache(cacheName string) error {
	return nil
}

//失效一个cache中的key
func (cacheManager *memeryCacheManager) evictKey(cacheName string, key string) error {
	return nil
}
