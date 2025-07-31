package wt

import (
	"errors"
	"net"
	"strings"

	"github.com/windf17/wt/models"
)

/**
 * ValidateConfig 验证配置参数
 * @param {*ConfigRaw} config 配置对象
 * @returns {error} 验证错误
 */
func ValidateConfig(config models.ConfigRaw) error {
	// 验证语言设置
	if config.Language !="zh"{
		config.Language = "en"
	}

	// 定义以下报错信息的中英文双语字符串

	// 验证最大token数量
	if config.MaxTokens<0 {
		return errors.New("最大token数量不能小于0")
	}

	// 验证分隔符
	if config.Delimiter == "" {
		return errors.New("分隔符不能为空")
	}

	// 验证Token续期时间
	if err := validateTokenRenewTime(config.TokenRenewTime); err != nil {
		return err
	}

	return nil
}

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
	
	return nil
}

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
