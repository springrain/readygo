package cache

//缓存管理器接口.缓存的结构是map[cacheName string]map[key string]value interface{}
type ICacheManager interface {
	//从cache中获取key的值
	GetFromCache(cacheName string, key string) (interface{}, error)
	//设置指定cache中的key值
	PutToCache(cacheName string, key string, value interface{}) error
	//清理cache
	ClearCache(cacheName string) error
	//失效一个cache中的key
	EvictKey(cacheName string, key string) error
}
