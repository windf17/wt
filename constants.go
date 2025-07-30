package wtoken

// 文件权限常量
const (
	// DIR_PERM 目录创建权限
	DIR_PERM = 0755
	// FILE_PERM 文件创建权限
	FILE_PERM = 0644
)

// Token相关常量
const (
	// TOKEN_BYTE_SIZE Token随机字节大小
	TOKEN_BYTE_SIZE = 24
	// TIMESTAMP_BYTE_SIZE 时间戳字节大小
	TIMESTAMP_BYTE_SIZE = 8
)

// 默认配置常量
const (

	// DEFAULT_MAX_TOKENS 默认最大Token数量
	DEFAULT_MAX_TOKENS = 10000
	// DEFAULT_DELIMITER 默认分隔符
	DEFAULT_DELIMITER = " "
	// DEFAULT_TOKEN_RENEW_TIME 默认Token续期时间
	DEFAULT_TOKEN_RENEW_TIME = "10m"
)

// 时间单位常量
const (
	// SECONDS_PER_MINUTE 每分钟秒数
	SECONDS_PER_MINUTE = 60
	// SECONDS_PER_HOUR 每小时秒数
	SECONDS_PER_HOUR = 3600
	// SECONDS_PER_DAY 每天秒数
	SECONDS_PER_DAY = 86400
)

// 日志级别常量
const (
	// DEBUG 调试级别
	DEBUG = 2
	// INFO 信息级别
	INFO = 1
	// ERROR 错误级别
	ERROR = 0
)