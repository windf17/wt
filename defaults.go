package wtoken

// Token相关默认值
const (
	// DefaultTokenExpireSeconds 默认Token有效期（秒）
	// 默认24小时 = 24 * 3600秒
	DefaultTokenExpireSeconds = 24 * 3600

	// DefaultTokenRefreshSeconds 默认Token刷新时间（秒）
	// 默认30分钟 = 30 * 60秒
	DefaultTokenRefreshSeconds = 30 * 60

	// DefaultMaxTokens 默认最大Token数量，小于等于0表示不限制
	DefaultMaxTokens = 10000

	// DefaultCacheFilePath 默认缓存文件路径，包含文件名
	DefaultCacheFilePath = "./token.cache"

	// DefaultLogPath 默认日志文件路径，不包含文件名
	DefaultLogPath = "./log"
)

// 用户组相关默认值
const (
	// DefaultAllowMultipleLogin 默认是否允许多设备登录
	DefaultAllowMultipleLogin = false

	// DefaultUseChinese 默认是否使用中文
	DefaultUseChinese = false

	// DefaultDebug 默认是否开启调试模式
	DefaultDebug = false

	// DefaultDelimiter 默认分隔符
	DefaultDelimiter = " "
)

// DefaultPublicAPIs 默认公开API列表
var DefaultPublicAPIs = []string{
	"/api/public/*",
	"/api/health",
	"/api/version",
}
