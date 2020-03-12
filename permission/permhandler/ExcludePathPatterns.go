/*
 * @Author: your name
 * @Date: 2020-03-11 22:06:03
 * @LastEditTime: 2020-03-11 22:25:21
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \readygo\permission\permhandler\ExcludePathPatterns.go
 */

package permhandler

import (
	"strings"
)

var excludePathPatterns []string

func init() {
	excludePathPatterns = make([]string, 16)
	AddExcluePath("/login")
}

// AddExcluePath 添加排除目录
func AddExcluePath(path string) {
	excludePathPatterns = append(excludePathPatterns, strings.ToLower(path))
}

// IsExcludePath 是否排除目录
func IsExcludePath(path string) bool {
	for _, item := range excludePathPatterns {
		if item == strings.ToLower(path) {
			return true
		}
	}
	return false
}
