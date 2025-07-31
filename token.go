package wt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/windf17/wt/models"
)

/**
 * GetToken 获取token数据
 * @param {string} key token键
 * @returns {*models.Token[T], error} token数据和错误信息
 */
func (tm *Manager[T]) GetToken(key string) (*models.Token[T], error) {
	// 性能监控

	// 先用读锁检查token是否存在
	tm.rLock()
	t := tm.tokens[key]
	if t == nil {
		tm.rUnlock()
		return nil, errors.New(getErrorMessage(tm.config.Language, "invalid_token"))
	}

	// 检查是否过期
	if t.IsExpired() {
		tm.rUnlock()

		// 使用单独的方法处理过期token删除，避免锁升级死锁
		tm.removeExpiredTokenSafe(key)
		return nil, errors.New(getErrorMessage(tm.config.Language, "token_expired"))
	}

	// 创建token副本，避免返回指针导致的并发问题
	tokenCopy := *t
	tm.rUnlock()

	// 更新最后访问时间（用于LRU策略）
	tm.lock()
	if currentToken, exists := tm.tokens[key]; exists && !currentToken.IsExpired() {
		currentToken.LastAccessTime = time.Now()
		tm.tokens[key] = currentToken
		// 更新副本中的访问时间
		tokenCopy.LastAccessTime = currentToken.LastAccessTime
	}
	tm.unlock()

	return &tokenCopy, nil
}

/**
 * AddToken 新增token，通过它申请token，不存储用户数据，存储用户数据另外用SetUserData
 * @param {uint} userID 用户ID
 * @param {uint} groupID 用户组ID
 * @param {string} clientIp 客户端IP地址
 * @returns {string, error} token字符串和错误信息
 */
func (tm *Manager[T]) AddToken(userID uint, groupID uint, clientIp string) (string, error) {
	if userID < 1 {
		return "", errors.New(getErrorMessage(tm.config.Language, "user_invalid"))
	}
	if groupID < 1 {
		return "", errors.New(getErrorMessage(tm.config.Language, "group_invalid"))
	}
	if err := ValidateIPAddress(clientIp); err != nil {
		return "", errors.New(getErrorMessage(tm.config.Language, "invalid_ip"))
	}

	// 首先检查用户组是否存在
	tm.rLock()
	g := tm.groups[groupID]
	if g == nil {
		tm.rUnlock()
		return "", errors.New(getErrorMessage(tm.config.Language, "group_not_found"))
	}
	tm.rUnlock()

	// 获取写锁进行token操作
	tm.lock()
	defer tm.unlock()

	// 如果不允许多设备登录，则清理该用户在其他设备上的token
	if !g.AllowMultipleLogin {
		expiredDeleted := 0
		activeDeleted := 0
		for t, ut := range tm.tokens {
			if ut.UserID == userID {
				// 检查token是否过期
				if ut.IsExpired() {
					expiredDeleted++
				} else {
					activeDeleted++
				}
				delete(tm.tokens, t)
			}
		}
		// 直接更新统计信息，避免重复加锁
		if activeDeleted > 0 {
			tm.stats.TotalTokens -= activeDeleted
			tm.stats.ActiveTokens -= activeDeleted
			tm.stats.LastUpdateTime = time.Now()
		}
		if expiredDeleted > 0 {
			// 对于过期token，只减少总数，不减少过期token计数
			tm.stats.TotalTokens -= expiredDeleted
			tm.stats.LastUpdateTime = time.Now()
		}
	}

	// 生成token
	tokenKey, er := tm.GenerateToken()
	if er != nil {
		return "", errors.New(getErrorMessage(tm.config.Language, "token_generate"))
	}

	// 创建用户tokens数据
	now := time.Now()
	var zero T
	tokenData := models.Token[T]{
		UserID:         userID,
		GroupID:        groupID,
		LoginTime:      now,
		LastAccessTime: now,
		ExpireSeconds:  g.ExpireSeconds,
		UserData:       zero,
		IP:             clientIp,
	}

	// 如果配置了最大token数量，先清理过期token
	if tm.config.MaxTokens > 0 {
		tm.cleanExpiredTokensInternal()

		// 检查清理后的token数量是否仍然达到上限
		if len(tm.tokens) >= tm.config.MaxTokens {
			// 清理最久没有使用的token（LRU策略）
			tm.cleanOldestTokensInternal(1)
		}
	}

	// 存储token
	tm.tokens[tokenKey] = &tokenData
	// 直接更新统计信息，避免重复加锁
	tm.stats.TotalTokens += 1
	tm.stats.ActiveTokens += 1
	tm.stats.LastUpdateTime = time.Now()

	return tokenKey, nil
}

