package models

// ApiRule 定义API规则
type ApiRule struct {
	// 路径
	Path []string `json:"path"`
	// 规则：true表示允许，false表示禁止
	Rule bool `json:"rule"`
}
// Group 用户组配置
type Group struct {
	// 名称
	Name string `json:"name"`
	// api权限规则
	ApiRules []ApiRule `json:"apiRules"`
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
