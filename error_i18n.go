package wtoken

// Language 定义语言类型
type Language string

const (
	// LangChinese 中文
	LangChinese Language = "zh"
	// LangEnglish 英文
	LangEnglish Language = "en"
)

// ErrorI18nMessages 定义多语言错误信息映射
var ErrorI18nMessages = map[Language]map[ErrorCode]string{
	LangChinese: {
		ErrSuccess:              "成功",
		ErrInvalidToken:         "无效的token",
		ErrTokenExpired:         "token已过期",
		ErrTokenNotFound:        "token不存在",
		ErrTokenLimitExceeded:   "超出token数量限制",
		ErrInvalidUserID:        "无效的用户ID",
		ErrUserIDNotFound:       "用户ID不存在",
		ErrGroupNotFound:        "用户组不存在",
		ErrInvalidGroupID:       "无效的用户组ID",
		ErrInvalidIP:            "无效的IP地址",
		ErrIPMismatch:           "IP地址不匹配",
		ErrInvalidConfig:        "无效的配置",
		ErrCacheFileLoadFailed:  "加载缓存文件失败",
		ErrCacheFileParseFailed: "缓存文件解析错误",
		ErrAccessDenied:         "无权访问该API",
		ErrInternalError:        "内部错误",
		ErrTypeAssertionError:   "类型断言错误",
		ErrInvalidURL:           "无效的URL",
		ErrNoAPIPermission:      "该用户的用户组没有定制API访问权限",
		ErrUnknown:              "未知错误",
	},
	LangEnglish: {
		ErrSuccess:              "Success",
		ErrInvalidToken:         "Invalid token",
		ErrTokenExpired:         "Token expired",
		ErrTokenNotFound:        "Token not found",
		ErrTokenLimitExceeded:   "Token limit exceeded",
		ErrInvalidUserID:        "Invalid user ID",
		ErrUserIDNotFound:       "User ID not found",
		ErrGroupNotFound:        "User group not found",
		ErrInvalidGroupID:       "Invalid user group ID",
		ErrInvalidIP:            "Invalid IP address",
		ErrIPMismatch:           "IP address mismatch",
		ErrInvalidConfig:        "Invalid configuration",
		ErrCacheFileLoadFailed:  "Cache file load failed",
		ErrCacheFileParseFailed: "Cache file parse failed",
		ErrAccessDenied:         "API access not allowed",
		ErrInternalError:        "Internal error",
		ErrTypeAssertionError:   "Type assertion error",
		ErrInvalidURL:           "Invalid URL",
		ErrNoAPIPermission:      "The user group does not have customized API access permissions",
		ErrUnknown:              "Unknown error",
	},
}
