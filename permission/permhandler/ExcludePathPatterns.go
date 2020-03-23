package permhandler

//全局排除的路径
var excludePathPatterns map[string]bool

//登录用户默认能访问的路径
var userDefaultPathPatterns map[string]bool

func init() {
	excludePathPatterns = make(map[string]bool)
	AddExcluePath("/login")
}

// AddExcluePath 添加排除目录
func AddExcluePath(path string) {
	excludePathPatterns[path] = true
}

// AddUserDefaulPath 添加登录用户默认能访问的路径
func AddUserDefaulPath(path string) {
	userDefaultPathPatterns[path] = true
}

// isExcludePath 是否排除目录
func isExcludePath(path string) bool {
	has := excludePathPatterns[path]
	return has
}

// isUserDefaultPath 是否登录用户默认能访问的路径
func isUserDefaultPath(path string) bool {
	has := userDefaultPathPatterns[path]
	return has
}
