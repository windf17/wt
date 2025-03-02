package wtoken

import (
	"strconv"

	"github.com/windf17/wtoken/utility"
)

// Group 用户组配置
type Group struct {
	// 名称
	Name string `json:"name"`
	// 允许访问的API列表
	AllowedAPIs []string `json:"allowedApis"`
	// 禁止访问的API列表
	DeniedAPIs []string `json:"deniedApis"`
	// Token过期时间（秒），0表示永不过期
	ExpireSeconds int64 `json:"tokenExpireSeconds"`
	// 允许多设备登录。为true时允许同一用户在多个设备上登录，不校验IP；为false时只允许在一个设备上登录，会校验IP
	AllowMultipleLogin bool `json:"allowMultipleLogin"`
}

// 用户组原型
type GroupRaw struct {
	// 组ID
	ID uint `json:"id"`
	// 组名称
	Name string `json:"name"`
	// 允许访问的API列表
	AllowedAPIs string `json:"allowedApis"`
	// 禁止访问的API列表
	DeniedAPIs string `json:"deniedApis"`
	// Token过期时间（秒），0表示永不过期
	TokenExpire string `json:"tokenExpire"`
	// 允许多设备登录。为1时允许同一用户在多个设备上登录，不校验IP；为false时只允许在一个设备上登录，会校验IP
	AllowMultipleLogin int `json:"allowMultipleLogin"`
}

// ContGroup 将GroupRaw转换为Group
func ConvGroup(raw GroupRaw, delimiter string) *Group {
	g := Group{}

	// 处理 AllowMultipleLogin
	if raw.AllowMultipleLogin == 1 {
		g.AllowMultipleLogin = true
	} else {
		g.AllowMultipleLogin = false
	}
	// 处理 Name
	g.Name = raw.Name
	// 处理 TokenExpire
	g.ExpireSeconds = parseTokenExpire(raw.TokenExpire)
	// 处理 AllowedAPIs
	g.AllowedAPIs = utility.ParseAPIs(raw.AllowedAPIs, delimiter)
	// 处理 DeniedAPIs
	g.DeniedAPIs = utility.ParseAPIs(raw.DeniedAPIs, delimiter)
	return &g
}

// parseTokenExpire 将TokenExpire字符串转换为秒数
func parseTokenExpire(tokenExpire string) int64 {
	if tokenExpire == "" {
		return 0
	}

	// 获取单位
	unit := tokenExpire[len(tokenExpire)-1]
	valueStr := tokenExpire[:len(tokenExpire)-1]

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0
	}

	switch unit {
	case 'h', 'H': // 时间单位是小时
		return int64(value * 3600)
	case 'm', 'M': // 时间单位是分钟
		return int64(value * 60)
	case 'd', 'D': // 时间单位是天
		return int64(value * 86400)
	default:
		return int64(value)
	}
}
