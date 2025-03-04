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
		return ErrInvalidToken
	}
	if url == "" {
		return ErrInvalidURL
	}
	if ip == "" {
		return ErrInvalidIP
	}

	t, exists := tm.tokens[key]
	if !exists {
		return ErrTokenNotFound
	}
	if t == nil {
		// token的key存在但值为nil，则删除该token
		delete(tm.tokens, key)
		return ErrInvalidToken
	}
	if t.IsExpired() {
		// token已过期，则删除该token
		delete(tm.tokens, key)
		return ErrTokenExpired
	}

	// 获取用户组配置
	g := tm.groups[t.GroupID]
	if g == nil {
		return ErrGroupNotFound
	}

	// 如果当前用户组不允许同一IP多次登录，则检查当前IP与申请token时的ip是否一致
	if !g.AllowMultipleLogin && t.IP != ip {
		return (ErrInvalidIP)
	}

	// 处理url，获取api路径
	apiPath := utility.NormalizeAPIPath(url)
	totalLen := len(g.AllowedAPIs) + len(g.DeniedAPIs)
	// 如果用户组没有允许的API路径和拒绝的API路径，则拒绝访问
	if totalLen == 0 {
		return ErrNoAPIPermission
	}
	// 检查API路径权限
	if !utility.HasPermission(apiPath, g.AllowedAPIs, g.DeniedAPIs) {
		return ErrAccessDenied
	}

	// 更新最后访问时间
	t.LastAccessTime = time.Now()
	if t.ExpireSeconds > 0 {
		t.ExpireSeconds += tm.config.TokenRenewTime
	}
	tm.tokens[key] = t

	return ErrSuccess
}
