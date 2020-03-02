package cache

import (
	"encoding/json"
	"errors"
	"sync"
)

//memeryCacheManager 内存的缓存管理器.缓存的结构是map[cacheName string]map[key stringvalue interface{}
//缓存实现小写保护,避免外部直接使用实现而不使用函数,避免多个缓存实现混杂在业务中.
type memeryCacheManager struct {
	// 用于缓存反射的信息,sync.Map内部处理了并发锁.用指针地址
	//为什么不使用指针也可以直接Load获取值啊?golang里的struct对象能直接调用指针的方法吗?
	//参照:https://blog.csdn.net/qq_31930499/article/details/93335096      https://blog.csdn.net/suiban7403/article/details/78899671
	memeryCacheMap *sync.Map
}

//NewMemeryCacheManager 创建内存管理器,需要给CacheManger中的cacheManager变量赋值
func NewMemeryCacheManager() error {
	newMemeryCacheManager := memeryCacheManager{}
	newMemeryCacheManager.memeryCacheMap = &sync.Map{}
	//赋值变量,cacheManager只能初始化一次,后面的会覆盖前面的,作为缓存实现
	cacheManager = &newMemeryCacheManager
	return nil
}

//getFromCache 从cache中获取key的
//valuePtr形参是接收值的对象指针,例如 &user
//取出json的[]byte进行转化,小写的属性json无法转化,struct需要实现MarshalJSON和UnmarshalJSON的接口方法
func (cacheManager *memeryCacheManager) getFromCache(cacheName string, key string, valuePtr interface{}) error {
	if cacheName == "" || key == "" || valuePtr == nil {
		return errors.New("值不能为空")
	}
	//获取cache
	cache, errCache := cacheManager.getCache(cacheName)
	if errCache != nil {
		return errCache
	}
	//获取cache中json的[]byte格式数据
	jsonData, _ := cache.Load(key)
	//转换成json的[]byte
	jsonBytes, jsonOK := jsonData.([]byte)
	if !jsonOK { //取值失败
		return errors.New("缓存中的格式值错误")
	}
	if len(jsonBytes) < 1 { //缓存中没有值
		return nil
	}
	//赋值
	errJSON := json.Unmarshal(jsonBytes, valuePtr)
	return errJSON
}

//putToCache 设置指定cache中的key
//值变成json的[]byte进行保存,小写的属性json无法转化,struct需要实现MarshalJSON和UnmarshalJSON的接口方法
func (cacheManager *memeryCacheManager) putToCache(cacheName string, key string, valuePtr interface{}) error {
	if cacheName == "" || key == "" || valuePtr == nil {
		return errors.New("值不能为空")
	}
	//获取cache
	cache, errCache := cacheManager.getCache(cacheName)
	if errCache != nil {
		return errCache
	}
	//key值不能为空
	if len(key) < 1 {
		return errors.New("key值不能为空")
	}

	//对象转成json的[]byte
	jsonData, errJSON := json.Marshal(valuePtr)
	if errJSON != nil {
		return errJSON
	}

	//map赋值
	cache.Store(key, jsonData)
	return nil
}

//clearCahe 清理cache
func (cacheManager *memeryCacheManager) clearCache(cacheName string) error {
	//cacheName值不能为空
	if len(cacheName) < 1 {
		return errors.New("cacheName值不能为空")
	}
	cacheManager.memeryCacheMap.Delete(cacheName)
	return nil
}

//evictKey 失效一个cache中的key
func (cacheManager *memeryCacheManager) evictKey(cacheName string, key string) error {
	if cacheName == "" || key == "" {
		return errors.New("值不能为空")
	}
	//获取cache
	cache, errCache := cacheManager.getCache(cacheName)
	if errCache != nil {
		return errCache
	}
	//key值不能为空
	if len(key) < 1 {
		return errors.New("key值不能为空")
	}
	//删除Key值
	cache.Delete(key)
	return nil
}

//获取cacheName,这个方法不在接口内,避免直接获取到cache对象
func (cacheManager *memeryCacheManager) getCache(cacheName string) (*sync.Map, error) {
	if len(cacheName) < 1 {
		return nil, errors.New("cacheName为空")
	}
	cacheMapInterface, ok := cacheManager.memeryCacheMap.Load(cacheName)
	if ok { //如果cacheManager中有值
		cacheMap, mapOK := cacheMapInterface.(*sync.Map)
		if !mapOK { //如果类型转化失败
			return nil, errors.New("memeryCacheManager中,从memeryCacheMap取值map类型转化失败")
		}
		return cacheMap, nil
	}

	//cacheManager中没值,初始化一个sync.Map放进去,返回这个map
	cache := &sync.Map{}
	cacheManager.memeryCacheMap.Store(cacheName, cache)
	return cache, nil
}
