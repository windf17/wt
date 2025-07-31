package test

import (
	"os"
	"testing"
	"time"

	"github.com/windf17/wt"
	"github.com/windf17/wt/models"
)

/**
 * TestConfigValidation 测试配置验证功能
 */
func TestConfigValidation(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		config := &models.ConfigRaw{

		MaxTokens:      1000,
		Delimiter:      ",",
		TokenRenewTime: "1h",
		Language:       "zh",
	}

		// 测试配置是否能正常初始化
		tm, err := wt.InitTM[string](*config, []models.GroupRaw{})
		if err != nil {
			t.Errorf("Valid config should initialize successfully, got error: %v", err)
		}
		if tm == nil {
			t.Error("Valid config should initialize successfully")
		} else {

		}
	})

	t.Run("InvalidMaxTokens", func(t *testing.T) {
		config := &models.ConfigRaw{
			MaxTokens:      -1, // 无效值（负数）
			Delimiter:      ",",
			TokenRenewTime: "1h",
			Language:       "zh",
		}

		// 测试无效配置会返回nil（统一错误处理）
		tm, _ := wt.InitTM[string](*config, []models.GroupRaw{})
		if tm != nil {
			t.Error("Expected nil for invalid config, but initialization succeeded")

		} else {
			t.Log("Invalid config correctly returned nil")
		}
	})

	t.Run("EmptyDelimiter", func(t *testing.T) {
		config := &models.ConfigRaw{
			MaxTokens:      1000,
			Delimiter:      "", // 空分隔符
			TokenRenewTime: "1h",
			Language:       "zh",
		}

		// 测试空分隔符配置会返回nil（统一错误处理）
		tm, _ := wt.InitTM[string](*config, []models.GroupRaw{})
		if tm != nil {
			t.Error("Expected nil for empty delimiter config, but initialization succeeded")

		} else {
			t.Log("Empty delimiter config correctly returned nil")
		}
	})
}

/**
 * TestCompleteWorkflow 测试完整的工作流程
 */
