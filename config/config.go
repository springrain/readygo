package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    Redis    RedisConfig    `yaml:"redis"`
	Jwt JwtConfig `yaml:"jwt"`
}

type ServerConfig struct {
    Port int `yaml:"port"`
	BasePath string `yaml:"basePath"`
}

type JwtConfig struct {
	TokenName string `yaml:"tokenName"`
	Secret string `yaml:"secret"`
	Timeout int `yaml:"timeout"`
}

type DatabaseConfig struct {
    DSN             string `yaml:"dsn"`
    DriverName      string `yaml:"driverName"`
    Dialect           string `yaml:"dialect"`
    MaxOpenConns    int    `yaml:"maxOpenConns"`
    MaxIdleConns    int    `yaml:"maxIdleConns"`
    ConnMaxLifetimeSecond int    `yaml:"connMaxLifetimeSecond"`
    SlowSQLMillis         int    `yaml:"slowSQLMillis"`
}

type RedisConfig struct {
    Addr     string `yaml:"addr"`
    Password string `yaml:"password"`
    PoolSize int    `yaml:"poolSize"`
	MinIdleConns int `yaml:"minIdleConns"`
}

// 全局变量，其他包可直接访问 config.Cfg.XXX
var Cfg Config

// init 函数：在 main 之前自动执行
func init() {
    if err := Load("config.yaml"); err != nil {
        panic("加载配置失败: " + err.Error())
    }
}

// Load 加载配置文件
func Load(filename string) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        return err
    }
    return yaml.Unmarshal(data, &Cfg)
}