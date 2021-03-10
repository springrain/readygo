/*
 * @Author: your name
 * @Date: 2021-03-04 11:36:07
 * @LastEditTime: 2021-03-05 17:21:03
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \readygo\permission\permhandler\ExcludePathPatterns.go
 */
package permhandler

//bug(springrain) 需要支持url的正则表达式
//全局排除的路径
var excludePathPatterns map[string]bool

//登录用户默认能访问的路径
var userDefaultPathPatterns map[string]bool

func init() {
	excludePathPatterns = make(map[string]bool)
	AddExcluePath("/login")
	AddExcluePath("/swagger/")
}

// AddExcluePath 添加排除目录
func AddExcluePath(path string) {
	excludePathPatterns[path] = true
}

// AddUserDefaulPath 添加登录用户默认能访问的路径
func AddUserDefaulPath(path string) {
	userDefaultPathPatterns[path] = true
}

// isExcludePath 是否排除目录,需要验证正则
func isExcludePath(path string) bool {
	has := excludePathPatterns[path]
	return has
}

// isUserDefaultPath 是否登录用户默认能访问的路径
func isUserDefaultPath(path string) bool {
	has := userDefaultPathPatterns[path]
	return has
}
