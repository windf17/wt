package wt

import (
	"strings"
	"time"

	"github.com/windf17/wt/utility"
)

/**
 * Auth 专门负责对客户端访问指定API进行鉴权
 * 包含完整的鉴权流程：Token验证、IP验证和API权限验证
 * @param {string} key token字符串
 * @param {string} clientIp 客户端IP地址
 * @param {string} api 请求的API地址
 * @returns {ErrorCode} 鉴权结果
 */
func (tm *Manager[T]) Auth(key string, clientIp string, api string) ErrorCode {
	// 输入参数验证
	if strings.TrimSpace(key) == "" {
		return E_InvalidToken // 无效Token
	}
	if clientIp == "" {
		return E_InvalidIP // 无效IP
	}
	// 检查是否启用了鉴权功能（如果没有配置任何用户组，则禁用鉴权）
	tm.rLock()
	if len(tm.groups) == 0 {
		tm.rUnlock()
		// 没有配置用户组，跳过权限验证，直接返回成功
		return E_Success
	}

	// 第一阶段：Token验证（防止盗用）
	t, exists := tm.tokens[key]
	if !exists {
		tm.rUnlock()
		return E_InvalidToken // 无效Token
	}
	if t == nil {
		// token的key存在但值为nil，需要删除该token
		tm.rUnlock()
		// 升级为写锁进行删除操作
		tm.lock()
		// 重新检查token是否仍然存在且为nil
		if currentToken, exists := tm.tokens[key]; exists && currentToken == nil {
			delete(tm.tokens, key)
		}
		tm.unlock()
		return E_InvalidToken // 无效Token
	}
	if t.IsExpired() {
		// token已过期，需要删除该token
		tm.rUnlock()
		// 升级为写锁进行删除操作
		tm.lock()
		// 重新检查token是否仍然过期
		if currentToken, exists := tm.tokens[key]; exists && currentToken != nil && currentToken.IsExpired() {
			delete(tm.tokens, key)
		}
		tm.unlock()
		return E_TokenExpired // Token过期，拒绝访问
	}

	// IP验证：若token有效但是ip不匹配，则判断为token被盗用
	if t.IP != clientIp {
		tm.rUnlock()
		return E_Forbidden // IP不匹配，token被盗用，禁止访问
	}

	// 获取用户组配置
	g := tm.groups[t.GroupID]
	if g == nil {
		tm.rUnlock()
		return E_Forbidden // 用户组不存在，拒绝访问
	}

	// 第二阶段：API权限验证
	// 如果用户组没有配置任何API规则，则拒绝访问
	if len(g.ApiRules) == 0 {
		tm.rUnlock()
		return E_Unauthorized // 无权访问
	}
	// 检查API路径权限
	if !utility.HasPermission(api, g.ApiRules) {
		tm.rUnlock()
		return E_Unauthorized // 无权访问
	}

	// 第三阶段：更新最后访问时间
	// 释放读锁，升级为写锁
	tm.rUnlock()
	tm.lock()

	// 重新验证token是否仍然有效（防止在锁切换期间token被删除）
	if currentToken, exists := tm.tokens[key]; exists && currentToken != nil && !currentToken.IsExpired() {
		// 更新最后访问时间
		currentToken.LastAccessTime = time.Now()
		tm.tokens[key] = currentToken
		tm.unlock()
		return E_Success
	} else {
		// Token在锁切换期间被删除或过期
		tm.unlock()
		return E_Forbidden // Token无效，拒绝访问
	}
}

/**
 * BatchAuth 批量API权限检查
 * 用于前端一次性检查多个API的访问权限
 * @param {string} key token字符串
 * @param {string} clientIp 客户端IP地址
 * @param {[]string} apis 需要检查的API地址数组
 * @returns {[]bool} 对应每个API的权限检查结果数组，true表示有权限，false表示无权限
 */
func (tm *Manager[T]) BatchAuth(key string, clientIp string, apis []string) []bool {
	// 初始化结果数组
	results := make([]bool, len(apis))

	// 对每个API进行权限检查
	for i, api := range apis {
		// 调用单个API的权限检查方法
		authResult := tm.Auth(key, clientIp, api)
		// 只有返回E_Success时才表示有权限
		results[i] = (authResult == E_Success)
	}

	return results
}