/**
 * DelToken 删除指定的token
 * @param {string} key token键
 * @returns {error} 操作结果错误信息
 */
func (tm *Manager[T]) DelToken(key string) error {
	tm.lock()
	defer tm.unlock()
	token, exists := tm.tokens[key]
	if !exists {
		return errors.New(getErrorMessage(tm.config.Language, "invalid_token"))
	}

	// 检查token是否过期
	isExpired := token.IsExpired()
	delete(tm.tokens, key)
	// 直接更新统计信息，避免重复加锁
	if isExpired {
		// 对于过期token，只减少总数
		tm.stats.TotalTokens -= 1
		tm.stats.LastUpdateTime = time.Now()
	} else {
		// 减少总数和活跃数
		tm.stats.TotalTokens -= 1
		tm.stats.ActiveTokens -= 1
		tm.stats.LastUpdateTime = time.Now()
	}

	return nil
}

/**
 * DelTokensByUserID 删除指定用户的所有token
 * @param {uint} userID 用户ID
 * @returns {error} 操作结果错误信息
 */
func (tm *Manager[T]) DelTokensByUserID(userID uint) error {
	if userID == 0 {
		return errors.New(getErrorMessage(tm.config.Language, "user_invalid"))
	}
	tm.lock()
	defer tm.unlock()
	expiredDeleted := 0
	activeDeleted := 0
	for token, ut := range tm.tokens {
		if ut.UserID == userID {
			// 检查token是否过期
			if ut.IsExpired() {
				expiredDeleted++
			} else {
				activeDeleted++
			}
			delete(tm.tokens, token)
		}
	}
	// 直接更新统计信息，避免重复加锁
	if activeDeleted > 0 {
		tm.stats.TotalTokens -= activeDeleted
		tm.stats.ActiveTokens -= activeDeleted
		tm.stats.LastUpdateTime = time.Now()
	}
	if expiredDeleted > 0 {
		// 对于过期token，只减少总数，不减少过期token计数
		tm.stats.TotalTokens -= expiredDeleted
		tm.stats.LastUpdateTime = time.Now()
	}
	return nil
}

/**
 * DelTokensByGroupID 删除指定用户组的所有token
 * @param {uint} groupID 用户组ID
 * @returns {error} 操作结果错误信息
 */
func (tm *Manager[T]) DelTokensByGroupID(groupID uint) error {
	if groupID == 0 {
		return errors.New(getErrorMessage(tm.config.Language, "group_invalid"))
	}
	tm.lock()
	defer tm.unlock()
	// 检查用户组id是不是存在
	if _, exists := tm.groups[groupID]; !exists {
		return errors.New(getErrorMessage(tm.config.Language, "group_not_found"))
	}
	expiredDeleted := 0
	activeDeleted := 0
	for token, ut := range tm.tokens {
		if ut.GroupID == groupID {
			// 检查token是否过期
			if ut.IsExpired() {
				expiredDeleted++
			} else {
				activeDeleted++
			}
			delete(tm.tokens, token)
		}
	}
	// 直接更新统计信息，避免重复加锁
	if activeDeleted > 0 {
		tm.stats.TotalTokens -= activeDeleted
		tm.stats.ActiveTokens -= activeDeleted
		tm.stats.LastUpdateTime = time.Now()
	}
	if expiredDeleted > 0 {
		// 对于过期token，只减少总数，不减少过期token计数
		tm.stats.TotalTokens -= expiredDeleted
		tm.stats.LastUpdateTime = time.Now()
	}
	return nil
}

// UpdateToken 更新指定的token
func (tm *Manager[T]) UpdateToken(key string, token *models.Token[T]) error {
	tm.lock()
	defer tm.unlock()
	if _, exists := tm.tokens[key]; !exists {
		return errors.New(getErrorMessage(tm.config.Language, "invalid_token"))
	}
	if token == nil {
		return errors.New(getErrorMessage(tm.config.Language, "invalid_token"))
	}
	token.LastAccessTime = time.Now()
	tm.tokens[key] = token

	return nil
}

/**
 * CleanExpiredTokens 清理过期token并更新缓存文件
 */
