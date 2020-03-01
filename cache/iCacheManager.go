package cache

//iCacheManage 缓存管理器接口.缓存的结构是map[cacheName string]map[key string]valu interface{}
//内部接口,方便直接使用方法
type iCacheManage interface {
	//getFromCache 从cache中获取key值
	//缓存的结构是map[cacheName string]map[key string]valu interface{}
	getFromCache(cachName string, key string) (interface{}, error)

	//putToCache 设置指定cache中的key
	//缓存的结构是map[cacheName string]map[key string]valu interface{}
	putToCache(cacheName string, key string, value interface{}) error

	//clearCache 清理cache
	//缓存的结构是map[cacheName string]map[key string]valu interface{}
	clearCache(cacheame string) error

	//evictKey 失效一个cache中的key
	//缓存的结构是map[cacheName string]map[key string]valu interface{}
	evictKey(cacheName string, key string) error
}
