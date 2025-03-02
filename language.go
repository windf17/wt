package wtoken

// Language 定义语言类型
type Language string

const (
	// LangChinese 中文
	LangChinese Language = "zh"
	// LangEnglish 英文
	LangEnglish Language = "en"
)

// GetCode 实现ILanguage接口
func (l Language) GetCode() string {
	return string(l)
}

// RegisterLanguage 注册新的语言类型
func RegisterLanguage(lang string) Language {
	return Language(lang)
}
