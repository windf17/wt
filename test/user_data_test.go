package test

import (
	"os"
	"testing"

	wtoken "github.com/windf17/wtoken"
	"github.com/windf17/wtoken/models"
)

// UserInfo 定义用户信息结构体
type UserInfo struct {
	Username string
	Role     string
	Age      int
}

/**
 * TestUserDataOperations 测试用户数据相关操作
 */
func TestUserDataOperations(t *testing.T) {
	// 创建临时缓存文件
	tempFile := "test_user_cache.json"
	defer os.Remove(tempFile)

	// 初始化配置
	config := wtoken.ConfigRaw{

		MaxTokens:      1000,
		Delimiter:      "|",
		TokenRenewTime: "3600s",
		Language:       wtoken.LangChinese,
	}

	// 配置用户组
	groups := []models.GroupRaw{
		{
			ID:                 1,
			Name:               "user",
			AllowedAPIs:        "/api/user|/api/profile",
			DeniedAPIs:         "/api/admin",
			TokenExpire:        "3600s",
			AllowMultipleLogin: 0,
		},
	}

	// 创建指定UserInfo类型的token管理器
	tokenManager := wtoken.InitTM[UserInfo](&config, groups, nil)
	if tokenManager == nil {
		t.Fatalf("Failed to initialize token manager")
	}
	defer tokenManager.Close()

	// 测试生成用户token
	tokenKey, errData := tokenManager.AddToken(1001, 1, "192.168.1.100")
	if errData != wtoken.E_Success {
		t.Fatalf("生成token失败：%v", errData)
	}
	t.Logf("生成token成功：%s", tokenKey)

	// 测试保存用户数据
	userData := UserInfo{
		Username: "张三",
		Role:     "user",
		Age:      25,
	}
	err := tokenManager.SetUserData(tokenKey, userData)
	if err != wtoken.E_Success {
		t.Fatalf("保存用户数据失败：%v", err)
	}
	t.Log("保存用户数据成功")

	// 测试获取用户数据
	retrievedData, err := tokenManager.GetUserData(tokenKey)
	if err != wtoken.E_Success {
		t.Fatalf("获取用户数据失败：%v", err)
	}

	// 验证获取的数据是否正确
	if retrievedData.Username != userData.Username ||
		retrievedData.Role != userData.Role ||
		retrievedData.Age != userData.Age {
		t.Errorf("获取的用户数据与保存的不匹配，期望：%+v，实际：%+v",
			userData, retrievedData)
	}

	// 测试更新用户数据
	userData.Role = "admin"
	userData.Age = 26
	err = tokenManager.SetUserData(tokenKey, userData)
	if err != wtoken.E_Success {
		t.Fatalf("更新用户数据失败：%v", err)
	}
	t.Log("更新用户数据成功")

	// 测试获取更新后的数据
	updatedData, err := tokenManager.GetUserData(tokenKey)
	if err != wtoken.E_Success {
		t.Fatalf("获取更新后的用户数据失败：%v", err)
	}

	// 验证更新后的数据是否正确
	if updatedData.Role != "admin" || updatedData.Age != 26 {
		t.Errorf("更新后的用户数据不正确，期望Role=admin,Age=26，实际：Role=%s,Age=%d",
			updatedData.Role, updatedData.Age)
	}

	// 测试删除token
	errData = tokenManager.DelToken(tokenKey)
	if errData != wtoken.E_Success {
		t.Fatalf("删除token失败：%v", errData)
	}
	t.Log("删除token成功")

	// 验证token已被删除
	_, errData = tokenManager.GetToken(tokenKey)
	if errData == wtoken.E_Success {
		t.Errorf("期望token不存在，但获取到了token")
	}
}

/**
 * TestUserDataErrorCases 测试用户数据操作的错误情况
 */
