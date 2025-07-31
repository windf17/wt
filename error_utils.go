package wt

/**
 * getErrorMessage 根据语言获取错误信息
 * @param {string} language 语言类型，"zh"为中文，其他为英文
 * @param {string} errorKey 错误键值
 * @returns {string} 错误信息
 */
func getErrorMessage(language, errorKey string) string {
	if language == "zh" {
		return getChineseErrorMessage(errorKey)
	}
	return getEnglishErrorMessage(errorKey)
}

/**
 * getChineseErrorMessage 获取中文错误信息
 * @param {string} errorKey 错误键值
 * @returns {string} 中文错误信息
 */
func getChineseErrorMessage(errorKey string) string {
	errorMessages := map[string]string{
		"success":             "成功",
		"internal":            "系统内部错误",
		"system":              "系统错误",
		"config_invalid":      "配置无效",
		"invalid_params":      "无效参数",
		"api_not_found":       "接口不存在",
		"unauthorized":        "未授权访问",
		"forbidden":           "禁止访问",
		"invalid_token":       "无效令牌",
		"token_expired":       "令牌过期",
		"token_limit":         "令牌数量超限",
		"token_generate":      "令牌生成失败",
		"invalid_auth_header": "请求头认证格式错误",
		"captcha_invalid":     "验证码错误",
		"invalid_ip":          "无效IP地址",
		"ip_mismatch":         "IP地址不匹配",
		"ip_not_allowed":      "IP未授权",
		"not_login":           "用户未登录",
		"user_logout":         "用户已登出",
		"user_not_found":      "用户不存在",
		"user_invalid":        "用户无效",
		"user_exists":         "用户已存在",
		"user_disabled":       "用户已禁用",
		"user_pending":        "用户审核中",
		"group_not_found":     "用户组不存在",
		"group_invalid":       "用户组无效",
		"group_exists":        "用户组已存在",
		"group_disabled":      "用户组已禁用",
		"group_pending":       "用户组审核中",
		"password_invalid":    "密码错误",
		"password_expired":    "密码已过期",
		"password_weak":       "密码强度不足",
		"password_reuse":      "密码重复使用",
		"password_locked":     "密码已锁定",
		"db_connect":          "数据库连接失败",
		"db_query":            "数据库查询失败",
		"db_query_not_found":  "记录不存在",
		"db_insert":           "数据库插入失败",
		"db_update":           "数据库更新失败",
		"db_delete":           "数据库删除失败",
		"db_transaction":      "事务操作失败",
		"db_deadlock":         "数据库死锁",
		"db_timeout":          "数据库操作超时",
		"db_duplicate":        "唯一键冲突",
		"db_foreign_key":      "外键约束违反",
		"unknown":             "未知错误",
	}

	if msg, exists := errorMessages[errorKey]; exists {
		return msg
	}
	return "未知错误"
}

/**
 * getEnglishErrorMessage 获取英文错误信息
 * @param {string} errorKey 错误键值
 * @returns {string} 英文错误信息
 */
func getEnglishErrorMessage(errorKey string) string {
	errorMessages := map[string]string{
		"success":             "Success",
		"internal":            "Internal system error",
		"system":              "System error",
		"config_invalid":      "Invalid configuration",
		"invalid_params":      "Invalid parameters",
		"api_not_found":       "API not found",
		"unauthorized":        "Unauthorized access",
		"forbidden":           "Access forbidden",
		"invalid_token":       "Invalid token",
		"token_expired":       "Token expired",
		"token_limit":         "Token limit exceeded",
		"token_generate":      "Token generation failed",
		"invalid_auth_header": "Invalid authorization header",
		"captcha_invalid":     "Invalid captcha",
		"invalid_ip":          "Invalid IP address",
		"ip_mismatch":         "IP address mismatch",
		"ip_not_allowed":      "IP not allowed",
		"not_login":           "User not logged in",
		"user_logout":         "User logged out",
		"user_not_found":      "User not found",
		"user_invalid":        "Invalid user",
		"user_exists":         "User already exists",
		"user_disabled":       "User disabled",
		"user_pending":        "User verification pending",
		"group_not_found":     "User group not found",
		"group_invalid":       "Invalid user group",
		"group_exists":        "User group already exists",
		"group_disabled":      "User group disabled",
		"group_pending":       "User group verification pending",
		"password_invalid":    "Invalid password",
		"password_expired":    "Password expired",
		"password_weak":       "Password strength insufficient",
		"password_reuse":      "Password reused",
		"password_locked":     "Password locked",
		"db_connect":          "Database connection failed",
		"db_query":            "Database query failed",
		"db_query_not_found":  "Record not found",
		"db_insert":           "Database insert failed",
		"db_update":           "Database update failed",
		"db_delete":           "Database delete failed",
		"db_transaction":      "Transaction failed",
		"db_deadlock":         "Database deadlock",
		"db_timeout":          "Database operation timeout",
		"db_duplicate":        "Duplicate key violation",
		"db_foreign_key":      "Foreign key violation",
		"unknown":             "Unknown error",
	}

	if msg, exists := errorMessages[errorKey]; exists {
		return msg
	}
	return "Unknown error"
}