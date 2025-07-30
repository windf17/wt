package wt

import "time"

// Stats token统计信息
type Stats struct {
	// TotalTokens token总数
	TotalTokens int `json:"total_tokens"`
	// ActiveTokens 活跃token数量
	ActiveTokens int `json:"active_tokens"`
	// ExpiredTokens 过期token数量
	ExpiredTokens int `json:"expired_tokens"`
	// LastUpdateTime 最后更新时间
	LastUpdateTime time.Time `json:"last_update_time"`
}
