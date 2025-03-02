package wtoken

import (
	"time"

	"github.com/windf17/wtoken/utility"
)

// Authenticate 鉴权，key即为token，url为请求的API地址，ip为客户端IP地址
func (tm *Manager[T]) Authenticate(key string, url string, ip string) ErrorData {
	tm.rLock()
	defer tm.rUnlock()

	if key == "" {
		return tm.NewError(ErrCodeInvalidToken)
	}
	if url == "" {
		return tm.NewError(ErrCodeInvalidURL)
	}
	if ip == "" {
		return tm.NewError(ErrCodeInvalidIP)
	}

	t, exists := tm.tokens[key]
	if !exists {
		return tm.NewError(ErrCodeTokenNotFound)
	}
	if t == nil {
		// token的key存在但值为nil，则删除该token
		delete(tm.tokens, key)
		return tm.NewError(ErrCodeInvalidToken)
	}
	if t.IsExpired() {
		// token已过期，则删除该token
		delete(tm.tokens, key)
		return tm.NewError(ErrCodeTokenExpired)
	}

	// 获取用户组配置
	g := tm.groups[t.GroupID]
	if g == nil {
		return tm.NewError(ErrCodeGroupNotFound)
	}

	// 如果当前用户组不允许同一IP多次登录，则检查当前IP与申请token时的ip是否一致
	if !g.AllowMultipleLogin && t.IP != ip {
		return tm.NewError(ErrCodeInvalidIP)
	}

	// 处理url，获取api路径
	apiPath := utility.NormalizeAPIPath(url)

	// 检查API路径权限
	if !utility.HasPermission(apiPath, g.AllowedAPIs, g.DeniedAPIs) {
		return tm.NewError(ErrCodeAccessDenied)
	}

	// 更新最后访问时间
	t.LastAccessTime = time.Now()
	if t.ExpireTime > 0 && g.ExpireSeconds > 0 {
		t.ExpireTime += g.ExpireSeconds
	}
	tm.tokens[key] = t

	return tm.NewError(ErrCodeSuccess)
}