func (tm *Manager[T]) CleanExpiredTokens() {
	tm.lock()
	defer tm.unlock()
	tm.cleanExpiredTokensInternal()
}

/**
 * cleanExpiredTokensInternal 内部清理过期token函数（不获取锁）
 */
func (tm *Manager[T]) cleanExpiredTokensInternal() {
	expiredCount := 0
	nullCount := 0
	for key, token := range tm.tokens {
		if token == nil {
			delete(tm.tokens, key)
			nullCount++
		} else if token.IsExpired() {
			delete(tm.tokens, key)
			expiredCount++
		}
	}

	// 直接更新统计信息，避免重复加锁
	if expiredCount > 0 {
		tm.stats.TotalTokens -= expiredCount
		tm.stats.ExpiredTokens += expiredCount
		tm.stats.LastUpdateTime = time.Now()
	}
	if nullCount > 0 {
		tm.stats.TotalTokens -= nullCount
		tm.stats.ActiveTokens -= nullCount
		tm.stats.LastUpdateTime = time.Now()
	}
}

/**
 * GenerateToken 生成随机token
 * @returns {string, error} token字符串和错误
 */
func (tm *Manager[T]) GenerateToken() (string, error) {
	b := make([]byte, TOKEN_BYTE_SIZE)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	now := time.Now().UnixNano()
	nowBytes := make([]byte, TIMESTAMP_BYTE_SIZE)
	for i := 0; i < TIMESTAMP_BYTE_SIZE; i++ {
		nowBytes[i] = byte(now >> uint(i*8))
	}
	b = append(nowBytes, b...)
	return base64.URLEncoding.EncodeToString(b), nil
}

/**
 * removeExpiredTokenSafe 安全删除过期token，避免死锁
 * @param {string} key token键
 */
func (tm *Manager[T]) removeExpiredTokenSafe(key string) {
	tm.lock()
	defer tm.unlock()

	// 重新检查，因为在锁切换期间可能有变化
	if t, exists := tm.tokens[key]; exists && t != nil && t.IsExpired() {
		delete(tm.tokens, key)
		// 删除过期token时，只减少总数，不减少过期token计数（过期token数是累计统计）
		tm.stats.TotalTokens -= 1
		tm.stats.LastUpdateTime = time.Now()
	}
}

/**
 * cleanOldestTokensInternal 清理最久没有使用的token（LRU策略）
 * @param {int} count 要清理的token数量
 */
func (tm *Manager[T]) cleanOldestTokensInternal(count int) {
	if count <= 0 || len(tm.tokens) == 0 {
		return
	}

	// 创建token切片用于排序
	type tokenInfo struct {
		key        string
		lastAccess time.Time
	}

	tokensToSort := make([]tokenInfo, 0, len(tm.tokens))
	for key, token := range tm.tokens {
		tokensToSort = append(tokensToSort, tokenInfo{
			key:        key,
			lastAccess: token.LastAccessTime,
		})
	}

	// 按LastAccessTime升序排序（最久的在前面）
	for i := 0; i < len(tokensToSort)-1; i++ {
		for j := i + 1; j < len(tokensToSort); j++ {
			if tokensToSort[i].lastAccess.After(tokensToSort[j].lastAccess) {
				tokensToSort[i], tokensToSort[j] = tokensToSort[j], tokensToSort[i]
			}
		}
	}

	// 删除最久没有使用的token
	deleteCount := min(count, len(tokensToSort))
	expiredDeleted := 0
	activeDeleted := 0

	for i := 0; i < deleteCount; i++ {
		key := tokensToSort[i].key
		// 检查token是否过期，以便正确更新统计信息
		if token, exists := tm.tokens[key]; exists {
			isExpired := token.IsExpired()
			delete(tm.tokens, key)
			if isExpired {
				expiredDeleted++
			} else {
				activeDeleted++
			}
		}
	}

	// 直接更新统计信息，避免重复加锁
	if expiredDeleted > 0 {
		tm.stats.TotalTokens -= expiredDeleted
		tm.stats.ExpiredTokens += expiredDeleted
		tm.stats.LastUpdateTime = time.Now()
	}
	if activeDeleted > 0 {
		tm.stats.TotalTokens -= activeDeleted
		tm.stats.ActiveTokens -= activeDeleted
		tm.stats.LastUpdateTime = time.Now()
	}
}
