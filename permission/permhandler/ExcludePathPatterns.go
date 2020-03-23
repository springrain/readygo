/*
 * @Author: your name
 * @Date: 2020-03-11 22:06:03
 * @LastEditTime: 2020-03-11 22:25:21
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \readygo\permission\permhandler\ExcludePathPatterns.go
 */

package permhandler

var excludePathPatterns map[string]bool

func init() {
	excludePathPatterns = make(map[string]bool)
	AddExcluePath("/login")
}

// AddExcluePath 添加排除目录
func AddExcluePath(path string) {
	excludePathPatterns[path] = true
}

// isExcludePath 是否排除目录
func isExcludePath(path string) bool {
	has := excludePathPatterns[path]
	return has
}
