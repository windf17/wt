package wtoken

// Config 定义了Token管理器的配置
type Config struct {
	// Language：错误信息语言类型
	Language Language
	// MaxTokens：最大缓存Token数量，如果小于等于0则不限制数量
	MaxTokens int
	// Delimiter：分隔符，权限字符串分割符，默认是空格
	Delimiter string
	// TokenRenewTime：Token续期时间，单位秒，默认10分钟
	TokenRenewTime int64
	// MinTokenExpire：最小token过期时间（秒），默认60秒
	MinTokenExpire int
	// MaxTokenExpire：最大token过期时间（秒），默认1年
	MaxTokenExpire int

}

type ConfigRaw struct {
	// Language：错误信息语言类型，支持中文(zh)和英文(en)
	Language Language `json:"language"`
	// Delimiter：分隔符，权限字符串分割符，默认是空格
	Delimiter string `json:"delimiter"`
	// MaxTokens：最大缓存Token数量，如果小于等于0则不限制数量
	MaxTokens int `json:"maxTokens"`
	// TokenRenewTime：Token续期时间，单位秒，默认10分钟
	TokenRenewTime string `json:"tokenRenewTime"`
	// MinTokenExpire：最小token过期时间（秒），默认60秒
	MinTokenExpire int `json:"minTokenExpire"`
	// MaxTokenExpire：最大token过期时间（秒），默认30天
	MaxTokenExpire int `json:"maxTokenExpire"`

}

var DefaultConfigRaw = &ConfigRaw{
	Language:       LangChinese,
	MaxTokens:      DEFAULT_MAX_TOKENS,
	Delimiter:      DEFAULT_DELIMITER,
	TokenRenewTime: DEFAULT_TOKEN_RENEW_TIME,
	MinTokenExpire: 60,                   // 60秒
	MaxTokenExpire: 2592000,             // 30天

}
