package models

// Config 定义了Token管理器的配置
type Config struct {
	// Language：错误信息语言类型，"zh"为中文，其他为英文
	Language string
	// MaxTokens：最大缓存Token数量，如果小于等于0则不限制数量
	MaxTokens int
	// Delimiter：分隔符，权限字符串分割符，默认是空格
	Delimiter string
	// TokenRenewTime：Token续期时间，单位秒，默认10分钟
	TokenRenewTime int64
}

type ConfigRaw struct {
	// Language：错误信息语言类型，"zh"为中文，其他为英文
	Language string `json:"language"`
	// Delimiter：分隔符，权限字符串分割符，默认是空格
	Delimiter string `json:"delimiter"`
	// MaxTokens：最大缓存Token数量，如果小于等于0则不限制数量
	MaxTokens int `json:"maxTokens"`
	// TokenRenewTime：Token续期时间，单位秒，默认10分钟
	TokenRenewTime string `json:"tokenRenewTime"`
}