func TestCompleteWorkflow(t *testing.T) {
	// 创建临时缓存文件
	tempFile := "test_cache.json"
	defer os.Remove(tempFile)

	// 初始化配置
	config := &models.ConfigRaw{

		MaxTokens:      10,
		Delimiter:      ",",
		TokenRenewTime: "30m",
		Language:       "zh",
	}

	// 验证配置
	if err := wt.ValidateConfig(*config); err != nil {
		t.Fatalf("Config validation failed: %v", err)
	}

	groups := []models.GroupRaw{
		{
			ID:                 1,
			Name:               "admin",
			AllowedAPIs:        "/api/admin,/api/user,/api/admin/users",
			DeniedAPIs:         "/api/admin/delete",
			TokenExpire:        "2h",
			AllowMultipleLogin: 0,
		},
		{
			ID:                 2,
			Name:               "user",
			AllowedAPIs:        "/api/user/profile",
			DeniedAPIs:         "/api/user/admin",
			TokenExpire:        "1h",
			AllowMultipleLogin: 1,
		},
	}

	// 验证用户组配置
	for _, group := range groups {
		if err := wt.ValidateGroupRaw(group); err != nil {
			t.Fatalf("Group validation failed for group %s: %v", group.Name, err)
		}
	}

	// 初始化token管理器
	tm, err := wt.InitTM[map[string]any](*config, groups)
	if err != nil {
		t.Fatalf("Failed to initialize token manager: %v", err)
	}
	if tm == nil {
		t.Fatalf("Failed to initialize token manager")
	}

	// 测试用户组管理
	t.Run("GroupManagement", func(t *testing.T) {
		// 获取用户组
		group, err := tm.GetGroup(1)
		if err != nil {
			t.Errorf("Failed to get group: %v", err)
		}
		if group.Name != "admin" {
			t.Errorf("Expected group name 'admin', got '%s'", group.Name)
		}

		// 添加新用户组
		newGroup := &models.GroupRaw{
			ID:                 3,
			Name:               "guest",
			AllowedAPIs:        "/api/public/*",
			DeniedAPIs:         "",
			TokenExpire:        "30m",
			AllowMultipleLogin: 1,
		}
		err = tm.AddGroup(newGroup)
		if err != nil {
			t.Errorf("Failed to add group: %v", err)
		}

		// 验证新用户组
		addedGroup, err := tm.GetGroup(3)
		if err != nil {
			t.Errorf("Failed to get added group: %v", err)
		}
		if addedGroup.Name != "guest" {
			t.Errorf("Expected group name 'guest', got '%s'", addedGroup.Name)
		}
	})

	// 测试token管理
	t.Run("TokenManagement", func(t *testing.T) {
		// 添加token
		token1, err := tm.AddToken(1, 1, "192.168.1.100")
		if err != nil {
			t.Errorf("Failed to add token: %v", err)
			return
		}

		// 验证token
		tokenData, err := tm.GetToken(token1)
		if err != nil {
			t.Errorf("Failed to get token: %v", err)
			return
		}
		if tokenData.UserID != 1 {
			t.Errorf("Expected user ID 1, got %d", tokenData.UserID)
		}
		if tokenData.GroupID != 1 {
			t.Errorf("Expected group ID 1, got %d", tokenData.GroupID)
		}

		// 设置用户数据前先验证token1是否有效
		_, err = tm.GetToken(token1)
		if err != nil {
			t.Errorf("Token1 is invalid before SetUserData: %v", err)
			return
		}

		// 设置用户数据
		userData := map[string]any{
			"username": "admin_user",
			"email":    "admin@example.com",
			"role":     "administrator",
		}
		// 打印token1的详细信息用于调试
		t.Logf("Token1: %s, Length: %d", token1, len(token1))

		err = tm.SetUserData(token1, userData)
		if err != nil {
			t.Errorf("Failed to set user data: %v (token: %s)", err, token1)
			return
		}

		// 获取用户数据
		retrievedData, err := tm.GetUserData(token1)
		if err != nil {
			t.Errorf("Failed to get user data: %v", err)
			return
		}
		if retrievedData["username"] != "admin_user" {
			t.Errorf("Expected username 'admin_user', got '%v'", retrievedData["username"])
		}

		// 测试同一用户的多设备登录限制
		token3, err := tm.AddToken(1, 1, "192.168.1.102")
		if err != nil {
			t.Errorf("Failed to add third token: %v", err)
			return
		}

		// 第一个token应该被删除（同一用户，不允许多设备登录）
		_, err = tm.GetToken(token1)
		if err == nil {
			t.Errorf("Expected first token to be invalid after adding third token, but it's still valid")
		}

		// 第三个token应该有效
		_, err = tm.GetToken(token3)
		if err != nil {
			t.Errorf("Expected third token to be valid, got: %v", err)
		}

		// 测试不同用户的token（应该不受影响）
		token2, err := tm.AddToken(2, 1, "192.168.1.101")
		if err != nil {
			t.Errorf("Failed to add second token: %v", err)
			return
		}

		// token2应该有效（不同用户）
		_, err = tm.GetToken(token2)
		if err != nil {
			t.Errorf("Expected second token to be valid, got: %v", err)
		}
	})

	// 测试API权限验证
	t.Run("APIPermissions", func(t *testing.T) {
		// 创建用户token
		userToken, err := tm.AddToken(2, 2, "192.168.1.102")
		if err != nil {
			t.Errorf("Failed to add user token: %v", err)
		}

		// 测试允许的API
		err = tm.Auth(userToken, "192.168.1.102", "/api/user/profile")
		if err != nil {
			t.Errorf("Expected API access to be allowed, got: %v", err)
		}

		// 测试被拒绝的API
		err = tm.Auth(userToken, "192.168.1.102", "/api/admin/delete")
		if err == nil {
			t.Errorf("Expected API access to be denied, but it was allowed")
		}
	})

	// 测试批量操作
	t.Run("BatchOperations", func(t *testing.T) {
		// 创建多个用户的token
		userIDs := []uint{10, 11, 12}
		for _, userID := range userIDs {
			_, err := tm.AddToken(userID, 2, "192.168.1.200")
			if err != nil {
				t.Errorf("Failed to add token for user %d: %v", userID, err)
			}
		}

		// 批量删除用户token
		err := tm.BatchDeleteTokensByUserIDs(userIDs)
		if err != nil {
			t.Errorf("Failed to batch delete tokens: %v", err)
		}

		// 验证token已被删除
		tokens := tm.GetTokensByUserID(10)
		if len(tokens) != 0 {
			t.Errorf("Expected no tokens for user 10, got %d", len(tokens))
		}
	})

	// 测试token过期
	t.Run("TokenExpiration", func(t *testing.T) {
		// 创建短期token
		shortToken, err := tm.AddToken(20, 2, "192.168.1.200") // 60秒过期
		if err != nil {
			t.Errorf("Failed to add short-term token: %v", err)
			return
		}

		// 手动设置token为过期状态（通过修改过期时间）
		token, err := tm.GetToken(shortToken)
		if err != nil {
			t.Errorf("Failed to get token: %v", err)
			return
		}
		// 将过期时间设置为1秒前
		token.ExpireSeconds = 1
		token.LoginTime = time.Now().Add(-2 * time.Second)
		tm.UpdateToken(shortToken, token)

		// 验证token已过期
		_, err = tm.GetToken(shortToken)
		if err == nil {
			t.Errorf("Expected token to be expired, but it's still valid")
		}
	})

	// 测试统计信息
	t.Run("Statistics", func(t *testing.T) {
		stats := tm.GetStats()
		if stats.TotalTokens < 0 {
			t.Errorf("Invalid total tokens count: %d", stats.TotalTokens)
		}
		if stats.ActiveTokens < 0 {
			t.Errorf("Invalid active tokens count: %d", stats.ActiveTokens)
		}
		if stats.ExpiredTokens < 0 {
			t.Errorf("Invalid expired tokens count: %d", stats.ExpiredTokens)
		}
	})

	// 测试缓存文件持久化
	t.Run("CachePersistence", func(t *testing.T) {
		// 添加一些token
		for i := 30; i < 35; i++ {
			_, err := tm.AddToken(uint(i), 2, "192.168.1.100")
			if err != nil {
				t.Errorf("Failed to add token for persistence test: %v", err)
			}
		}

		// 关闭当前管理器以触发最终备份


		// 缓存功能已移除，跳过缓存文件检查

		// 重新初始化token管理器以测试加载
		tm2, err := wt.InitTM[map[string]any](*config, groups)
	if err != nil {
		t.Errorf("Failed to reinitialize token manager: %v", err)
	}
		if tm2 == nil {
			t.Errorf("Failed to reinitialize token manager")
		}


		// 验证缓存功能已移除，不会加载任何数据
		stats := tm2.GetStats()
		if stats.TotalTokens != 0 {
			t.Errorf("Tokens should not be loaded since cache functionality is removed, got %d tokens", stats.TotalTokens)
		} else {
			t.Logf("Cache functionality correctly removed: no tokens loaded")
		}

		// 重新赋值tm以便后续测试使用
		tm = tm2
	})

	// 清理过期token
	t.Run("CleanExpiredTokens", func(t *testing.T) {
		initialStats := tm.GetStats()
		tm.CleanExpiredTokens()
		finalStats := tm.GetStats()

		// 验证过期token被清理
		if finalStats.ExpiredTokens > initialStats.ExpiredTokens {
			t.Errorf("Expired tokens count increased after cleanup")
		}
	})
}

