package cache

import (
	"context"
	"encoding/json"
	"errors"
	"runtime"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig redis的配置
type RedisConfig struct {
	//Addr 连接字符串,例如:127.0.0.1:6379 或者 192.168.0.2:6379,192.168.0.3:6379
	Addr string
	//密码 默认"" 没有密码
	Password string
	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int
}

//单机
var redisClient *redis.Client = nil

//集群
var redisClusterClient *redis.ClusterClient = nil

//NewRedisClient 创建Redis客户端,一个项目认为只连接一个redis即可
func NewRedisClient(ctx context.Context, redisConfig *RedisConfig) error {

	if ctx == nil {
		return errors.New("ctx不能为nil")
	}

	if redisConfig == nil {
		return errors.New("配置文件不能为nil")
	}

	if len(redisConfig.Addr) < 1 {
		return errors.New("服务器地址不能为空")
	}

	if redisConfig.PoolSize == 0 { //默认每个CPU 10个连接
		redisConfig.PoolSize = runtime.NumCPU() * 10
	}
	if redisConfig.MinIdleConns == 0 { //默认最少10个连接
		if redisConfig.PoolSize < 10 {
			redisConfig.MinIdleConns = redisConfig.PoolSize
		} else {
			redisConfig.MinIdleConns = 10
		}

	}

	//分割连接地址,判断是单机还是集群cluster
	addrs := strings.Split(redisConfig.Addr, ",")
	if len(addrs) < 1 {
		return errors.New("服务器地址不能为空")
	}
	//只有一个地址
	if len(addrs) == 1 {
		redisClient = redis.NewClient(&redis.Options{
			Addr:         addrs[0],
			Password:     redisConfig.Password, // no password set
			PoolSize:     redisConfig.PoolSize,
			MinIdleConns: redisConfig.MinIdleConns,
		})

		//验证连接有效性
		_, err := redisClient.Ping(ctx).Result()
		if err != nil {
			return err
		}

	} else { //redis 集群
		redisClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:          addrs,
			RouteByLatency: true,                 //从最近的master或者slave读取
			Password:       redisConfig.Password, // no password set
			PoolSize:       redisConfig.PoolSize,
			MinIdleConns:   redisConfig.MinIdleConns,
		})

		//验证连接有效性
		_, err := redisClusterClient.Ping(ctx).Result()
		if err != nil {
			return err
		}
	}

	return nil
}

//redisHset 为redisCacheManager设置值,不再单独提供redis的API,统一为cacheManager接口
//值变成json的[]byte进行保存,小写的属性json无法转化,struct需要实现MarshalJSON和UnmarshalJSON的接口方法
func redisHset(ctx context.Context, hname string, key string, value interface{}) error {
	if hname == "" || key == "" || value == nil {
		return errors.New("值不能为空")
	}
	//把值转成JSON的[]byte格式
	jsonData, errJSON := json.Marshal(value)
	if errJSON != nil {
		return errJSON
	}
	_, errResult := RedisCMDContext(ctx, "hset", hname, key, jsonData)
	//_, errResult := RedisCMDContext(ctx, "hset", hname, key, string(jsonData))
	//获值错误
	if errResult != nil {
		return errResult
	}
	return nil

}

//redisHget 获取指定的值
//取出json的[]byte进行转化,小写的属性json无法转化,struct需要实现MarshalJSON和UnmarshalJSON的接口方法
func redisHget(ctx context.Context, hname string, key string, valuePtr interface{}) error {
	if hname == "" || key == "" || valuePtr == nil {
		return errors.New("值不能为空")
	}

	jsonData, errResult := RedisCMDContext(ctx, "hget", hname, key)
	//获值错误
	if errResult != nil {
		return errResult
	}
	//转换成json的[]byte
	//jsonBytes, jsonOK := jsonData.([]byte)
	jsonBytes, jsonOK := jsonData.(string)
	if !jsonOK { //取值失败
		return errors.New("缓存中的格式值错误")
	}
	if len(jsonBytes) < 1 { //缓存中没有值
		return nil
	}
	//赋值
	//errJSON := json.Unmarshal(jsonBytes, valuePtr)
	errJSON := json.Unmarshal([]byte(jsonBytes), valuePtr)
	return errJSON
}

