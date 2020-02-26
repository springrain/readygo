package cache

//缓存管理器接口.缓存的结构是map[cacheName string]map[key string]value interface{}

//内存的缓存管理器
type redisCacheManager struct {
}

//创建Redis缓存管理器,需要给CacheManager中的cacheManager变量赋值
func NewRedisCacheManager() error {

	return nil
}

//从cache中获取key的值
func (cacheManager *redisCacheManager) getFromCache(cacheName string, key string) (interface{}, error) {
	return nil, nil
}

//设置指定cache中的key值
func (cacheManager *redisCacheManager) putToCache(cacheName string, key string, value interface{}) error {
	return nil
}

//清理cache
func (cacheManager *redisCacheManager) clearCache(cacheName string) error {
	return nil
}

//失效一个cache中的key
func (cacheManager *redisCacheManager) evictKey(cacheName string, key string) error {
	return nil
}