/**
 * TestSecurityFeatures 测试安全功能
 */
func TestSecurityFeatures(t *testing.T) {
	// 测试token格式验证
	t.Run("TokenFormatValidation", func(t *testing.T) {
		validToken := "dGVzdF90b2tlbl9mb3JfdmFsaWRhdGlvbl90ZXN0aW5n" // "test_token_for_validation_testing" base64编码，长度32字节
		invalidToken := "invalid_token!"

		if !wt.ValidateTokenFormat(validToken) {
			t.Errorf("Valid token format was rejected")
		}

		if wt.ValidateTokenFormat(invalidToken) {
			t.Errorf("Invalid token format was accepted")
		}
	})

	// 测试输入清理
	t.Run("InputSanitization", func(t *testing.T) {
		dangerousInput := "test<script>alert('xss')</script>"
		sanitized := wt.SanitizeInput(dangerousInput)
		expected := "testscriptalertxssscript"

		if sanitized != expected {
			t.Errorf("Expected '%s', got '%s'", expected, sanitized)
		}
	})

	// 测试加密功能
	t.Run("EncryptionDecryption", func(t *testing.T) {
		sm := wt.NewSecurityManager("test_password")
		originalData := "sensitive_token_data"

		// 加密
		encrypted, err := sm.EncryptToken(originalData)
		if err != nil {
			t.Errorf("Encryption failed: %v", err)
		}

		// 解密
		decrypted, err := sm.DecryptToken(encrypted)
		if err != nil {
			t.Errorf("Decryption failed: %v", err)
		}

		if decrypted != originalData {
			t.Errorf("Expected '%s', got '%s'", originalData, decrypted)
		}
	})

	// 测试哈希功能
	t.Run("HashFunction", func(t *testing.T) {
		sm := wt.NewSecurityManager("test_password")
		data := "sensitive_data"

		hash1 := sm.HashSensitiveData(data)
		hash2 := sm.HashSensitiveData(data)

		// 相同数据应该产生相同哈希
		if hash1 != hash2 {
			t.Errorf("Hash function is not deterministic")
		}

		// 不同数据应该产生不同哈希
		hash3 := sm.HashSensitiveData("different_data")
		if hash1 == hash3 {
			t.Errorf("Different data produced same hash")
		}
	})
}

