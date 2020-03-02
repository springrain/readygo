package cache

import (
	"encoding/json"
	"errors"
	"runtime"
	"strings"

	"github.com/go-redis/redis/v7"
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
func NewRedisClient(redisConfig *RedisConfig) error {

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
		_, err := redisClient.Ping().Result()
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
		_, err := redisClusterClient.Ping().Result()
		if err != nil {
			return err
		}
	}

	return nil
}

//hset 为redisCacheManager设置值,不再单独提供redis的API,统一为cacheManager接口
//值变成json的[]byte进行保存,小写的属性json无法转化,struct需要实现MarshalJSON和UnmarshalJSON的接口方法
func hset(hname string, key string, valuePtr interface{}) error {
	if hname == "" || key == "" || valuePtr == nil {
		return errors.New("值不能为空")
	}
	//把值转成JSON的[]byte格式
	value, errJSON := json.Marshal(valuePtr)
	if errJSON != nil {
		return errJSON
	}
	var errResult error
	if redisClient != nil { //单机redis
		_, errResult = redisClient.Do("hset", hname, key, value).Result()
	} else if redisClusterClient != nil { //集群Redis
		_, errResult = redisClusterClient.Do("hset", hname, key, value).Result()
	} else {
		return errors.New("没有redisClient或redisClusterClient实现")
	}
	//获值错误
	if errResult != nil {
		return errResult
	}
	return nil

}

//hget 获取指定的值
//取出json的[]byte进行转化,小写的属性json无法转化,struct需要实现MarshalJSON和UnmarshalJSON的接口方法
func hget(hname string, key string, valuePtr interface{}) error {
	if hname == "" || key == "" || valuePtr == nil {
		return errors.New("值不能为空")
	}
	var errResult error
	var jsonData interface{}
	if redisClient != nil { //单机redis
		jsonData, errResult = redisClient.Do("hget", hname, key).Result()
	} else if redisClusterClient != nil { //集群Redis
		jsonData, errResult = redisClusterClient.Do("hget", hname, key).Result()
	} else {
		return errors.New("没有redisClient或redisClusterClient实现")
	}
	//获值错误
	if errResult != nil {
		return errResult
	}
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

//hdel 删除一个map中的key
func hdel(hname string, key string) error {
	if hname == "" || key == "" {
		return errors.New("值不能为空")
	}
	var errResult error
	if redisClient != nil { //单机redis
		_, errResult = redisClient.HDel(hname, key).Result()
	} else if redisClusterClient != nil { //集群Redis
		_, errResult = redisClusterClient.HDel(hname, key).Result()
	} else {
		return errors.New("没有redisClient或redisClusterClient实现")
	}
	//获值错误
	if errResult != nil {
		return errResult
	}
	return nil
}

//del 删除缓存
func del(hname string) error {

	if hname == "" {
		return errors.New("值不能为空")
	}
	var errResult error
	if redisClient != nil { //单机redis
		_, errResult = redisClient.Del(hname).Result()
	} else if redisClusterClient != nil { //集群Redis
		_, errResult = redisClusterClient.Del(hname).Result()
	} else {
		return errors.New("没有redisClient或redisClusterClient实现")
	}
	//获值错误
	if errResult != nil {
		return errResult
	}
	return nil
}

//func Lock()
