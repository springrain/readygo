package cache

import (
	"context"
	"errors"
)

// 内部使用cacheManger,对外暴露使用方法
// 初始化的时候赋值这个变量,只能有一个cacheManager
// 缓存的结构是map[cacheName string]map[key string]value interface{}
var cacheManager iCacheManage = nil

// cacheManager为nil
var errNilManager error = errors.New("cacheManager为nil,请先调用NewMemeryCacheManager或者NewRedisCacheManager方法,初始化cacheManager")

// GetFromCache 从cache中获取key的值.默认使用
// 缓存的结构是map[cacheName string]map[key string]value interface{}
// valuePtr形参是接收值的对象指针,例如 &user
func GetFromCache(ctx context.Context, cacheName string, key string, valuePtr interface{}) error {
	if cacheManager == nil {
		return errNilManager
	}
	return cacheManager.getFromCache(ctx, cacheName, key, valuePtr)
}

// PutToCache 设置指定cache中的key值
// 缓存的结构是map[cacheName string]map[key string]value interface{}
func PutToCache(ctx context.Context, cacheName string, key string, value interface{}) error {
	if cacheManager == nil {
		return errNilManager
	}
	return cacheManager.putToCache(ctx, cacheName, key, value)
}

// ClearCache 清理cache
// 缓存的结构是map[cacheName string]map[key string]value interface{}
func ClearCache(ctx context.Context, cacheName string) error {
	if cacheManager == nil {
		return errNilManager
	}
	return cacheManager.clearCache(ctx, cacheName)
}

// EvictKey 失效一个cache中的key
// 缓存的结构是map[cacheName string]map[key string]value interface{}
func EvictKey(ctx context.Context, cacheName string, key string) error {
	if cacheManager == nil {
		return errNilManager
	}
	return cacheManager.evictKey(ctx, cacheName, key)
}
