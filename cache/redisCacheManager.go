package cache

//redisCacheManager redis的缓存接口实现,缓存的结构是map[cacheName string]map[key string]value interface{}
//缓存实现小写保护,避免外部直接使用实现而不使用函数,避免多个缓存实现混杂在业务中.
type redisCacheManager struct {
}

//NewRedisCacheManager 创建Redis缓存管理器,需要给CacheManager中的cacheManager变量赋值
func NewRedisCacheManager() error {

	return nil
}

//getFromCache 从cache中获取key的值
func (cacheManager *redisCacheManager) getFromCache(cacheName string, key string) (interface{}, error) {
	return nil, nil
}

//putToCache 设置指定cache中的key值
func (cacheManager *redisCacheManager) putToCache(cacheName string, key string, value interface{}) error {
	return nil
}

//clearCache 清理cache
func (cacheManager *redisCacheManager) clearCache(cacheName string) error {
	return nil
}

//evictKey 失效一个cache中的key
func (cacheManager *redisCacheManager) evictKey(cacheName string, key string) error {
	return nil
}
