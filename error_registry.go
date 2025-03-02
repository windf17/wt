package wtoken

import "fmt"

// ILanguage 定义语言接口
type ILanguage interface {
	// GetCode 获取语言代码
	GetCode() string
}

// ErrorRegistry 错误信息注册中心
type ErrorRegistry struct {
	// messages 存储所有错误信息
	// map[Language]map[ErrorCode]string
	messages map[string]map[ErrorCode]string
	// defaultLanguage 默认语言
	defaultLanguage string
	// unknownErrorCode 未知错误码
	unknownErrorCode ErrorCode
}

// defaultRegistry 默认的错误信息注册中心
var defaultRegistry = NewErrorRegistry()

// NewErrorRegistry 创建错误信息注册中心
func NewErrorRegistry() *ErrorRegistry {
	return &ErrorRegistry{
		messages:         make(map[string]map[ErrorCode]string),
		defaultLanguage:  string(LangEnglish),
		unknownErrorCode: ErrorCode(9999),
	}
}

// RegisterLanguage 注册新的语言
func (r *ErrorRegistry) RegisterLanguage(lang ILanguage) {
	if _, exists := r.messages[lang.GetCode()]; !exists {
		r.messages[lang.GetCode()] = make(map[ErrorCode]string)
	}
}

// RegisterErrorMessage 注册错误信息
func (r *ErrorRegistry) RegisterErrorMessage(lang ILanguage, code ErrorCode, message string) {
	r.RegisterLanguage(lang)
	r.messages[lang.GetCode()][code] = message
}

// GetErrorMessage 获取错误信息
func (r *ErrorRegistry) GetErrorMessage(lang ILanguage, code ErrorCode) string {
	// 检查语言是否存在
	if messages, ok := r.messages[lang.GetCode()]; ok {
		// 检查错误码是否存在
		if message, ok := messages[code]; ok {
			return message
		}
	}

	// 如果找不到对应的错误信息，返回未知错误信息
	if unknownMsg, ok := r.messages[r.defaultLanguage][r.unknownErrorCode]; ok {
		return fmt.Sprintf("%s (Language: %s, Code: %d)", unknownMsg, lang.GetCode(), code)
	}

	// 如果连未知错误信息都找不到，返回固定的错误信息
	return fmt.Sprintf("Unknown error code (Language: %s, Code: %d)", lang.GetCode(), code)
}

// SetDefaultLanguage 设置默认语言
func (r *ErrorRegistry) SetDefaultLanguage(lang ILanguage) {
	r.defaultLanguage = lang.GetCode()
}

// SetUnknownErrorCode 设置未知错误码
func (r *ErrorRegistry) SetUnknownErrorCode(code ErrorCode) {
	r.unknownErrorCode = code
}

// RegisterDefaultMessages 注册默认的错误信息
func (r *ErrorRegistry) RegisterDefaultMessages() {
	// 注册未知错误码
	r.RegisterErrorMessage(LangChinese, ErrorCode(9999), "未知错误码")
	r.RegisterErrorMessage(LangEnglish, ErrorCode(9999), "Unknown error code")

	// 注册现有的错误信息
	for lang, messages := range ErrorI18nMessages {
		for code, message := range messages {
			r.RegisterErrorMessage(lang, code, message)
		}
	}
}
