package util

import (
	"github.com/vmihailenco/msgpack/v5"
)

// 用于隔离序列化工具

// init 初始化msgpack,设置使用json标签
func init() {
	// 创建一个编码器并设置使用 json 标签
	enc := msgpack.NewEncoder(nil)
	enc.SetCustomStructTag("json") // 关键配置：指定使用 json 标签
}

// Marshal 序列化对象
func Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

// Unmarshal 反序列化数据到给定的指针对象
func Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}