/**
 * TestCompareApiRules 测试API规则排序功能
 */
func TestCompareApiRules(t *testing.T) {
	tests := []struct {
		name     string
		pathA    []string
		pathB    []string
		expected bool // true表示pathA应该排在pathB前面
	}{
		{
			name:     "长度不同-A更长",
			pathA:    []string{"api", "v1", "users"},
			pathB:    []string{"api", "v1"},
			expected: true,
		},
		{
			name:     "长度不同-B更长",
			pathA:    []string{"api", "v1"},
			pathB:    []string{"api", "v1", "users"},
			expected: false,
		},
		{
			name:     "长度相同-字符串长度不同",
			pathA:    []string{"api", "version"},
			pathB:    []string{"api", "v1"},
			expected: true, // "version"比"v1"长
		},
		{
			name:     "长度相同-字符串长度相同-字典序",
			pathA:    []string{"api", "v1"},
			pathB:    []string{"api", "v2"},
			expected: true, // "v1" < "v2"
		},
		{
			name:     "完全相同的路径",
			pathA:    []string{"api", "v1", "users"},
			pathB:    []string{"api", "v1", "users"},
			expected: false, // 保持原有顺序
		},
		{
			name:     "*号优先级最低-A是*",
			pathA:    []string{"api", "*"},
			pathB:    []string{"api", "v1"},
			expected: false, // *优先级低
		},
		{
			name:     "*号优先级最低-B是*",
			pathA:    []string{"api", "v1"},
			pathB:    []string{"api", "*"},
			expected: true, // *优先级低
		},
		{
			name:     "都是*号",
			pathA:    []string{"api", "*", "users"},
			pathB:    []string{"api", "*", "posts"},
			expected: false, // "users" > "posts"，但稳定排序返回false
		},
		{
			name:     "复杂情况-多个*号",
			pathA:    []string{"*", "v1", "users"},
			pathB:    []string{"api", "*", "users"},
			expected: false, // 第一个段A是*，B不是*
		},
		{
			name:     "混合长度和*号",
			pathA:    []string{"api", "version", "*"},
			pathB:    []string{"api", "v1", "users"},
			expected: true, // "version"比"v1"长
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 由于compareApiRules是包内函数，我们通过ConvGroup来间接测试
			// 这里我们需要创建一个测试辅助函数
			result := testCompareApiRules(tt.pathA, tt.pathB)
			if result != tt.expected {
				t.Errorf("compareApiRules(%v, %v) = %v, expected %v", tt.pathA, tt.pathB, result, tt.expected)
			}
		})
	}
}

/**
 * testCompareApiRules 测试辅助函数，通过ConvGroup间接测试排序逻辑
 * @param {[]string} pathA 第一个路径段数组
 * @param {[]string} pathB 第二个路径段数组
 * @returns {bool} 排序结果
 */
func testCompareApiRules(pathA, pathB []string) bool {
	// 对于完全相同的路径，直接返回false（保持原有顺序）
	if equalStringSlices(pathA, pathB) {
		return false
	}

	// 创建两个不同的API规则，通过ConvGroup的排序来测试
	// 为了避免重复路径问题，我们使用一个允许规则和一个拒绝规则
	allowedPath := joinPath(pathA)
	deniedPath := joinPath(pathB)

	raw := models.GroupRaw{
		Name:               "test",
		AllowedAPIs:        allowedPath,
		DeniedAPIs:         deniedPath,
		TokenExpire:        "3600",
		AllowMultipleLogin: 1,
	}

	group := wt.ConvGroup(raw, "|")

	// 检查排序结果
	if len(group.ApiRules) >= 2 {
		// 如果第一个规则的路径等于pathA，说明pathA排在前面
		return equalStringSlices(group.ApiRules[0].Path, pathA)
	}

	return false
}

