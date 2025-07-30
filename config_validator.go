package wtoken

import (
	"errors"
	"strings"

	"github.com/windf17/wtoken/models"
)

/**
 * ValidateConfig 验证配置参数
 * @param {*ConfigRaw} config 配置对象
 * @returns {error} 验证错误
 */
func ValidateConfig(config *ConfigRaw) error {
	if config == nil {
		return errors.New("配置不能为空")
	}

	// 验证最大token数量
	if err := validateMaxTokens(config.MaxTokens); err != nil {
		return err
	}

	// 验证分隔符
	if err := validateDelimiter(config.Delimiter); err != nil {
		return err
	}

	// 验证Token续期时间
	if err := validateTokenRenewTime(config.TokenRenewTime); err != nil {
		return err
	}

	// 验证语言设置
	if err := validateLanguage(config.Language); err != nil {
		return err
	}

	return nil
}



/**
 * validateMaxTokens 验证最大token数量
 * @param {int} maxTokens 最大token数量
 * @returns {error} 验证错误
 */
func validateMaxTokens(maxTokens int) error {
	const (
		MIN_MAX_TOKENS = 1
		MAX_MAX_TOKENS = 1000000
	)
	if err := ValidateIntRange(maxTokens, "最大token数量", MIN_MAX_TOKENS, MAX_MAX_TOKENS); err != nil {
		return err
	}
	return nil
}

/**
 * validateDelimiter 验证分隔符
 * @param {string} delimiter 分隔符
 * @returns {error} 验证错误
 */
func validateDelimiter(delimiter string) error {
	if delimiter == "" {
		return errors.New("分隔符不能为空")
	}
	return nil
}



// ValidateIntRange函数已在validation.go中定义，此处移除重复定义

/**
 * validateTokenRenewTime 验证Token续期时间格式
 * @param {string} renewTime 续期时间字符串
 * @returns {error} 验证错误
 */
func validateTokenRenewTime(renewTime string) error {
	if len(renewTime) < 2 {
		return errors.New("TokenRenewTime格式错误，至少需要2个字符")
	}

	// 获取单位
	unit := renewTime[len(renewTime)-1]
	valueStr := renewTime[:len(renewTime)-1]

	// 验证数值部分
	if valueStr == "" {
		return errors.New("TokenRenewTime缺少数值部分")
	}

	// 验证单位
	switch unit {
	case 's', 'S', 'm', 'M', 'h', 'H', 'd', 'D':
		// 合法单位
	default:
		return errors.New("TokenRenewTime单位不支持，支持的单位：s(秒)、m(分钟)、h(小时)、d(天)")
	}

	return nil
}

/**
 * validateLanguage 验证语言设置
 * @param {Language} lang 语言类型
 * @returns {error} 验证错误
 */
func validateLanguage(lang Language) error {
	// 允许任何语言类型，包括自定义语言
	if lang == "" {
		return nil // 空值使用默认语言
	}
	// 验证语言代码格式（2-5个字符的字母）
	langStr := string(lang)
	if len(langStr) < 2 || len(langStr) > 5 {
		return errors.New("语言代码长度应为2-5个字符")
	}
	for _, char := range langStr {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') {
			return errors.New("语言代码只能包含字母")
		}
	}
	return nil
}



/**
 * ValidateGroupRaw 验证用户组配置
 * @param {GroupRaw} group 用户组配置
 * @returns {error} 验证错误
 */
func ValidateGroupRaw(group models.GroupRaw) error {
	if group.ID == 0 {
		return errors.New("用户组ID不能为0")
	}

	if strings.TrimSpace(group.Name) == "" {
		return errors.New("用户组名称不能为空")
	}

	// 验证Token过期时间
	if group.TokenExpire != "" {
		if err := validateTokenRenewTime(group.TokenExpire); err != nil {
			return errors.New("用户组TokenExpire格式错误: " + err.Error())
		}
	}

	// 验证多设备登录设置
	if group.AllowMultipleLogin != 0 && group.AllowMultipleLogin != 1 {
		return errors.New("AllowMultipleLogin只能为0或1")
	}

	return nil
}