package cache

import "sync"

//memeryCacheManager 内存的缓存管理器.缓存的结构是map[cacheName string]map[key stringvalue interface{}
type memeryCacheManager struct {
	// 用于缓存反射的信息,sync.Map内部处理了并发锁.用指针地址.
	memeryCacheMap *sync.Map
	//cache的Map,用于实际存储map缓存.用指针地址.
	cacheMap *map[string]interface{}
}

//NewMemeryCacheManager 创建内存管理器,需要给CacheManger中的cacheManager变量赋值
func NewMemeryCacheManager() error {
	newMemeryCacheManager := memeryCacheManager{}
	newmap := make(map[string]interface{})
	newMemeryCacheManager.cacheMap = &newmap
	//赋值变量,cacheManager只能初始化一次,后面的会覆盖前面的,作为缓存实现
	cacheManager = &newMemeryCacheManager
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
