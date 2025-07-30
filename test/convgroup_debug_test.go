package test

import (
	"os"
	"strings"
	"testing"

	"github.com/windf17/wt"
	"github.com/windf17/wt/models"
	"github.com/windf17/wt/utility"
)

/**
 * TestConvGroupDebug 专门调试ConvGroup函数的行为
 * 分析为什么/api/user/profile会被解析为/api/user/*
 */
func TestConvGroupDebug(t *testing.T) {
	// 测试单独的路径解析
	t.Run("TestPathParsing", func(t *testing.T) {
		testPaths := []string{
			"/api/user/profile",
			"/api/user/admin",
			"/api/admin,/api/user,/api/admin/users",
			"/api/admin/delete",
		}

		for _, path := range testPaths {
			segments := utility.ParsePathToSegments(path)
			t.Logf("Path '%s' -> Segments: %v", path, segments)
		}
	})

	// 测试分隔符分割
	t.Run("TestDelimiterSplit", func(t *testing.T) {
		allowedAPIs := "/api/user/profile"
		deniedAPIs := "/api/user/admin"
		delimiter := ","

		t.Logf("Original strings:")
		t.Logf("  AllowedAPIs: '%s'", allowedAPIs)
		t.Logf("  DeniedAPIs: '%s'", deniedAPIs)
		t.Logf("  Delimiter: '%s'", delimiter)

		// 手动分割测试
		allowedSplit := strings.Split(allowedAPIs, delimiter)
		deniedSplit := strings.Split(deniedAPIs, delimiter)

		t.Logf("\nSplit results:")
		t.Logf("  AllowedAPIs split: %v", allowedSplit)
		t.Logf("  DeniedAPIs split: %v", deniedSplit)

		// 测试每个分割后的API
		for i, api := range allowedSplit {
			trimmed := strings.TrimSpace(api)
			segments := utility.ParsePathToSegments(trimmed)
			t.Logf("  Allowed[%d]: '%s' -> trimmed: '%s' -> segments: %v", i, api, trimmed, segments)
		}

		for i, api := range deniedSplit {
			trimmed := strings.TrimSpace(api)
			segments := utility.ParsePathToSegments(trimmed)
			t.Logf("  Denied[%d]: '%s' -> trimmed: '%s' -> segments: %v", i, api, trimmed, segments)
		}
	})

	// 测试ConvGroup函数的完整流程
	t.Run("TestConvGroupFlow", func(t *testing.T) {
		raw := models.GroupRaw{
			ID:                 2,
			Name:               "user",
			AllowedAPIs:        "/api/user/profile",
			DeniedAPIs:         "/api/user/admin",
			TokenExpire:        "1h",
			AllowMultipleLogin: 1,
		}

		t.Logf("Input GroupRaw:")
		t.Logf("  ID: %d", raw.ID)
		t.Logf("  Name: %s", raw.Name)
		t.Logf("  AllowedAPIs: '%s'", raw.AllowedAPIs)
		t.Logf("  DeniedAPIs: '%s'", raw.DeniedAPIs)
		t.Logf("  TokenExpire: %s", raw.TokenExpire)
		t.Logf("  AllowMultipleLogin: %d", raw.AllowMultipleLogin)

		delimiter := ","
		convertedGroup := wt.ConvGroup(raw, delimiter)

		t.Logf("\nOutput Group:")
		t.Logf("  Name: %s", convertedGroup.Name)
		t.Logf("  ExpireSeconds: %d", convertedGroup.ExpireSeconds)
		t.Logf("  AllowMultipleLogin: %v", convertedGroup.AllowMultipleLogin)
		t.Logf("  ApiRules count: %d", len(convertedGroup.ApiRules))

		for i, rule := range convertedGroup.ApiRules {
			t.Logf("  Rule %d: Path=%v, Rule=%v", i, rule.Path, rule.Rule)
		}
	})

	// 测试与core_functionality_test.go完全相同的配置
	t.Run("TestExactCoreConfig", func(t *testing.T) {
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

		delimiter := ","

		for _, raw := range groups {
			t.Logf("\n=== Processing Group %d ===", raw.ID)
			t.Logf("Input: AllowedAPIs='%s', DeniedAPIs='%s'", raw.AllowedAPIs, raw.DeniedAPIs)

			convertedGroup := wt.ConvGroup(raw, delimiter)

			t.Logf("Output: %d ApiRules", len(convertedGroup.ApiRules))
			for i, rule := range convertedGroup.ApiRules {
				t.Logf("  Rule %d: Path=%v, Rule=%v", i, rule.Path, rule.Rule)
			}
		}
	})

	// 测试缓存文件的影响
	t.Run("TestCacheFileImpact", func(t *testing.T) {
		// 创建临时缓存文件（与core_functionality_test.go完全相同）
		tempFile := "test_cache.json"
		cacheFile := tempFile + ".cache"

		t.Logf("Cache file path: %s", cacheFile)

		// 检查缓存文件是否存在
		if _, err := os.Stat(cacheFile); err == nil {
			t.Logf("Cache file exists")
			// 读取缓存文件内容
			var data []byte
			if data, err = os.ReadFile(cacheFile); err == nil {
				t.Logf("Cache file content (first 500 chars): %s", string(data[:min(len(data), 500)]))
			} else {
				t.Logf("Failed to read cache file: %v", err)
			}
		} else {
			t.Logf("Cache file does not exist: %v", err)
		}
	})
}

// min 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
