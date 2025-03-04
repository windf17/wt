package wtoken

// ErrorCode 定义错误码类型
type ErrorCode int

const (
	// ErrSuccess 成功
	ErrSuccess ErrorCode = 0
	// ErrUnknown 未知错误
	ErrUnknown ErrorCode = 9999

	// token错误码1101开头
	// ErrInvalidToken 无效的token
	ErrInvalidToken ErrorCode = 1101
	// ErrTokenExpired token已过期
	ErrTokenExpired ErrorCode = 1102
	// ErrTokenNotFound token不存在
	ErrTokenNotFound ErrorCode = 1103
	// ErrTokenLimitExceeded 超出token数量限制
	ErrTokenLimitExceeded ErrorCode = 1104
	// ErrAddToken 生成token错误
	ErrAddToken ErrorCode = 1105

	// 用户错误码，1201开头
	// ErrInvalidUserID 无效的用户ID
	ErrInvalidUserID ErrorCode = 1201
	// ErrUserIDNotFound 用户ID不存在
	ErrUserIDNotFound ErrorCode = 1202
	// ErrTypeAssertionError 类型断言错误
	ErrTypeAssertionError ErrorCode = 1203

	// 用户组错误码，1301开头
	// ErrGroupNotFound 用户组ID不存在
	ErrGroupNotFound ErrorCode = 1301
	// ErrInvalidGroupID 无效的用户组ID
	ErrInvalidGroupID ErrorCode = 1302

	// IP错误码，1401开头
	// ErrInvalidIP 无效的IP地址
	ErrInvalidIP ErrorCode = 1401
	// ErrIPMismatch IP地址不匹配
	ErrIPMismatch ErrorCode = 1402

	// 配置错误码，1501开头
	// ErrInvalidConfig 无效的配置
	ErrInvalidConfig ErrorCode = 1501

	// 缓存错误码，1601开头
	// ErrCacheFileLoadFailed 加载缓存文件失败
	ErrCacheFileLoadFailed ErrorCode = 1601
	// ErrCacheFileParseFailed 缓存文件解析错误
	ErrCacheFileParseFailed ErrorCode = 1602

	// API错误码，1700开头
	// ErrAccessDenied 访问被禁止的API
	ErrAccessDenied ErrorCode = 1701
	// ErrInvalidURL 无效的URL
	ErrInvalidURL ErrorCode = 1702
	// 该用户的用户组没有定制API访问权限
	ErrNoAPIPermission ErrorCode = 1703

	// ErrInternalError 内部错误，1901开头
	ErrInternalError ErrorCode = 1901
)

func (e ErrorCode) Error() string {
	return defaultRegistry.getErrorMessage(e)
}

func (e ErrorCode) String() string {
	return defaultRegistry.getErrorMessage(e)
}
