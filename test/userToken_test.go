package test

import (
	"testing"
	"time"

	"github.com/windf17/wtoken"
)

func TestToken_IsExpired(t *testing.T) {
	// 设置基准时间为10分钟前
	tenMinutesAgo := time.Now().Add(-10 * time.Minute)

	tests := []struct {
		name        string
		token       *wtoken.Token[any]
		wantExpired bool
	}{
		{
			name: "永不过期的token",
			token: &wtoken.Token[any]{
				UserID:     1,
				LoginTime:  tenMinutesAgo,
				ExpireTime: 0, // 0表示永不过期
			},
			wantExpired: false,
		},
		{
			name: "已过期的token(10分钟前登录，有效期5分钟)",
			token: &wtoken.Token[any]{
				UserID:     2,
				LoginTime:  tenMinutesAgo,
				ExpireTime: 300, // 300秒 = 5分钟
			},
			wantExpired: true,
		},
		{
			name: "未过期的token(10分钟前登录，有效期11分钟)",
			token: &wtoken.Token[any]{
				UserID:     3,
				LoginTime:  tenMinutesAgo,
				ExpireTime: 660, // 660秒 = 11分钟
			},
			wantExpired: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isExpired := tt.token.IsExpired()
			if isExpired != tt.wantExpired {
				t.Errorf("Token.IsExpired() = %v, 期望 %v", isExpired, tt.wantExpired)
			}
		})
	}
}
