package wtoken

import (
	"time"

	"github.com/windf17/wtoken/utility"
)

// Authenticate 鉴权，key即为token，url为请求的API地址，ip为客户端IP地址
func (tm *Manager[T]) Authenticate(key string, url string, ip string) ErrorCode {
	tm.rLock()
	defer tm.rUnlock()

	if key == "" {
		return E_InvalidToken
	}
	if url == "" {
		return E_APINotFound
	}
	if ip == "" {
		return E_InvalidIP
	}

	t, exists := tm.tokens[key]
	if !exists {
		return E_InvalidToken
	}
	if t == nil {
		// token的key存在但值为nil，则删除该token
		delete(tm.tokens, key)
		return E_InvalidToken
	}
	if t.IsExpired() {
		// token已过期，则删除该token
		delete(tm.tokens, key)
		return E_TokenExpired
	}

	// 获取用户组配置
	g := tm.groups[t.GroupID]
	if g == nil {
		return E_GroupNotFound
	}

	// 如果当前用户组不允许同一IP多次登录，则检查当前IP与申请token时的ip是否一致
	if !g.AllowMultipleLogin && t.IP != ip {
		return (E_IPMismatch)
	}

	// 处理url，获取api路径
	apiPath := utility.NormalizeAPIPath(url)
	totalLen := len(g.AllowedAPIs) + len(g.DeniedAPIs)
	// 如果用户组没有允许的API路径和拒绝的API路径，则拒绝访问
	if totalLen == 0 {
		return E_Forbidden
	}
	// 检查API路径权限
	if !utility.HasPermission(apiPath, g.AllowedAPIs, g.DeniedAPIs) {
		return E_Forbidden
	}

	// 更新最后访问时间
	t.LastAccessTime = time.Now()
	if t.ExpireSeconds > 0 {
		t.ExpireSeconds += tm.config.TokenRenewTime
	}
	tm.tokens[key] = t

	return E_Success
}