/**
 * joinPath 将路径段数组连接为路径字符串
 * @param {[]string} path 路径段数组
 * @returns {string} 路径字符串
 */
func joinPath(path []string) string {
	if len(path) == 0 {
		return ""
	}
	result := "/"
	for i, segment := range path {
		if i > 0 {
			result += "/"
		}
		result += segment
	}
	return result
}

/**
 * equalStringSlices 比较两个字符串切片是否相等
 * @param {[]string} a 第一个字符串切片
 * @param {[]string} b 第二个字符串切片
 * @returns {bool} 是否相等
 */
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

/**
 * TestConvGroup 测试ConvGroup函数的完整功能
 */
func TestConvGroup(t *testing.T) {
	tests := []struct {
		name      string
		raw       models.GroupRaw
		delimiter string
		expected  *models.Group
	}{
		{
			name: "基本功能测试",
			raw: models.GroupRaw{
				Name:               "TestGroup",
				AllowedAPIs:        "/api/v1/users|/api/v2/posts",
				DeniedAPIs:         "/admin/*|/system/config",
				TokenExpire:        "3600",
				AllowMultipleLogin: 1,
			},
			delimiter: "|",
			expected: &models.Group{
				Name: "TestGroup",
				ApiRules: []models.ApiRule{
					{Path: []string{"api", "v1", "users"}, Rule: true},
					{Path: []string{"api", "v2", "posts"}, Rule: true},
					{Path: []string{"system", "config"}, Rule: false},
					{Path: []string{"admin", "*"}, Rule: false},
				},
				ExpireSeconds:      3600,
				AllowMultipleLogin: true,
			},
		},
		{
			name: "排序测试-长度优先",
			raw: models.GroupRaw{
				Name:               "SortTest",
				AllowedAPIs:        "/api|/api/v1|/api/v1/users/profile",
				DeniedAPIs:         "",
				TokenExpire:        "0",
				AllowMultipleLogin: 0,
			},
			delimiter: "|",
			expected: &models.Group{
				Name: "SortTest",
				ApiRules: []models.ApiRule{
					{Path: []string{"api", "v1", "users", "profile"}, Rule: true},
					{Path: []string{"api", "v1"}, Rule: true},
					{Path: []string{"api"}, Rule: true},
				},
				ExpireSeconds:      0,
				AllowMultipleLogin: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wt.ConvGroup(tt.raw, tt.delimiter)

			// 验证基本属性
			if result.Name != tt.expected.Name {
				t.Errorf("Name = %v, expected %v", result.Name, tt.expected.Name)
			}
			if result.ExpireSeconds != tt.expected.ExpireSeconds {
				t.Errorf("ExpireSeconds = %v, expected %v", result.ExpireSeconds, tt.expected.ExpireSeconds)
			}
			if result.AllowMultipleLogin != tt.expected.AllowMultipleLogin {
				t.Errorf("AllowMultipleLogin = %v, expected %v", result.AllowMultipleLogin, tt.expected.AllowMultipleLogin)
			}

			// 验证API规则数量
			if len(result.ApiRules) != len(tt.expected.ApiRules) {
				t.Errorf("ApiRules length = %v, expected %v", len(result.ApiRules), len(tt.expected.ApiRules))
				return
			}

			// 验证每个API规则
			for i, rule := range result.ApiRules {
				expectedRule := tt.expected.ApiRules[i]
				if rule.Rule != expectedRule.Rule {
					t.Errorf("ApiRules[%d].Rule = %v, expected %v", i, rule.Rule, expectedRule.Rule)
				}
				if !equalStringSlices(rule.Path, expectedRule.Path) {
					t.Errorf("ApiRules[%d].Path = %v, expected %v", i, rule.Path, expectedRule.Path)
				}
			}
		})
	}
}

/**
 * equalStringSlices 比较两个字符串切片是否相等
 */

/**
 * TestConvGroupEdgeCases 测试ConvGroup函数的边界情况
 */
func TestConvGroupEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		raw       models.GroupRaw
		delimiter string
	}{
		{
			name: "空API列表",
			raw: models.GroupRaw{
				Name:               "EmptyAPIs",
				AllowedAPIs:        "",
				DeniedAPIs:         "",
				TokenExpire:        "3600",
				AllowMultipleLogin: 1,
			},
			delimiter: "|",
		},
		{
			name: "包含空白字符的API",
			raw: models.GroupRaw{
				Name:               "WhitespaceAPIs",
				AllowedAPIs:        " /api/v1 | /api/v2 | ",
				DeniedAPIs:         " /admin | /system ",
				TokenExpire:        "7200",
				AllowMultipleLogin: 0,
			},
			delimiter: "|",
		},
		{
			name: "包含*号的API",
			raw: models.GroupRaw{
				Name:               "WildcardAPIs",
				AllowedAPIs:        "/api/*/users|/api/v1/*|/*",
				DeniedAPIs:         "/admin/*/config",
				TokenExpire:        "1800",
				AllowMultipleLogin: 1,
			},
			delimiter: "|",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wt.ConvGroup(tt.raw, tt.delimiter)

			// 基本验证：确保函数不会崩溃并返回有效结果
			if result == nil {
				t.Error("ConvGroup returned nil")
				return
			}

			if result.Name != tt.raw.Name {
				t.Errorf("Name = %v, expected %v", result.Name, tt.raw.Name)
			}

			// 验证API规则是否正确排序（长度递减）
			for i := 1; i < len(result.ApiRules); i++ {
				prevLen := len(result.ApiRules[i-1].Path)
				currLen := len(result.ApiRules[i].Path)
				if prevLen < currLen {
					t.Errorf("API rules not properly sorted: rule %d has length %d, rule %d has length %d", i-1, prevLen, i, currLen)
				}
			}
		})
	}
}

/**
 * TestGroupValidation 测试用户组验证功能
 */
func TestGroupValidation(t *testing.T) {
	config := &models.ConfigRaw{

		MaxTokens:      100,
		Delimiter:      ",",
		TokenRenewTime: "1h",
		Language:       "zh",
	}

	t.Run("ValidGroup", func(t *testing.T) {
		groups := []models.GroupRaw{
			{
				ID:                 1,
				Name:               "test_group",
				AllowedAPIs:        "/api/user,/api/profile",
				DeniedAPIs:         "/api/admin",
				TokenExpire:        "1h",
				AllowMultipleLogin: 1,
			},
		}

		// 测试有效用户组配置是否能正常初始化
		tm, err := wt.InitTM[string](*config, groups)
	if err != nil {
		t.Fatalf("Failed to initialize token manager: %v", err)
	}
		if tm == nil {
			t.Error("Valid group should initialize successfully")
		} else {
			// 测试能否获取用户组信息
			group, err := tm.GetGroup(1)
			if err != nil {
				t.Errorf("Should be able to get valid group: %v", err)
			} else if group.Name != "test_group" {
				t.Errorf("Expected group name 'test_group', got '%s'", group.Name)
			}

		}
	})

	t.Run("InvalidGroupID", func(t *testing.T) {
		groups := []models.GroupRaw{
			{
				ID:                 0, // 无效ID
				Name:               "test_group",
				AllowedAPIs:        "/api/user",
				DeniedAPIs:         "",
				TokenExpire:        "1h",
				AllowMultipleLogin: 1,
			},
		}

		// 测试无效用户组ID会返回nil（统一错误处理）
		tm, _ := wt.InitTM[string](*config, groups)
		if tm != nil {
			t.Error("Expected nil for invalid group ID, but initialization succeeded")

		} else {
			t.Log("Invalid group ID correctly returned nil")
		}
	})

	t.Run("EmptyGroupName", func(t *testing.T) {
		groups := []models.GroupRaw{
			{
				ID:                 1,
				Name:               "", // 空名称
				AllowedAPIs:        "/api/user",
				DeniedAPIs:         "",
				TokenExpire:        "1h",
				AllowMultipleLogin: 1,
			},
		}

		// 测试空用户组名称会返回nil（统一错误处理）
		tm, _ := wt.InitTM[string](*config, groups)
		if tm != nil {
			t.Error("Expected nil for empty group name, but initialization succeeded")

		} else {
			t.Log("Empty group name correctly returned nil")
		}
	})

	t.Run("InvalidTokenExpire", func(t *testing.T) {
		groups := []models.GroupRaw{
			{
				ID:                 1,
				Name:               "test_group",
				AllowedAPIs:        "/api/user",
				DeniedAPIs:         "",
				TokenExpire:        "invalid_expire", // 无效过期时间
				AllowMultipleLogin: 1,
			},
		}

		// 测试无效过期时间会返回nil（统一错误处理）
		tm, _ := wt.InitTM[string](*config, groups)
		if tm != nil {
			t.Error("Expected nil for invalid token expire, but initialization succeeded")

		} else {
			t.Log("Invalid token expire correctly returned nil")
		}
	})
}

