package wtoken

import (
	"fmt"
	"runtime"
	"time"
)

// ErrorLevel 错误级别
type ErrorLevel int

const (
	ERROR_LEVEL_DEBUG ErrorLevel = iota
	ERROR_LEVEL_INFO
	ERROR_LEVEL_WARN
	ERROR_LEVEL_ERROR
	ERROR_LEVEL_FATAL
)

// ErrorInfo 错误信息结构
type ErrorInfo struct {
	Code      ErrorCode  `json:"code"`
	Message   string     `json:"message"`
	Level     ErrorLevel `json:"level"`
	Timestamp time.Time  `json:"timestamp"`
	File      string     `json:"file,omitempty"`
	Line      int        `json:"line,omitempty"`
	Function  string     `json:"function,omitempty"`
	Context   any        `json:"context,omitempty"`
}

// ErrorHandler 错误处理器
type ErrorHandler struct {
	// 移除logger依赖，使用简化的错误处理
}

// NewErrorHandler 创建新的错误处理器
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

/**
 * HandleError 处理错误
 * @param {ErrorCode} code 错误码
 * @param {string} message 错误消息
 * @param {ErrorLevel} level 错误级别
 * @param {any} context 错误上下文
 * @returns {*ErrorInfo} 错误信息
 */
func (eh *ErrorHandler) HandleError(code ErrorCode, message string, level ErrorLevel, context any) *ErrorInfo {
	// 获取调用栈信息
	pc, file, line, ok := runtime.Caller(1)
	var function string
	if ok {
		function = runtime.FuncForPC(pc).Name()
	}

	errorInfo := &ErrorInfo{
		Code:      code,
		Message:   message,
		Level:     level,
		Timestamp: time.Now(),
		File:      file,
		Line:      line,
		Function:  function,
		Context:   context,
	}

	// 简化错误处理，不再记录日志
	// 可以在这里添加其他错误处理逻辑

	return errorInfo
}

/**
 * HandleValidationError 处理验证错误
 * @param {string} field 字段名
 * @param {string} value 字段值
 * @param {string} rule 验证规则
 * @returns {*ErrorInfo} 错误信息
 */
func (eh *ErrorHandler) HandleValidationError(field, value, rule string) *ErrorInfo {
	message := fmt.Sprintf("字段 '%s' 验证失败: 值 '%s' 不符合规则 '%s'", field, value, rule)
	context := map[string]any{
		"field": field,
		"value": value,
		"rule":  rule,
	}
	return eh.HandleError(E_InvalidParams, message, ERROR_LEVEL_WARN, context)
}

/**
 * HandleSystemError 处理系统错误
 * @param {error} err 系统错误
 * @param {string} component 组件名称
 * @returns {*ErrorInfo} 错误信息
 */
func (eh *ErrorHandler) HandleSystemError(err error, component string) *ErrorInfo {
	message := fmt.Sprintf("系统错误在组件 '%s': %v", component, err)
	context := map[string]any{
		"component": component,
		"error":     err.Error(),
	}
	return eh.HandleError(E_SystemError, message, ERROR_LEVEL_FATAL, context)
}

/**
 * HandleConcurrencyError 处理并发错误
 * @param {string} resource 资源名称
 * @param {string} operation 操作名称
 * @returns {*ErrorInfo} 错误信息
 */
func (eh *ErrorHandler) HandleConcurrencyError(resource, operation string) *ErrorInfo {
	message := fmt.Sprintf("并发冲突: 资源 '%s' 在操作 '%s' 时发生冲突", resource, operation)
	context := map[string]any{
		"resource":  resource,
		"operation": operation,
	}
	return eh.HandleError(E_DBDeadlock, message, ERROR_LEVEL_ERROR, context)
}

// convertToLogLevel 函数已移除，不再需要转换为LogLevel

/**
 * IsRecoverableError 判断错误是否可恢复
 * @param {ErrorCode} code 错误码
 * @returns {bool} 是否可恢复
 */
func IsRecoverableError(code ErrorCode) bool {
	switch code {
	case E_InvalidToken, E_TokenExpired, E_InvalidParams, E_UserInvalid, E_GroupInvalid:
		return true
	case E_Internal, E_System, E_SystemError, E_DatabaseError:
		return false
	default:
		return true
	}
}

/**
 * GetErrorSeverity 获取错误严重程度
 * @param {ErrorCode} code 错误码
 * @returns {ErrorLevel} 错误级别
 */
func GetErrorSeverity(code ErrorCode) ErrorLevel {
	switch code {
	case E_SystemError, E_DatabaseError:
		return ERROR_LEVEL_ERROR
	case E_InvalidToken, E_TokenExpired:
		return ERROR_LEVEL_WARN
	case E_InvalidParams, E_UserInvalid, E_GroupInvalid:
		return ERROR_LEVEL_INFO
	default:
		return ERROR_LEVEL_DEBUG
	}
}
