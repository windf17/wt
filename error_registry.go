package wtoken

import "fmt"

// ErrorRegistry 错误信息注册中心
type ErrorRegistry struct {
	// messages 存储所有错误信息
	messages map[Language]map[ErrorCode]string
	// language 默认语言
	language Language
}

// defaultRegistry 默认的错误信息注册中心
var defaultRegistry = &ErrorRegistry{
	messages: make(map[Language]map[ErrorCode]string),
	language: LangChinese, // 默认语言为中文
}

// registerErrorMessage 注册错误信息
func (r *ErrorRegistry) registerErrorMessage(lang Language, code ErrorCode, message string) {
	if _, exists := r.messages[lang]; !exists {
		r.messages[lang] = make(map[ErrorCode]string)
	}
	r.messages[lang][code] = message
}

// getErrorMessage 获取错误信息
func (r *ErrorRegistry) getErrorMessage(code ErrorCode) string {
	// 检查语言是否存在
	if messages, ok := r.messages[r.language]; ok {
		// 检查错误码是否存在
		if message, ok := messages[code]; ok {
			return message
		}
	}

	// 如果找不到对应的错误信息，返回未知错误信息
	if unknownMsg, ok := r.messages[r.language][E_Unknown]; ok {
		return fmt.Sprintf("%s, Code: %d)", unknownMsg, code)
	}

	// 如果连未知错误信息都找不到，返回固定的错误信息
	return fmt.Sprintf("Unknown error code , Code: %d)", code)
}

func init() {
	defaultRegistry.registerErrorMessage(LangChinese, E_Unknown, "未知错误码")
	defaultRegistry.registerErrorMessage(LangEnglish, E_Unknown, "Unknown error code")
	for lang, messages := range ErrorI18nMessages {
		for code, message := range messages {
			defaultRegistry.registerErrorMessage(lang, code, message)
		}
	}
}
