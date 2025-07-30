package wt

import (
	"time"
)

/**
 * BatchDeleteTokensByUserIDs 批量删除多个用户的所有token
 * @param {[]uint} userIDs 用户ID列表
 * @returns {ErrorCode} 操作结果错误码
 */
func (tm *Manager[T]) BatchDeleteTokensByUserIDs(userIDs []uint) ErrorCode {
	if len(userIDs) == 0 {
		return E_Success // 空列表被认为是成功的操作
	}

	// 验证用户ID
	for _, userID := range userIDs {
		if userID == 0 {
			return E_UserInvalid
		}
	}

	tm.lock()
	defer tm.unlock()

	// 创建用户ID集合用于快速查找
	userIDSet := make(map[uint]bool)
	for _, userID := range userIDs {
		userIDSet[userID] = true
	}

	// 统计每个用户删除的token数量
	deleteCount := make(map[uint]int)
	totalDeleted := 0
	expiredDeleted := 0
	activeDeleted := 0

	// 批量删除
	for token, ut := range tm.tokens {
		if userIDSet[ut.UserID] {
			// 检查token是否过期
			if ut.IsExpired() {
				expiredDeleted++
			} else {
				activeDeleted++
			}
			delete(tm.tokens, token)
			deleteCount[ut.UserID]++
			totalDeleted++
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

	return E_Success
}

/**
 * BatchDeleteTokensByGroupIDs 批量删除多个用户组的所有token
 * @param {[]uint} groupIDs 用户组ID列表
 * @returns {ErrorCode} 操作结果错误码
 */
func (tm *Manager[T]) BatchDeleteTokensByGroupIDs(groupIDs []uint) ErrorCode {
	if len(groupIDs) == 0 {
		return E_Success // 空列表被认为是成功的操作
	}

	// 验证用户组ID
	for _, groupID := range groupIDs {
		if groupID == 0 {
			return E_GroupInvalid
		}
	}

	tm.lock()
	defer tm.unlock()

	// 检查用户组是否存在
	for _, groupID := range groupIDs {
		if _, exists := tm.groups[groupID]; !exists {
			return E_GroupNotFound
		}
	}

	// 创建用户组ID集合用于快速查找
	groupIDSet := make(map[uint]bool)
	for _, groupID := range groupIDs {
		groupIDSet[groupID] = true
	}

	// 统计每个用户组删除的token数量
	deleteCount := make(map[uint]int)
	totalDeleted := 0
	expiredDeleted := 0
	activeDeleted := 0

	// 批量删除
	for token, ut := range tm.tokens {
		if groupIDSet[ut.GroupID] {
			// 检查token是否过期
			if ut.IsExpired() {
				expiredDeleted++
			} else {
				activeDeleted++
			}
			delete(tm.tokens, token)
			deleteCount[ut.GroupID]++
			totalDeleted++
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

	return E_Success
}

/**
 * BatchDeleteExpiredTokens 批量删除过期token
 * @returns {ErrorCode} 操作结果错误码
 */
func (tm *Manager[T]) BatchDeleteExpiredTokens() ErrorCode {
	tm.lock()
	defer tm.unlock()

	expiredTokens := make([]string, 0)

	// 收集过期token
	for key, token := range tm.tokens {
		if token == nil {
			expiredTokens = append(expiredTokens, key)
		} else if token.IsExpired() {
			expiredTokens = append(expiredTokens, key)
		}
	}

	// 批量删除
	for _, key := range expiredTokens {
		delete(tm.tokens, key)
	}

	deleteCount := len(expiredTokens)
	if deleteCount > 0 {
		// 删除过期token时，只减少总数，不减少过期token计数（过期token数是累计统计）
		tm.stats.TotalTokens -= deleteCount
		tm.stats.LastUpdateTime = time.Now()
	}

	return E_Success
}

/**
 * GetTokensByUserID 获取指定用户的所有token
 * @param {uint} userID 用户ID
 * @returns {[]*Token[T]} token列表
 */
func (tm *Manager[T]) GetTokensByUserID(userID uint) []*Token[T] {
	if userID == 0 {
		return nil
	}

	tm.rLock()
	defer tm.rUnlock()

	tokens := make([]*Token[T], 0)
	for _, token := range tm.tokens {
		if token.UserID == userID {
			tokens = append(tokens, token)
		}
	}

	return tokens
}

/**
 * GetTokensByGroupID 获取指定用户组的所有token
 * @param {uint} groupID 用户组ID
 * @returns {[]*Token[T]} token列表
 */
func (tm *Manager[T]) GetTokensByGroupID(groupID uint) []*Token[T] {
	if groupID == 0 {
		return nil
	}

	tm.rLock()
	defer tm.rUnlock()

	// 检查用户组是否存在
	if _, exists := tm.groups[groupID]; !exists {
		return nil
	}

	tokens := make([]*Token[T], 0)
	for _, token := range tm.tokens {
		if token.GroupID == groupID {
			tokens = append(tokens, token)
		}
	}

	return tokens
}