/**
 * TestTokenSecurity 测试Token安全功能
 */
func TestTokenSecurity(t *testing.T) {
	// 初始化测试环境
	config := &models.ConfigRaw{
		MaxTokens:      100,
		Delimiter:      ",",
		TokenRenewTime: "1h",
		Language:       "zh",
	}

	groups := []models.GroupRaw{
		{
			ID:                 1,
			Name:               "test_group",
			AllowedAPIs:        "/api/user",
			DeniedAPIs:         "/api/admin",
			TokenExpire:        "1s", // 设置为1秒过期，用于测试
			AllowMultipleLogin: 0,    // 不允许多次登录，启用IP验证
		},
	}

	tm, err := wt.InitTM[string](*config, groups)
	if err != nil {
		t.Fatalf("Failed to initialize token manager: %v", err)
	}
	defer

	t.Run("TokenFormat", func(t *testing.T) {
		// 测试Token生成和验证
		token, err := tm.GenerateToken()
		if err != nil {
			t.Errorf("Failed to generate token: %v", err)
			return
		}

		// 验证生成的Token格式
		if len(token) == 0 {
			t.Error("Generated token should not be empty")
		}

		// 测试无效Token验证
		invalidTokens := []string{
			"",         // 空Token
			"abc",      // 太短
			"invalid!", // 包含非法字符
		}

		for _, invalidToken := range invalidTokens {
			err := tm.Auth(invalidToken, "192.168.1.1", "/api/user")
			if err == nil {
				t.Errorf("Invalid token should be rejected: %s", invalidToken)
			}
		}
	})

	t.Run("TokenExpiration", func(t *testing.T) {
		// 添加一个短期Token
		token, err := tm.AddToken(1, 1, "192.168.1.1") // 1秒过期，不允许多设备登录
		if err != nil {
			t.Fatalf("Failed to add token: %v", err)
		}

		// 立即验证应该成功
		checkResult := tm.Auth(token, "192.168.1.1", "/api/user")
		if checkResult != nil {
			t.Errorf("Token should be valid immediately after creation")
		}

		// 等待Token过期
		time.Sleep(2 * time.Second)

		// 验证应该失败
		checkResult = tm.Auth(token, "192.168.1.1", "/api/user")
		if checkResult == nil {
			t.Error("Expired token should be rejected")
		}
	})

	t.Run("IPValidation", func(t *testing.T) {
		// 添加绑定IP的Token
		token, err := tm.AddToken(1, 1, "192.168.1.100") // false表示不允许多设备登录，会验证IP
		if err != nil {
			t.Fatalf("Failed to add token: %v", err)
		}

		// 验证Token是否有效
		authResult := tm.Auth(token, "192.168.1.100", "/api/user")
		if authResult != nil {
			t.Error("Valid token should pass authentication")
		}

		// 获取Token信息验证IP绑定
		tokenInfo, getErr := tm.GetToken(token)
		if getErr != nil {
			t.Errorf("Failed to get token info: %v", getErr)
		} else if tokenInfo.IP != "192.168.1.100" {
			t.Errorf("Expected IP 192.168.1.100, got %s", tokenInfo.IP)
		}

		// 测试IP验证的重要性：验证token与IP的绑定
		if tokenInfo.IP == "192.168.1.100" {
			t.Log("Token+IP验证功能正常：Token已正确绑定到指定IP")
		} else {
			t.Errorf("IP绑定失败：期望IP为192.168.1.100，实际为%s", tokenInfo.IP)
		}

		// 测试从不同IP验证同一token（应该失败）
		// 使用Auth方法进行IP验证
		err2 := tm.Auth(token, "192.168.1.200", "/api/user") // 不同的IP
		if err2 == nil {
			t.Log("IP验证安全功能正常：Token在非绑定IP上验证失败，返回E_Forbidden")
			t.Error("安全漏洞：Token应该只能在绑定的IP上使用，但在不同IP上验证成功")
		} else {
			t.Logf("IP验证返回其他错误码: %v", err2)
		}

		// 测试用正确的IP验证token（应该成功）
		err3 := tm.Auth(token, "192.168.1.100", "/api/user") // 正确的IP
		if err3 == nil {
			t.Log("IP验证功能正常：Token在绑定IP上验证成功")
		} else {
			t.Errorf("Token在绑定IP上验证失败: %v", err3)
		}

		// 清理
		tm.DelToken(token)
	})
}

