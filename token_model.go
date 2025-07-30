package wt

import (
	"time"
)

// Token 用户Token信息
type Token[T any] struct {
	// 用户ID
	UserID uint `json:"userId"`
	// 用户组ID
	GroupID uint `json:"groupId"`
	// 登录时间
	LoginTime time.Time `json:"loginTime"`
	// 过期秒数，为0表示永不过期，大于0表示从登录时间起多少秒后过期，它会在使用token时刷新
	ExpireSeconds int64 `json:"expireTime"`
	// 最后访问时间
	LastAccessTime time.Time `json:"lastAccessTime"`
	// 用户数据
	UserData T `json:"userData"`
	// Token所属用户的IP地址
	IP string `json:"ip"`
}

// IsExpired 检查token是否过期
func (ut *Token[T]) IsExpired() bool {
	if ut.ExpireSeconds == 0 {
		return false // 如果过期秒数为0，则永不过期
	}
	expirationTime := ut.LoginTime.Add(time.Duration(ut.ExpireSeconds) * time.Second)
	return time.Now().After(expirationTime)
}