func TestUserDataErrorCases(t *testing.T) {
	// 创建临时缓存文件
	tempFile := "test_user_error_cache.json"
	defer os.Remove(tempFile)

	// 配置用户组
	groups := []models.GroupRaw{
		{
			ID:                 1,
			Name:               "user",
			AllowedAPIs:        "/api/user",
			TokenExpire:        "3600s",
			AllowMultipleLogin: 0,
		},
	}

	// 创建token管理器
	tokenManager := wtoken.InitTM[UserInfo](&wtoken.ConfigRaw{

		MaxTokens:      1000,
		Delimiter:      "|",
		TokenRenewTime: "3600s",
		Language:       wtoken.LangChinese,
	}, groups, nil)
	if tokenManager == nil {
		t.Fatalf("Failed to initialize token manager")
	}
	defer tokenManager.Close()

	// 测试无效的token
	_, err := tokenManager.GetUserData("invalid_token")
	if err == wtoken.E_Success {
		t.Error("期望获取无效token数据失败，但成功了")
	}

	// 测试使用无效的用户组ID
	_, errData := tokenManager.AddToken(1001, 999, "192.168.1.100")
	if errData == wtoken.E_Success {
		t.Error("期望使用无效的用户组ID失败，但成功了")
	}

	// 测试使用无效的IP地址
	_, errData = tokenManager.AddToken(1001, 1, "")
	if errData == wtoken.E_Success {
		t.Error("期望使用空IP地址失败，但成功了")
	}

	// 测试使用无效的用户ID
	_, errData = tokenManager.AddToken(0, 1, "192.168.1.100")
	if errData == wtoken.E_Success {
		t.Error("期望使用无效的用户ID失败，但成功了")
	}
}

/**
 * TestUserDataWithDifferentTypes 测试不同类型的用户数据
 */
func TestUserDataWithDifferentTypes(t *testing.T) {
	// 创建临时缓存文件
	tempFile := "test_user_types_cache.json"
	defer os.Remove(tempFile)

	// 配置用户组
	groups := []models.GroupRaw{
		{
			ID:                 1,
			Name:               "user",
			AllowedAPIs:        "/api/user",
			TokenExpire:        "3600s",
			AllowMultipleLogin: 1,
		},
	}

	// 测试map[string]any类型
	t.Run("MapStringInterface", func(t *testing.T) {
		tm := wtoken.InitTM[map[string]any](&wtoken.ConfigRaw{

			MaxTokens:      100,
			Delimiter:      "|",
			TokenRenewTime: "3600s",
			Language:       wtoken.LangChinese,
		}, groups, nil)
		if tm == nil {
			t.Fatalf("Failed to initialize token manager")
		}
		defer tm.Close()
		defer os.Remove(tempFile + "_map")

		token, err := tm.AddToken(1001, 1, "192.168.1.100")
		if err != wtoken.E_Success {
			t.Fatalf("Failed to add token: %v", err)
		}

		userData := map[string]any{
			"username": "test_user",
			"age":      30,
			"active":   true,
		}

		err = tm.SetUserData(token, userData)
		if err != wtoken.E_Success {
			t.Fatalf("Failed to set user data: %v", err)
		}

		retrievedData, err := tm.GetUserData(token)
		if err != wtoken.E_Success {
			t.Fatalf("Failed to get user data: %v", err)
		}

		if retrievedData["username"] != "test_user" {
			t.Errorf("Expected username 'test_user', got %v", retrievedData["username"])
		}
	})

	// 测试string类型
	t.Run("StringType", func(t *testing.T) {
		tm := wtoken.InitTM[string](&wtoken.ConfigRaw{

			MaxTokens:      100,
			Delimiter:      "|",
			TokenRenewTime: "3600s",
			Language:       wtoken.LangChinese,
		}, groups, nil)
		if tm == nil {
			t.Fatalf("Failed to initialize token manager")
		}
		defer tm.Close()
		defer os.Remove(tempFile + "_string")

		token, err := tm.AddToken(1002, 1, "192.168.1.101")
		if err != wtoken.E_Success {
			t.Fatalf("Failed to add token: %v", err)
		}

		userData := "simple_user_data"
		err = tm.SetUserData(token, userData)
		if err != wtoken.E_Success {
			t.Fatalf("Failed to set user data: %v", err)
		}

		retrievedData, err := tm.GetUserData(token)
		if err != wtoken.E_Success {
			t.Fatalf("Failed to get user data: %v", err)
		}

		if retrievedData != "simple_user_data" {
			t.Errorf("Expected 'simple_user_data', got %v", retrievedData)
		}
	})
}
