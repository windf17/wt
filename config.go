package wtoken

// Config 定义了Token管理器的配置
type Config struct {
	// CacheFilePath：缓存文件路径及文件名，为空则不启用缓存
	CacheFilePath string
	// MaxTokens：最大缓存Token数量，如果小于等于0则不限制数量
	MaxTokens int
	// Debug：是否启用调试模式，启用后会输出详细的日志信息
	Debug bool
	// Delimiter：分隔符，权限字符串分割符，默认是空格
	Delimiter string
	// TokenRenewTime：Token续期时间，单位秒，默认10分钟
	TokenRenewTime int64
}

type ConfigRaw struct {
	// CacheFilePath：缓存文件路径及文件名，为空则不启用缓存
	CacheFilePath string `json:"cacheFilePath"`
	// Language：错误信息语言类型，支持中文(zh)和英文(en)
	Language Language `json:"language"`
	// MaxTokens：最大缓存Token数量，如果小于等于0则不限制数量
	MaxTokens int `json:"maxTokens"`
	// Debug：是否启用调试模式，启用后会输出详细的日志信息
	Debug bool `json:"debug"`
	// Delimiter：分隔符，权限字符串分割符，默认是空格
	Delimiter string `json:"delimiter"`
	// TokenRenewTime：Token续期时间，单位秒，默认10分钟
	TokenRenewTime string `json:"tokenRenewTime"`
}

var DefaultConfigRaw = &ConfigRaw{
	CacheFilePath:  "./token.cache",
	Language:       LangChinese,
	MaxTokens:      10000,
	Debug:          false,
	Delimiter:      " ",
	TokenRenewTime: "10m",
}
