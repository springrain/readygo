package cache

import (
	"errors"
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

	} else { //redis 集群
		redisClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:          addrs,
			RouteByLatency: true,
			Password:       redisConfig.Password, // no password set
			PoolSize:       redisConfig.PoolSize,
			MinIdleConns:   redisConfig.MinIdleConns,
		})
	}

	return nil
}
