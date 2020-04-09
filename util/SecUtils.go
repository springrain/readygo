package util

import (
	"crypto/md5"
	"encoding/hex"

	"gitee.com/chunanyong/gouuid"
)

// GenerateUUIDString 获取一个UUID字符串
func GenerateUUIDString() string {
	uuid, err := gouuid.NewV4()
	if err != nil {
		return ""
	}
	return uuid.String()
}

//GenerateMD5 生成md5字符串
func GenerateMD5(str string) string {
	ctx := md5.New()
	ctx.Write([]byte(str))
	return hex.EncodeToString(ctx.Sum(nil))
}
