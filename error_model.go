package wtoken

type ErrorData struct {
	Lan  Language  // 错误语言
	Code ErrorCode // 错误码
}

func (e ErrorData) Error() string {
	return defaultRegistry.GetErrorMessage(e.Lan, e.Code)
}

func (e ErrorData) String() string {
	return defaultRegistry.GetErrorMessage(e.Lan, e.Code)
}
