package util

import (
	"crypto/md5"
	"encoding/hex"
)

//GenerateMD5 生成md5字符串
func GenerateMD5(str string) string {
	ctx := md5.New()
	ctx.Write([]byte(str))
	return hex.EncodeToString(ctx.Sum(nil))
}
