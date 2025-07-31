package test

import (
	"testing"

	"github.com/windf17/wt"
	"github.com/windf17/wt/models"
	"github.com/windf17/wt/utility"
)

/**
 * TestHasPermissionWithApiRules 测试HasPermission函数的API规则匹配逻辑
 * @param {*} t 测试对象
 */
func TestHasPermissionWithApiRules(t *testing.T) {
	tests := []struct {
		name     string
		urlStr   string
		apiRules []models.ApiRule
		expected bool
	}{
		// 基础匹配测试
		{
			name:   "完全匹配-允许规则",
			urlStr: "/api/v1/users",
			apiRules: []models.ApiRule{
				{Path: []string{"api", "v1", "users"}, Rule: true},
			},
			expected: true,
		},
		{
			name:   "完全匹配-拒绝规则",
			urlStr: "/api/v1/admin",
			apiRules: []models.ApiRule{
				{Path: []string{"api", "v1", "admin"}, Rule: false},
			},
			expected: false,
		},
		{
			name:   "前缀匹配-允许规则",
			urlStr: "/api/v1/users/123",
			apiRules: []models.ApiRule{
				{Path: []string{"api", "v1", "users"}, Rule: true},
			},
			expected: true,
		},
		{
			name:   "前缀匹配-拒绝规则",
			urlStr: "/api/v1/admin/delete",
			apiRules: []models.ApiRule{
				{Path: []string{"api", "v1", "admin"}, Rule: false},
			},
			expected: false,
		},

		// 关键测试案例：/api/logout 应该匹配 /api 但不匹配 /api/admin
		{
			name:   "logout接口权限测试-应该匹配api而非admin",
			urlStr: "/api/logout",
			apiRules: []models.ApiRule{
				{Path: []string{"api"}, Rule: true},
				{Path: []string{"api", "admin"}, Rule: false},
			},
			expected: true, // 应该匹配 /api 规则，而不是 /api/admin
		},
		{
			name:   "admin接口权限测试-应该匹配admin规则",
			urlStr: "/api/admin/users",
			apiRules: []models.ApiRule{
				{Path: []string{"api"}, Rule: true},
				{Path: []string{"api", "admin"}, Rule: false},
			},
			expected: false, // 应该匹配 /api/admin 规则
		},

		// 多规则优先级测试
		{
			name:   "最长匹配优先-拒绝规则胜出",
			urlStr: "/api/v1/users/profile",
			apiRules: []models.ApiRule{
				{Path: []string{"api", "v1"}, Rule: true},
				{Path: []string{"api", "v1", "users"}, Rule: false},
			},
			expected: false, // 更长的匹配优先
		},
		{
			name:   "最长匹配优先-允许规则胜出",
			urlStr: "/api/v1/users/profile",
			apiRules: []models.ApiRule{
				{Path: []string{"api", "v1"}, Rule: false},
				{Path: []string{"api", "v1", "users"}, Rule: true},
			},
			expected: true, // 更长的匹配优先
		},
		{
			name:   "相同长度规则-先出现的优先",
			urlStr: "/api/v1/users/delete",
			apiRules: []models.ApiRule{
				{Path: []string{"api", "v1"}, Rule: true},
				{Path: []string{"api", "v1"}, Rule: false},
			},
			expected: true, // 先出现的规则优先
		},

		// 复杂场景测试
		{
			name:   "三层规则优先级测试",
			urlStr: "/api/v2/products/123/reviews",
			apiRules: []models.ApiRule{
				{Path: []string{"api", "v2"}, Rule: true},
				{Path: []string{"api", "v2", "products"}, Rule: true},
				{Path: []string{"api", "v2", "products", "123", "reviews"}, Rule: false},
			},
			expected: false, // 最长匹配的拒绝规则生效
		},
		{
			name:   "部分匹配不生效测试",
			urlStr: "/api/v1/users",
			apiRules: []models.ApiRule{
				{Path: []string{"api", "v1", "users", "admin"}, Rule: false}, // 规则比请求路径长
				{Path: []string{"api", "v1"}, Rule: true},
			},
			expected: true, // 长规则不匹配，短规则生效
		},

		// 边界情况测试
		{
			name:   "无匹配规则-默认拒绝",
			urlStr: "/api/v3/unknown",
			apiRules: []models.ApiRule{
				{Path: []string{"api", "v1"}, Rule: true},
				{Path: []string{"api", "v2"}, Rule: true},
			},
			expected: false,
		},
		{
			name:   "空路径测试",
			urlStr: "",
			apiRules: []models.ApiRule{
				{Path: []string{"api"}, Rule: true},
			},
			expected: false,
		},
		{
			name:   "根路径测试",
			urlStr: "/",
			apiRules: []models.ApiRule{
				{Path: []string{}, Rule: true},
			},
			expected: false, // 空路径段不匹配
		},
		{
			name:   "单段路径匹配",
			urlStr: "/api",
			apiRules: []models.ApiRule{
				{Path: []string{"api"}, Rule: true},
			},
			expected: true,
		},
		{
			name:   "深层路径前缀匹配",
			urlStr: "/api/v1/users/123/posts/456/comments/789",
			apiRules: []models.ApiRule{
				{Path: []string{"api", "v1", "users"}, Rule: true},
			},
			expected: true,
		},

		// 特殊字符和编码测试
		{
			name:   "带查询参数的路径",
			urlStr: "/api/users?id=123&name=test",
			apiRules: []models.ApiRule{
				{Path: []string{"api", "users"}, Rule: true},
			},
			expected: true, // 查询参数应该被忽略
		},
		{
			name:   "带锚点的路径",
			urlStr: "/api/users#section1",
			apiRules: []models.ApiRule{
				{Path: []string{"api", "users"}, Rule: true},
			},
			expected: true, // 锚点应该被忽略
		},

		// 空规则测试
		{
			name:     "空规则列表",
			urlStr:   "/api/users",
			apiRules: []models.ApiRule{},
			expected: false, // 无规则默认拒绝
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utility.HasPermission(tt.urlStr, tt.apiRules)
			if result != tt.expected {
				t.Errorf("HasPermission(%s, %v) = %v, expected %v", tt.urlStr, tt.apiRules, result, tt.expected)
			}
		})
	}
}

