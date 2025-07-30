package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	wtoken "github.com/windf17/wtoken"
	"github.com/windf17/wtoken/models"
)

/**
 * TestBackupFileLoading 测试从备份文件加载数据功能
 */
func TestBackupFileLoading(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	backupFile := filepath.Join(tempDir, "test_backup.json")

	// 准备测试数据 - 使用正确的Token结构
	testData := map[string]any{
		"tokens": map[string]any{
			"test_token_1": map[string]any{
				"UserID":         uint(1),
				"GroupID":        uint(1),
				"IP":             "192.168.1.1",
				"UserData":       "test data 1",
				"LoginTime":      time.Now().Format(time.RFC3339),
				"LastAccessTime": time.Now().Format(time.RFC3339),
				"ExpireSeconds":  int64(3600),
			},
			"test_token_2": map[string]any{
				"UserID":         uint(2),
				"GroupID":        uint(1),
				"IP":             "192.168.1.2",
				"UserData":       "test data 2",
				"LoginTime":      time.Now().Format(time.RFC3339),
				"LastAccessTime": time.Now().Format(time.RFC3339),
				"ExpireSeconds":  int64(3600),
			},
		},
		"groups": map[string]any{
			"1": map[string]any{
				"Name":               "test_group",
				"AllowedAPIs":        []string{"/api/*"},
				"DeniedAPIs":         []string{},
				"ExpireSeconds":      int64(3600),
				"AllowMultipleLogin": true,
			},
		},
		"stats": map[string]any{
			"TotalTokens":   2,
			"ActiveTokens":  2,
			"ExpiredTokens": 0,
			"DeletedTokens": 0,
			"LastCleanTime": time.Now().Format(time.RFC3339),
		},
	}

	// 将测试数据写入备份文件
	data, err := json.Marshal(testData)
	assert.NoError(t, err)
	err = os.WriteFile(backupFile, data, 0644)
	assert.NoError(t, err)

	// 测试1: 配置了备份文件路径时应该加载数据
	t.Run("LoadFromBackupFile", func(t *testing.T) {
		config := &wtoken.ConfigRaw{
			Language:       "zh",
			MaxTokens:      1000,
			Delimiter:      "|",
			TokenRenewTime: "24h",
		}

		groups := []models.GroupRaw{
			{
				ID:                 1,
				Name:               "test_group",
				AllowedAPIs:        "/api/*",
				DeniedAPIs:         "",
				TokenExpire:        "1h",
				AllowMultipleLogin: 1,
			},
		}

		tm := wtoken.InitTM[string](config, groups, nil)
		defer tm.Close()

		// 验证tokens是否正确加载
		token1, code := tm.GetToken("test_token_1")
		if code == wtoken.E_Success {
			assert.NotNil(t, token1)
			assert.Equal(t, uint(1), token1.UserID)
			assert.Equal(t, uint(1), token1.GroupID)
			assert.Equal(t, "192.168.1.1", token1.IP)
			assert.Equal(t, "test data 1", token1.UserData)
		}

		// 验证groups是否正确加载
		group, code := tm.GetGroup(1)
		if code == wtoken.E_Success {
			assert.NotNil(t, group)
			assert.Equal(t, "test_group", group.Name)
		}
	})

	// 测试2: 未配置备份文件路径时不应该加载数据
	t.Run("NoBackupFileConfigured", func(t *testing.T) {
		config := &wtoken.ConfigRaw{
			Language:       wtoken.LangChinese,
			MaxTokens:      1000,
			Delimiter:      "|",
			TokenRenewTime: "24h",
		}

		groups := []models.GroupRaw{
			{
				ID:                 1,
				Name:               "test_group",
				AllowedAPIs:        "/api/*",
				DeniedAPIs:         "",
				TokenExpire:        "1h",
				AllowMultipleLogin: 1,
			},
		}

		tm := wtoken.InitTM[string](config, groups, nil)
		defer tm.Close()

		// 验证没有加载任何数据
		_, code := tm.GetToken("test_token_1")
		assert.Equal(t, wtoken.E_InvalidToken, code)

		_, code = tm.GetToken("test_token_2")
		assert.Equal(t, wtoken.E_InvalidToken, code)

		// 验证用户组配置仍然有效（从groups参数加载）
		group, code := tm.GetGroup(1)
		assert.Equal(t, wtoken.E_Success, code)
		assert.NotNil(t, group)
	})

	// 测试3: 备份文件不存在时应该正常启动
	t.Run("BackupFileNotExists", func(t *testing.T) { // 测试不存在的备份文件
		config := &wtoken.ConfigRaw{
			Language:       "zh",
			MaxTokens:      1000,
			Delimiter:      "|",
			TokenRenewTime: "24h",
		}

		groups := []models.GroupRaw{
			{
				ID:                 1,
				Name:               "test_group",
				AllowedAPIs:        "/api/*",
				DeniedAPIs:         "",
				TokenExpire:        "1h",
				AllowMultipleLogin: 1,
			},
		}

		tm := wtoken.InitTM[string](config, groups, nil)
		defer tm.Close()

		// 验证没有加载任何数据，但管理器正常工作
		_, code := tm.GetToken("test_token_1")
		assert.Equal(t, wtoken.E_InvalidToken, code)

		// 验证可以正常添加新token
		newToken, code := tm.AddToken(1, 1, "192.168.1.100")
		assert.Equal(t, wtoken.E_Success, code)
		assert.NotEmpty(t, newToken)

		token, code := tm.GetToken(newToken)
		assert.Equal(t, wtoken.E_Success, code)
		assert.NotNil(t, token)
		assert.Equal(t, uint(1), token.UserID)
	})
}

/**
 * TestBackupFileLoadingWithInvalidData 测试加载无效备份文件数据
 */
func TestBackupFileLoadingWithInvalidData(t *testing.T) {
	tempDir := t.TempDir()
	backupFile := filepath.Join(tempDir, "invalid_backup.json")

	// 写入无效的JSON数据
	err := os.WriteFile(backupFile, []byte("invalid json data"), 0644)
	assert.NoError(t, err)

	config := &wtoken.ConfigRaw{
		Language:       "zh",
		MaxTokens:      1000,
		Delimiter:      "|",
		TokenRenewTime: "24h",
	}

	groups := []models.GroupRaw{
		{
			ID:                 1,
			Name:               "test_group",
			AllowedAPIs:        "/api/*",
			DeniedAPIs:         "",
			TokenExpire:        "1h",
			AllowMultipleLogin: 1,
		},
	}

	// 应该能正常启动，但不会加载任何数据
	tm := wtoken.InitTM[string](config, groups, nil)
	defer tm.Close()

	// 验证没有加载任何数据
	_, code := tm.GetToken("any_token")
	assert.Equal(t, wtoken.E_InvalidToken, code)

	// 验证管理器仍然可以正常工作
	newToken, code := tm.AddToken(1, 1, "192.168.1.100")
	assert.Equal(t, wtoken.E_Success, code)
	assert.NotEmpty(t, newToken)

	// 验证新添加的token可以正常获取
	token, code := tm.GetToken(newToken)
	assert.Equal(t, wtoken.E_Success, code)
	assert.NotNil(t, token)
	assert.Equal(t, uint(1), token.UserID)
}
