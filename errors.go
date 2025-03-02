package wtoken

import "strconv"

// ErrorCode 定义错误码类型
type ErrorCode int

const (
	// ErrCodeSuccess 成功
	ErrCodeSuccess ErrorCode = 0

	// token错误码1101开头
	// ErrCodeInvalidToken 无效的token
	ErrCodeInvalidToken ErrorCode = 1101
	// ErrCodeTokenExpired token已过期
	ErrCodeTokenExpired ErrorCode = 1102
	// ErrCodeTokenNotFound token不存在
	ErrCodeTokenNotFound ErrorCode = 1103
	// ErrCodeTokenLimitExceeded 超出token数量限制
	ErrCodeTokenLimitExceeded ErrorCode = 1104
	// ErrCodeAddToken 生成token错误
	ErrCodeAddToken ErrorCode = 1105

	// 用户错误码，1201开头
	// ErrCodeInvalidUserID 无效的用户ID
	ErrCodeInvalidUserID ErrorCode = 1201
	// ErrCodeUserIDNotFound 用户ID不存在
	ErrCodeUserIDNotFound ErrorCode = 1202
	// ErrCodeTypeAssertionError 类型断言错误
	ErrCodeTypeAssertionError ErrorCode = 1203

	// 用户组错误码，1301开头
	// ErrCodeGroupNotFound 用户组ID不存在
	ErrCodeGroupNotFound ErrorCode = 1301
	// ErrCodeInvalidGroupID 无效的用户组ID
	ErrCodeInvalidGroupID ErrorCode = 1302

	// IP错误码，1401开头
	// ErrCodeInvalidIP 无效的IP地址
	ErrCodeInvalidIP ErrorCode = 1401
	// ErrCodeIPMismatch IP地址不匹配
	ErrCodeIPMismatch ErrorCode = 1402

	// 配置错误码，1501开头
	// ErrCodeInvalidConfig 无效的配置
	ErrCodeInvalidConfig ErrorCode = 1501

	// 缓存错误码，1601开头
	// ErrCodeCacheFileLoadFailed 加载缓存文件失败
	ErrCodeCacheFileLoadFailed ErrorCode = 1601
	// ErrCodeCacheFileParseFailed 缓存文件解析错误
	ErrCodeCacheFileParseFailed ErrorCode = 1602

	// API错误码，1700开头
	// ErrCodeAccessDenied 访问被禁止的API
	ErrCodeAccessDenied ErrorCode = 1701
	// ErrCodeInvalidURL 无效的URL
	ErrCodeInvalidURL ErrorCode = 1702

	// ErrCodeInternalError 内部错误，1901开头
	ErrCodeInternalError ErrorCode = 1901
)

func (c ErrorCode) String() string {
	return strconv.Itoa(int(c))
}
