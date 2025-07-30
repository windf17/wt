package wtoken

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

// tokenRegex Token格式验证正则表达式
var tokenRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// 文件权限常量
const (
	// FILE_PERM_DIR 目录权限
	FILE_PERM_DIR = 0755
	// FILE_PERM_FILE 文件权限
	FILE_PERM_FILE = 0644
	// FILE_PERM_PRIVATE 私有文件权限
	FILE_PERM_PRIVATE = 0600
)

// 验证规则常量
const (
	// MIN_TOKEN_EXPIRE 最小token过期时间（秒）
	MIN_TOKEN_EXPIRE = 60
	// MAX_TOKEN_EXPIRE 最大token过期时间（秒）
	MAX_TOKEN_EXPIRE = 86400 // 1天
	// MAX_USERNAME_LENGTH 最大用户名长度
	MAX_USERNAME_LENGTH = 50
	// MAX_EMAIL_LENGTH 最大邮箱长度
	MAX_EMAIL_LENGTH = 100
	// MAX_API_PATH_LENGTH 最大API路径长度
	MAX_API_PATH_LENGTH = 200
)

// // 正则表达式模式
// var (
// 	// emailRegex 邮箱格式正则表达式
// 	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
// 	// usernameRegex 用户名格式正则表达式（字母、数字、下划线、中划线）
// 	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
// 	// tokenRegex token格式正则表达式（Base64字符，包括URL安全字符）
// 	tokenRegex = regexp.MustCompile(`^[A-Za-z0-9+/_-]+=*$`)
// 	// apiPathRegex API路径格式正则表达式
// 	apiPathRegex = regexp.MustCompile(`^[a-zA-Z0-9/_.-]+$`)
// )

/**
 * ValidateIPAddress 验证IP地址格式
 * @param {string} ip IP地址字符串
 * @returns {error} 验证错误
 */
func ValidateIPAddress(ip string) error {
	if strings.TrimSpace(ip) == "" {
		return errors.New("IP地址不能为空")
	}
	if net.ParseIP(ip) == nil {
		return errors.New("IP地址格式无效")
	}
	return nil
}

/**
 * ValidateTokenExpire 验证token过期时间
 * @param {int} expire 过期时间（秒）
 * @param {*Config} config 配置信息
 * @returns {error} 验证错误
 */
func ValidateTokenExpire(expire int, config *Config) error {
	minExpire := MIN_TOKEN_EXPIRE
	maxExpire := MAX_TOKEN_EXPIRE
	if config != nil {
		if config.MinTokenExpire > 0 {
			minExpire = config.MinTokenExpire
		}
		if config.MaxTokenExpire > 0 {
			maxExpire = config.MaxTokenExpire
		}
	}
	if expire < minExpire {
		return fmt.Errorf("token过期时间不能少于%d秒", minExpire)

	}
	if expire > maxExpire {
		return fmt.Errorf("token过期时间不能超过%d秒", maxExpire)
	}
	return nil
}

/**
 * ValidateTokenKey 验证token键格式
 * @param {string} token token字符串
 * @param {*Config} config 配置信息
 * @returns {error} 验证错误
 */
func ValidateTokenKey(token string, config *Config) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("token不能为空")
	}
	if !tokenRegex.MatchString(token) {
		return errors.New("token格式无效")
	}
	return nil
}

/**
 * ValidateStringNotEmpty 验证字符串非空
 * @param {string} value 待验证的字符串
 * @param {string} fieldName 字段名称
 * @returns {error} 验证错误
 */
func ValidateStringNotEmpty(value string, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(fieldName + "不能为空")
	}
	return nil
}

/**
 * ValidateStringLength 验证字符串长度
 * @param {string} value 待验证的字符串
 * @param {string} fieldName 字段名称
 * @param {int} minLength 最小长度
 * @param {int} maxLength 最大长度
 * @returns {error} 验证错误
 */
func ValidateStringLength(value string, fieldName string, minLength int, maxLength int) error {
	value = strings.TrimSpace(value)
	if len(value) < minLength {
		return errors.New(fieldName + "长度不能少于" + string(rune(minLength)) + "个字符")
	}
	if len(value) > maxLength {
		return errors.New(fieldName + "长度不能超过" + string(rune(maxLength)) + "个字符")
	}
	return nil
}

/**
 * ValidatePositiveInt 验证正整数
 * @param {int} value 待验证的整数
 * @param {string} fieldName 字段名称
 * @returns {error} 验证错误
 */
func ValidatePositiveInt(value int, fieldName string) error {
	if value <= 0 {
		return errors.New(fieldName + "必须是正整数")
	}
	return nil
}

/**
 * ValidateIntRange 验证整数范围
 * @param {int} value 待验证的整数
 * @param {string} fieldName 字段名称
 * @param {int} min 最小值
 * @param {int} max 最大值
 * @returns {error} 验证错误
 */
func ValidateIntRange(value int, fieldName string, min int, max int) error {
	if value < min || value > max {
		return fmt.Errorf("%s必须在%d到%d之间", fieldName, min, max)
	}
	return nil
}

/**
 * ValidateSliceNotEmpty 验证切片非空
 * @param {[]any} slice 待验证的切片
 * @param {string} fieldName 字段名称
 * @returns {error} 验证错误
 */
func ValidateSliceNotEmpty(slice []any, fieldName string) error {
	if len(slice) == 0 {
		return errors.New(fieldName + "不能为空")
	}
	return nil
}

/**
 * ValidateMapNotEmpty 验证映射非空
 * @param {map[string]any} m 待验证的映射
 * @param {string} fieldName 字段名称
 * @returns {error} 验证错误
 */
func ValidateMapNotEmpty(m map[string]any, fieldName string) error {
	if len(m) == 0 {
		return errors.New(fieldName + "不能为空")
	}
	return nil
}
