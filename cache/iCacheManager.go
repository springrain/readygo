package cache

import "context"

//iCacheManage 缓存管理器接口.缓存的结构是map[cacheName string]map[key string]valu interface{}
//小写的属性json无法转化,struct需要实现MarshalJSON和UnmarshalJSON的接口方法
//内部接口,方便直接使用方法
type iCacheManage interface {
	//getFromCache 从cache中获取key值
	//缓存的结构是map[cacheName string]map[key string]valu interface{}
	//value形参是接收值的对象指针,例如 &user
	//小写的属性json无法转化,struct需要实现MarshalJSON和UnmarshalJSON的接口方法
	getFromCache(ctx context.Context, cachName string, key string, valuePtr interface{}) error

	//putToCache 设置指定cache中的key
	//缓存的结构是map[cacheName string]map[key string]valu interface{}
	//小写的属性json无法转化,struct需要实现MarshalJSON和UnmarshalJSON的接口方法
	putToCache(ctx context.Context, cacheName string, key string, value interface{}) error

	//clearCache 清理cache
	//缓存的结构是map[cacheName string]map[key string]valu interface{}
	clearCache(ctx context.Context, cacheame string) error

	//evictKey 失效一个cache中的key
	//缓存的结构是map[cacheName string]map[key string]valu interface{}
	evictKey(ctx context.Context, cacheName string, key string) error
}