/**
 * TestAuthWithNoGroups 测试当没有配置用户组时的鉴权行为
 * 验证当groups参数为nil时，应该禁用鉴权功能，所有请求都应该通过
 */
func TestAuthWithNoGroups(t *testing.T) {
	// 配置
	config := &models.ConfigRaw{
		Language:       "zh",
		MaxTokens:      100,
		Delimiter:      ",",
		TokenRenewTime: "1h",
	}

	// 初始化token管理器，groups参数为空数组
	tm, err := wt.InitTM[map[string]any](*config, []models.GroupRaw{})
	if err != nil {
		t.Fatalf("Failed to initialize token manager: %v", err)
	}

	// 测试1: 验证没有用户组时，Auth应该直接返回成功
	t.Run("AuthWithoutToken", func(t *testing.T) {
		// 直接调用Auth，不需要有效的token
		err := tm.Auth("any_token", "192.168.1.1", "/api/test")
		if err != nil {
			t.Errorf("Expected nil when no groups configured, got %v", err)
		}
	})

	// 测试2: 验证Auth方法在没有groups时的行为
	t.Run("AuthWithoutGroups", func(t *testing.T) {
		// 当没有groups配置时，Auth应该直接返回成功
		// 使用符合格式要求的token（至少32个字符的Base64格式）
		validToken := "dGVzdF90b2tlbl9mb3JfYXV0aF90ZXN0aW5nX3B1cnBvc2U="
		err := tm.Auth(validToken, "192.168.1.1", "/api/test")
		if err != nil {
			t.Errorf("Expected nil when no groups configured, got %v", err)
		}
	})
}

/**
 * TestAuthWithEmptyGroups 测试空用户组数组的情况
 */
func TestAuthWithEmptyGroups(t *testing.T) {
	config := &models.ConfigRaw{
		Language:       "zh",
		MaxTokens:      100,
		Delimiter:      ",",
		TokenRenewTime: "1h",
	}

	// 传入空的groups数组
	emptyGroups := []models.GroupRaw{}
	tm, err := wt.InitTM[map[string]any](*config, emptyGroups)
	if err != nil {
		t.Fatalf("Failed to initialize token manager: %v", err)
	}

	// 验证空用户组时的行为
	err = tm.Auth("any_token", "192.168.1.1", "/api/test")
	if err != nil {
		t.Errorf("Expected nil when empty groups configured, got %v", err)
	}
}

/**
 * TestAuthWithGroups 测试有用户组配置时的正常鉴权行为
 */
