package cache

import "sync"

//memeryCacheMap 缓存管理器接口.缓存的结构是map[cacheName string]map[key stringvalue interface{}
// 用于缓存反射的信息,sync.Map内部处理了并发锁.
var memeryCacheMap sync.Map

//memeryCaheManager 内存的缓存管理器
type memeryCacheManager struct {
}

//NewMemeryCacheManager 创建内存管理器,需要给CacheManger中的cacheManager变量赋值
func NewMemeryCacheManager() error {
	return nil
}

//getFromCache 从cache中获取key的
func (cacheManager *memeryCacheManager) getFromCache(cacheName string, key string) (interface{}, error) {
	return nil, nil
}

//putToCache 设置指定cache中的key
func (cacheManager *memeryCacheManager) putToCache(cacheName string, key string, value interface{}) error {
	return nil
}

//clearCahe 清理cache
func (cacheManager *memeryCacheManager) clearCache(cacheName string) error {
	return nil
}

//evictKey 失效一个cache中的ke
func (cacheManager *memeryCacheManager) evictKey(cacheName string, key string) error {
	return nil
}
