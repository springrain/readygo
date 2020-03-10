package cache

import (
	"context"
	"errors"
)

//redisCacheManager redis的缓存接口实现,缓存的结构是map[cacheName string]map[key string]value interface{}
//缓存实现小写保护,避免外部直接使用实现而不使用函数,避免多个缓存实现混杂在业务中.
type redisCacheManager struct {
}

//NewRedisCacheManager 创建Redis缓存管理器,需要给CacheManager中的cacheManager变量赋值
//需要先调用RedisClient文件中的NewRedisClient(redisConfig *RedisConfig)方法,初始化RedisClient
func NewRedisCacheManager() error {
	if redisClient == nil && redisClusterClient == nil {
		return errors.New("需要先调用RedisClient文件中的NewRedisClient(redisConfig *RedisConfig)方法,初始化RedisClient")
	}

	newRedisCacheManager := redisCacheManager{}
	//赋值变量,cacheManager只能初始化一次,后面的会覆盖前面的,作为缓存实现
	cacheManager = &newRedisCacheManager
	return nil
}

//getFromCache 从cache中获取key的值
//valuePtr形参是接收值的对象指针,例如 &user
//小写的属性json无法转化,struct需要实现MarshalJSON和UnmarshalJSON的接口方法
func (cacheManager *redisCacheManager) getFromCache(ctx context.Context, cacheName string, key string, valuePtr interface{}) error {
	return redisHget(ctx, cacheName, key, valuePtr)
}

//putToCache 设置指定cache中的key值
//小写的属性json无法转化,struct需要实现MarshalJSON和UnmarshalJSON的接口方法
func (cacheManager *redisCacheManager) putToCache(ctx context.Context, cacheName string, key string, valuePtr interface{}) error {
	return redisHset(ctx, cacheName, key, valuePtr)
}

//clearCache 清理cache
func (cacheManager *redisCacheManager) clearCache(ctx context.Context, cacheName string) error {
	return redisDel(ctx, cacheName)
}

//evictKey 失效一个cache中的key
func (cacheManager *redisCacheManager) evictKey(ctx context.Context, cacheName string, key string) error {
	return redisHdel(ctx, cacheName, key)
}
