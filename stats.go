package wtoken

import (
	"time"
)

// GetStats 获取token统计信息
func (tm *Manager[T]) GetStats() Stats {
	tm.rLock()
	defer tm.rUnlock()
	statsCopy := Stats{
		LastUpdateTime: tm.stats.LastUpdateTime,
		TotalTokens:    tm.stats.TotalTokens,
		ActiveTokens:   tm.stats.ActiveTokens,
		ExpiredTokens:  tm.stats.ExpiredTokens,
	}
	return statsCopy
}

// UpdateStats 更新统计信息
func (tm *Manager[T]) UpdateStats() {
	activeTokens := 0
	expiredTokens := 0

	// 只统计token状态，不做清理
	for _, ut := range tm.tokens {
		if ut.IsExpired() {
			expiredTokens++
		} else {
			activeTokens++
		}
	}
	// 更新统计信息
	tm.stats.TotalTokens = len(tm.tokens)
	tm.stats.ActiveTokens = activeTokens
	tm.stats.ExpiredTokens = expiredTokens
	tm.stats.LastUpdateTime = time.Now()
}
