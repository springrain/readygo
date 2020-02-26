package cache

//缓存管理器接口.缓存的结构是map[cacheName string]map[key string]value interface{}
//内部接口,方便直接使用方法
type icacheManage interface {
	//从cache中获取key的值
	getFromCache(cachName string, key string) (interface{}, error)
	//设置指定cache中的key值
	putToCach(cacheName string, key string, value interface{}) error
	//清理cache
	clearCache(cacheame string) error
	//失效一个cache中的key
	evictKey(cacheName string, key string) error
}