//redisHdel 删除一个map中的key
func redisHdel(ctx context.Context, hname string, key string) error {
	if hname == "" || key == "" {
		return errors.New("值不能为空")
	}
	_, errResult := RedisCMDContext(ctx, "hdel", hname, key)

	//获值错误
	if errResult != nil {
		return errResult
	}
	return nil
}

//redisDel 删除缓存
func redisDel(ctx context.Context, cacheName string) error {

	if cacheName == "" {
		return errors.New("值不能为空")
	}
	_, errResult := RedisCMDContext(ctx, "del", cacheName)
	//获值错误
	if errResult != nil {
		return errResult
	}
	return nil
}

func redisGet(ctx context.Context, cacheName string) (interface{}, error) {
	if cacheName == "" {
		return nil, errors.New("值不能为空")
	}
	result, errResult := RedisCMDContext(ctx, "get", cacheName)

	//获值错误
	if errResult != nil {
		return nil, errResult
	}
	return result, nil
}

//RedisINCR redis实现的计数器
func RedisINCR(ctx context.Context, cacheName string) (interface{}, error) {
	if cacheName == "" {
		return nil, errors.New("值不能为空")
	}
	result, errResult := RedisCMDContext(ctx, "incr", cacheName)
	//获值错误
	if errResult != nil {
		return nil, errResult
	}
	return result, nil
}

//RedisLock redis实现的分布式锁
//参数:lockName锁的名称,timeoutSecond超时秒数默认5秒,分布式锁内业务的匿名函数
//返回值:true获取锁成功,获取锁失败false,匿名函数返回值,错误信息
func RedisLock(ctx context.Context, lockName string, timeoutSecond int, doLock func() (interface{}, error)) (bool, interface{}, error) {
	if lockName == "" {
		return false, nil, errors.New("lockName值不能为空")
	}

	if timeoutSecond == 0 { //如果没有超时时间,默认5秒
		timeoutSecond = 5
	}

	//获取超时的时间,作为value
	value := time.Now().Unix() + int64(timeoutSecond)

	lockedStatus, errResult := RedisCMDContext(ctx, "set", lockName, value, "ex", timeoutSecond, "nx")
	locked, lockOK := lockedStatus.(int)
	if !lockOK { //结果异常
		return false, nil, errors.New("获取锁状态异常")
	}

	//获值错误或者没有获取到锁
	if errResult != nil || (locked == 0) {
		return false, nil, errResult
	}
	//确保解锁逻辑执行
	defer func() {
		//当前时间
		newValue := time.Now().Unix()
		if newValue >= value { //已经过了超时时间,不需要解锁
			return
		}

		lockValue, errLock := redisGet(ctx, lockName)
		if errLock != nil { //从redis获取值出现异常,解锁
			redisDel(ctx, lockName)
			return
		}

		oldValue, newOK := lockValue.(int64)
		if !newOK { //如果获取值异常,解锁
			redisDel(ctx, lockName)
			return
		}
		if oldValue != value { //值已经被其他的程序修改,已经不再是本程序的锁了,返回
			return
		}

		//其他情况,解锁
		redisDel(ctx, lockName)

	}()

	//调用业务逻辑
	result, errLock := doLock()

	//返回业务逻辑
	return locked == 1, result, errLock
}

//RedisCMDContext 运行redis指令
func RedisCMDContext(ctx context.Context, args ...interface{}) (interface{}, error) {
	var result interface{}
	var errResult error
	if redisClient != nil { //单机redis
		result, errResult = redisClient.Do(ctx, args...).Result()
	} else if redisClusterClient != nil { //集群Redis
		result, errResult = redisClusterClient.Do(ctx, args...).Result()
	} else {
		return nil, errors.New("没有redisClient或redisClusterClient实现")
	}
	//获值错误
	if errResult != nil {
		return nil, errResult
	}
	return result, nil
}
