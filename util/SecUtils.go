package util

import (
	"gitee.com/chunanyong/gouuid"
)

// GetUUIDString 获取一个UUID字符串
func GetUUIDString() string {
	uuid, err := gouuid.NewV4()
	if err != nil {
		return ""
	}
	return uuid
}