func TestAuthWithGroups(t *testing.T) {
	config := &models.ConfigRaw{
		Language:       "zh",
		MaxTokens:      100,
		Delimiter:      ",",
		TokenRenewTime: "1h",
	}

	groups := []models.GroupRaw{
		{
			ID:                 1,
			Name:               "test_group",
			AllowedAPIs:        "/api/user",
			DeniedAPIs:         "/api/admin",
			TokenExpire:        "3600s",
			AllowMultipleLogin: 1,
		},
	}

	tm, err := wt.InitTM[string](*config, groups)
	if err != nil {
		t.Errorf("Failed to initialize token manager: %v", err)
	}

	// 添加一个有效的token
	tokenKey, err := tm.AddToken(1, 1, "192.168.1.1")
	if err != nil {
		t.Fatalf("Failed to add token: %v", err)
	}

	// 测试有效token访问允许的API
	err = tm.Auth(tokenKey, "192.168.1.1", "/api/user/profile")
	if err != nil {
		t.Errorf("Expected nil for allowed API, got %v", err)
	}

	// 测试有效token访问被拒绝的API
	err = tm.Auth(tokenKey, "192.168.1.1", "/api/admin/delete")
	if err == nil {
		t.Error("Expected error for denied API, got nil")
	}

	// 测试访问不在任何规则中的API
	err = tm.Auth(tokenKey, "192.168.1.1", "/api/other/test")
	if err == nil {
		t.Error("Expected error for unmatched API, got nil")
	}

	// 测试无效token
	err = tm.Auth("invalid_token", "192.168.1.1", "/api/user/profile")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

/**
 * TestBatchAuth 测试批量API权限检查功能
 * @param {*} t 测试对象
 */
func TestBatchAuth(t *testing.T) {
	// 创建配置
	config := &models.ConfigRaw{
		Language:       "zh",
		MaxTokens:      1000,
		Delimiter:      ",",
		TokenRenewTime: "1h",
	}

	// 创建用户组
	groups := []models.GroupRaw{
		{
			ID:                 1,
			Name:               "测试用户组",
			TokenExpire:        "3600s",
			AllowMultipleLogin: 1,
			AllowedAPIs:        "/api/user,/api/public",
			DeniedAPIs:         "/api/admin",
		},
	}

	// 初始化token管理器
	tm, err := wt.InitTM[map[string]any](*config, groups)
	if err != nil {
		t.Fatalf("Failed to initialize token manager: %v", err)
	}

	// 添加token
	tokenKey, err := tm.AddToken(1, 1, "192.168.1.1")
	if err != nil {
		t.Fatalf("Failed to add token: %v", err)
	}

	// 测试批量权限检查
	apis := []string{
		"/api/user/del",
		"/api/user/add",
		"/api/user/update",
		"/api/user/get",
		"/api/admin/delete",
		"/api/public/info",
	}

	// 执行批量权限检查
	results := tm.BatchAuth(tokenKey, "192.168.1.1", apis)

	// 验证结果数组长度
	if len(results) != len(apis) {
		t.Errorf("Expected results length %d, got %d", len(apis), len(results))
	}

	// 验证每个API的权限结果
	expected := []bool{
		true,  // /api/user/del - 应该有权限（匹配/api/user前缀）
		true,  // /api/user/add - 应该有权限（匹配/api/user前缀）
		true,  // /api/user/update - 应该有权限（匹配/api/user前缀）
		true,  // /api/user/get - 应该有权限（匹配/api/user前缀）
		false, // /api/admin/delete - 应该无权限（被拒绝规则阻止）
		true,  // /api/public/info - 应该有权限（匹配/api/public前缀）
	}

	for i, expectedResult := range expected {
		if results[i] != expectedResult {
			t.Errorf("API %s: expected %v, got %v", apis[i], expectedResult, results[i])
		}
	}

	// 测试空API数组
	emptyResults := tm.BatchAuth(tokenKey, "192.168.1.1", []string{})
	if len(emptyResults) != 0 {
		t.Errorf("Expected empty results for empty API array, got length %d", len(emptyResults))
	}

	// 测试无效token的批量检查
	invalidResults := tm.BatchAuth("invalid_token", "192.168.1.1", apis)
	for i, result := range invalidResults {
		if result != false {
			t.Errorf("API %s with invalid token: expected false, got %v", apis[i], result)
		}
	}

	// 测试IP不匹配的批量检查
	ipMismatchResults := tm.BatchAuth(tokenKey, "192.168.1.2", apis)
	for i, result := range ipMismatchResults {
		if result != false {
			t.Errorf("API %s with mismatched IP: expected false, got %v", apis[i], result)
		}
	}
}
