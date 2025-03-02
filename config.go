package wtoken

// Config 定义了Token管理器的配置
type Config struct {
	// CacheFilePath：缓存文件路径及文件名，为空则不启用缓存
	CacheFilePath string
	// Language：错误信息语言类型，支持中文(zh)和英文(en)
	Language Language
	// MaxTokens：最大缓存Token数量，如果小于等于0则不限制数量
	MaxTokens int
	// Debug：是否启用调试模式，启用后会输出详细的日志信息
	Debug bool
	// Delimiter：分隔符，权限字符串分割符，默认是空格
	Delimiter string
}
