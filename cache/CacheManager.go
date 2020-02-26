package cache

//内部使用cachemanger,对外暴露使用方法
//初始化的时候赋值这个变量,只能有一个cacheManager
var cacheManager icacheManage

//GetFromCache 从cache中获取key的值.默认使用
func GetFromCache(cacheName string, key string) (interface{}, error) {
	return cacheManager.getFromCache(cacheName, key)
}

//设置指定cache中的key值
func PutToCache(cacheName string, key string, value interface{}) error {
	return cacheManager.putToCach(cacheName, key, value)
}

//清理cache
func ClearCache(cacheName string) error {
	return cacheManager.clearCache(cacheName)
}

//失效一个cache中的key
func EvictKey(cacheName string, key string) error {
	return cacheManager.evictKey(cacheName, key)
}
