package test

import (
	"errors"
	"testing"

	"github.com/windf17/wt"
	"github.com/windf17/wt/models"
)

/**
 * TestDebugPermission 调试权限验证问题
 */
func TestDebugPermission(t *testing.T) {
	// 初始化配置
	config := &models.ConfigRaw{

		MaxTokens:      10,
		Delimiter:      ",",
		TokenRenewTime: "30m",
		Language:       "zh",
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
	tm, err := wt.InitTM[map[string]any](*config, groups)
	if err != nil {
		t.Fatalf("Failed to initialize token manager: %v", err)
	}
	if tm == nil {
		t.Fatalf("Failed to initialize token manager")
	}


	// 创建用户token
	userToken, err := tm.AddToken(2, 2, "192.168.1.102")
	if err != nil {
		t.Errorf("Failed to add user token: %v", err)
	}

	t.Logf("Created token: %s", userToken)

	// 获取用户组信息
	group, err := tm.GetGroup(2)
	if err != nil {
		t.Errorf("Failed to get group: %v", err)
		return
	}

	t.Logf("Group API Rules:")
	for i, rule := range group.ApiRules {
		t.Logf("  Rule %d: %v (rule: %v)", i, rule.Path, rule.Rule)
	}

	// 测试API权限
	testCases := []struct {
		api      string
		expected error
		desc     string
	}{
		{"/api/user/profile", nil, "exact match /api/user/profile"},
		{"api/user/admin", errors.New("未授权访问"), "denied /api/user/admin"},
		{"/api/admin", errors.New("未授权访问"), "not allowed /api/admin"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := tm.Auth(userToken, "192.168.1.102", tc.api)
			t.Logf("API: %s, Expected: %v, Got: %v", tc.api, tc.expected, err)
			if tc.expected == nil {
				if err != nil {
					t.Errorf("API %s: expected success, got %v", tc.api, err)
				}
			} else {
				if err == nil || err.Error() != tc.expected.Error() {
					t.Errorf("API %s: expected %v, got %v", tc.api, tc.expected, err)
				}
			}
		})
	}
}