/**
 * TestBoundaryConditions 测试边界条件
 */
func TestBoundaryConditions(t *testing.T) {
	config := &models.ConfigRaw{

		MaxTokens:      2, // 设置很小的最大Token数
		Delimiter:      ",",
		TokenRenewTime: "1h",
		Language:       "zh",
	}

	groups := []models.GroupRaw{
		{
			ID:                 1,
			Name:               "test_group",
			AllowedAPIs:        "/api/user",
			DeniedAPIs:         "",
			TokenExpire:        "1h",
			AllowMultipleLogin: 0, // 不允许多重登录
		},
	}

	tm, err := wt.InitTM[string](*config, groups)
	if err != nil {
		t.Fatalf("Failed to initialize token manager: %v", err)
	}
	defer

	t.Run("MaxTokensLimit", func(t *testing.T) {
		// 添加Token直到达到限制
		token1, err1 := tm.AddToken(1, 1, "192.168.1.1")
		if err1 != nil {
			t.Errorf("First token should be added successfully")
		}

		token2, err2 := tm.AddToken(2, 1, "192.168.1.2")
		if err2 != nil {
			t.Errorf("Second token should be added successfully")
		}

		// 添加一个小延迟确保token的LastAccessTime不同
		time.Sleep(10 * time.Millisecond)

		// 访问第二个token来更新它的LastAccessTime，确保第一个token是最久没有使用的
		_, _ = tm.GetToken(token2)
		time.Sleep(10 * time.Millisecond)

		// 第三个Token应该能成功添加，但会触发LRU清理删除最久没使用的token（token1）
		token3, err3 := tm.AddToken(3, 1, "192.168.1.3")
		if err3 != nil {
			t.Errorf("Third token should be added successfully: %v", err3)
		} else {
			t.Logf("Third token successfully added: %s", token3)
		}

		// 验证第一个token应该被删除（LRU策略）
		_, getErr1 := tm.GetToken(token1)
		if getErr1 == nil {
			t.Errorf("First token should be deleted by LRU policy")
		}

		// 验证第二个token仍然有效
		_, getErr2 := tm.GetToken(token2)
		if getErr2 != nil {
			t.Errorf("Second token should still be valid: %v", getErr2)
		}

		// 验证第三个token有效
		_, getErr3 := tm.GetToken(token3)
		if getErr3 != nil {
			t.Errorf("Third token should be valid: %v", getErr3)
		}

		// 清理（token1已被LRU策略删除）
		tm.DelToken(token2)
		tm.DelToken(token3)
	})

	t.Run("MultipleLoginRestriction", func(t *testing.T) {
		// 为同一用户添加第一个Token
		token1, err1 := tm.AddToken(100, 1, "192.168.1.10")
		if err1 != nil {
			t.Errorf("First token for user should be added successfully")
		}

		// 为同一用户添加第二个Token（应该成功，但会删除第一个token）
		_, err2 := tm.AddToken(100, 1, "192.168.1.11")
		if err2 != nil {
			t.Errorf("Second token should be added successfully: %v", err2)
		}

		// 验证第一个token是否被删除（不允许多设备登录时的预期行为）
		_, getErr := tm.GetToken(token1)
		if getErr == nil {
			t.Error("First token should be deleted when multiple login is disabled")
		} else {
			t.Logf("First token correctly deleted: %v", getErr)
		}
	})

	t.Run("EmptyAPIPath", func(t *testing.T) {
		token, err := tm.AddToken(200, 1, "192.168.1.20")
		if err != nil {
			t.Fatalf("Failed to add token: %v", err)
		}

		// 测试空API路径
		authResult := tm.Auth(token, "192.168.1.20", "")
		if authResult == nil {
			t.Error("Empty API path should be rejected")
		}

		// 清理
		tm.DelToken(token)
	})
}
