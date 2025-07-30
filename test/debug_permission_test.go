package test

import (
	"testing"

	"github.com/windf17/wtoken"
	"github.com/windf17/wtoken/models"
)

/**
 * TestDebugPermission 调试权限验证问题
 */
func TestDebugPermission(t *testing.T) {
	// 初始化配置
	config := &wtoken.ConfigRaw{

		MaxTokens:      10,
		Delimiter:      ",",
		TokenRenewTime: "30m",
		Language:       wtoken.LangChinese,
	}

	groups := []models.GroupRaw{
		{
			ID:                 2,
			Name:               "user",
			AllowedAPIs:        "/api/user/profile",
			DeniedAPIs:         "/api/user/admin",
			TokenExpire:        "1h",
			AllowMultipleLogin: 1,
		},
	}

	// 初始化token管理器
	tm := wtoken.InitTM[map[string]any](config, groups, nil)
	if tm == nil {
		t.Fatalf("Failed to initialize token manager")
	}
	defer tm.Close()

	// 创建用户token
	userToken, errCode := tm.AddToken(2, 2, "192.168.1.102")
	if errCode != wtoken.E_Success {
		t.Errorf("Failed to add user token: %v", errCode)
	}

	t.Logf("Created token: %s", userToken)

	// 获取用户组信息
	group, errCode := tm.GetGroup(2)
	if errCode != wtoken.E_Success {
		t.Errorf("Failed to get group: %v", errCode)
		return
	}

	t.Logf("Group API Rules:")
	for i, rule := range group.ApiRules {
		t.Logf("  Rule %d: %v (rule: %v)", i, rule.Path, rule.Rule)
	}

	// 测试API权限
	testCases := []struct {
		api      string
		expected wtoken.ErrorCode
		desc     string
	}{
		{"/api/user/profile", wtoken.E_Success, "exact match /api/user/profile"},
		{"/api/user/admin", wtoken.E_Unauthorized, "denied /api/user/admin"},
		{"/api/admin", wtoken.E_Unauthorized, "not allowed /api/admin"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			errCode := tm.Auth(userToken, "192.168.1.102", tc.api)
			t.Logf("API: %s, Expected: %v, Got: %v", tc.api, tc.expected, errCode)
			if errCode != tc.expected {
				t.Errorf("API %s: expected %v, got %v", tc.api, tc.expected, errCode)
			}
		})
	}
}
