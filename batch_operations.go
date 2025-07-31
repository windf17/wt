package wt

import (
	"errors"
	"time"

	"github.com/windf17/wt/models"
)

/**
 * BatchDeleteTokensByUserIDs 批量删除多个用户的所有token
 * @param {[]uint} userIDs 用户ID列表
 * @returns {error} 操作结果错误信息
 */
func (tm *Manager[T]) BatchDeleteTokensByUserIDs(userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil // 空列表被认为是成功的操作
	}

	// 验证用户ID
	for _, userID := range userIDs {
		if userID == 0 {
			return errors.New(getErrorMessage(tm.config.Language, "user_invalid"))
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

	// 原子性更新统计信息
	if totalDeleted > 0 {
		tm.stats.TotalTokens -= totalDeleted
		if activeDeleted > 0 {
			tm.stats.ActiveTokens -= activeDeleted
		}
		tm.stats.LastUpdateTime = time.Now()
	}

	return nil
}

/**
 * BatchDeleteTokensByGroupIDs 批量删除多个用户组的所有token
 * @param {[]uint} groupIDs 用户组ID列表
 * @returns {error} 操作结果错误信息
 */
func (tm *Manager[T]) BatchDeleteTokensByGroupIDs(groupIDs []uint) error {
	if len(groupIDs) == 0 {
		return nil // 空列表被认为是成功的操作
	}

	// 验证用户组ID
	for _, groupID := range groupIDs {
		if groupID == 0 {
			return errors.New(getErrorMessage(tm.config.Language, "group_invalid"))
		}
	}

	tm.lock()
	defer tm.unlock()

	// 检查用户组是否存在
	for _, groupID := range groupIDs {
		if _, exists := tm.groups[groupID]; !exists {
			return errors.New(getErrorMessage(tm.config.Language, "group_not_found"))
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

	return nil
}

/**
 * BatchDeleteExpiredTokens 批量删除过期token
 * @returns {error} 操作结果错误信息
 */
func (tm *Manager[T]) BatchDeleteExpiredTokens() error {
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

	return nil
}

/**
 * GetTokensByUserID 获取指定用户的所有token
 * @param {uint} userID 用户ID
 * @returns {[]*models.Token[T]} token列表
 */
func (tm *Manager[T]) GetTokensByUserID(userID uint) []*models.Token[T] {
	if userID == 0 {
		return nil
	}

	tm.rLock()
	defer tm.rUnlock()

	tokens := make([]*models.Token[T], 0)
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
 * @returns {[]*models.Token[T]} token列表
 */
func (tm *Manager[T]) GetTokensByGroupID(groupID uint) []*models.Token[T] {
	if groupID == 0 {
		return nil
	}

	tm.rLock()
	defer tm.rUnlock()

	// 检查用户组是否存在
	if _, exists := tm.groups[groupID]; !exists {
		return nil
	}

	tokens := make([]*models.Token[T], 0)
	for _, token := range tm.tokens {
		if token.GroupID == groupID {
			tokens = append(tokens, token)
		}
	}

	return tokens
}
