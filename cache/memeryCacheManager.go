package cache

import (
	"errors"
	"sync"
)

//memeryCacheManager 内存的缓存管理器.缓存的结构是map[cacheName string]map[key stringvalue interface{}
type memeryCacheManager struct {
	// 用于缓存反射的信息,sync.Map内部处理了并发锁.用指针地址.
	memeryCacheMap *sync.Map
}

//NewMemeryCacheManager 创建内存管理器,需要给CacheManger中的cacheManager变量赋值
func NewMemeryCacheManager() error {
	newMemeryCacheManager := memeryCacheManager{}
	//newmap := make(map[string]interface{})
	//newMemeryCacheManager.cacheMap = &newmap
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

//evictKey 失效一个cache中的key
func (cacheManager *memeryCacheManager) evictKey(cacheName string, key string) error {
	return nil
}

//获取cacheName,这个方法不在接口内,避免直接获取到cache对象
func (cacheManager *memeryCacheManager) getCache(cacheName string) (*map[string]interface{}, error) {
	if len(cacheName) < 1 {
		return nil, errors.New("cacheName为空")
	}
	cacheMapInterface, ok := cacheManager.memeryCacheMap.Load(cacheName)
	if ok { //如果cacheManager中有值
		cacheMap, mapOK := cacheMapInterface.(*map[string]interface{})
		if !mapOK { //如果类型转化失败
			return nil, errors.New("memeryCacheManager中,从memeryCacheMap取值map类型转化失败")
		}
		return cacheMap, nil
	}

	//cacheManager中没值,初始化一个map放进去,返回这个map
	newmap := make(map[string]interface{})
	cacheManager.memeryCacheMap.Store(cacheName, &newmap)
	return &newmap, nil
}
