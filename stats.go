package wt

import (
	"time"

	"github.com/windf17/wt/models"
)

// GetStats 获取token统计信息
func (tm *Manager[T]) GetStats() models.Stats {
	tm.rLock()
	defer tm.rUnlock()
	statsCopy := models.Stats{
		LastUpdateTime: tm.stats.LastUpdateTime,
		TotalTokens:    tm.stats.TotalTokens,
		ActiveTokens:   tm.stats.ActiveTokens,
		ExpiredTokens:  tm.stats.ExpiredTokens,
	}
	return statsCopy
}

// updateStatsCount 更新token统计数量
// count为正数时增加统计数，为负数时减少统计数
// isExpired参数用于指定是否为过期token的统计
func (tm *Manager[T]) updateStatsCount(count int, isExpired bool) {
	tm.lock()
	defer tm.unlock()
	tm.stats.TotalTokens += count
	if isExpired {
		tm.stats.ExpiredTokens += count
	} else {
		tm.stats.ActiveTokens += count
	}
	tm.stats.LastUpdateTime = time.Now()
}
