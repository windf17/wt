package test

import (
	"os"
	"testing"

	wtoken "github.com/windf17/wtoken"
)

// UserInfo 定义用户信息结构体
type UserInfo struct {
	Username string
	Role     string
	Age      int
}

// TestUserDataOperations 测试用户数据相关操作
func TestUserDataOperations(t *testing.T) {
	// 1. 初始化配置
	config := wtoken.Config{
		CacheFilePath:   "./cache_tmp.json",
	}

	// 2. 配置用户组
	groups := []wtoken.GroupRaw{
		{
			ID:                 1,
			AllowedAPIs:        "/api/user /api/profile",
			DeniedAPIs:         "/api/admin",
			TokenExpire:        "3600",
			AllowMultipleLogin: 0,
		},
	}

	// 3. 创建指定UserInfo类型的token管理器
	tokenManager := wtoken.InitTM[UserInfo](nil, groups, nil)

	// 4. 测试生成用户token
	tokenKey, errData := tokenManager.AddToken(1001, 1, "192.168.1.100")
	if errData != wtoken.ErrSuccess {
		t.Fatalf("生成token失败：%v", errData.Error())
	}
	t.Logf("生成token成功：%s", tokenKey)

	// 5. 测试保存用户数据
	userData := UserInfo{
		Username: "张三",
		Role:     "user",
		Age:      25,
	}
	err := tokenManager.SaveData(tokenKey, userData)
	if err != wtoken.ErrSuccess {
		t.Fatalf("保存用户数据失败：%v", err)
	}
	t.Log("保存用户数据成功")

	// 6. 测试获取用户数据
	retrievedData, err := tokenManager.GetData(tokenKey)
	if err != wtoken.ErrSuccess {
		t.Fatalf("获取用户数据失败：%v", err)
	}

	// 验证获取的数据是否正确
	if retrievedData.Username != userData.Username ||
		retrievedData.Role != userData.Role ||
		retrievedData.Age != userData.Age {
		t.Errorf("获取的用户数据与保存的不匹配，期望：%+v，实际：%+v",
			userData, retrievedData)
	}

	// 7. 测试更新用户数据
	userData.Role = "admin"
	userData.Age = 26
	err = tokenManager.SaveData(tokenKey, userData)
	if err != wtoken.ErrSuccess {
		t.Fatalf("更新用户数据失败：%v", err)
	}
	t.Log("更新用户数据成功")

	// 8. 测试获取更新后的数据
	updatedData, err := tokenManager.GetData(tokenKey)
	if err != wtoken.ErrSuccess {
		t.Fatalf("获取更新后的用户数据失败：%v", err)
	}

	// 验证更新后的数据是否正确
	if updatedData.Role != "admin" || updatedData.Age != 26 {
		t.Errorf("更新后的用户数据不正确，期望Role=admin,Age=26，实际：Role=%s,Age=%d",
			updatedData.Role, updatedData.Age)
	}

	// 9. 测试删除token
	errData = tokenManager.DelToken(tokenKey)
	if errData != wtoken.ErrSuccess {
		t.Fatalf("删除token失败：%v", errData.Error())
	}
	t.Log("删除token成功")

	// 10. 验证token已被删除
	_, errData = tokenManager.GetToken(tokenKey)
	if errData != wtoken.ErrTokenNotFound {
		t.Errorf("期望token不存在，但获取到了token")
	}

	// 清理测试文件
	os.Remove(config.CacheFilePath)
}

// TestUserDataErrorCases 测试用户数据操作的错误情况
func TestUserDataErrorCases(t *testing.T) {
	// 1. 初始化配置

	// 2. 配置用户组
	groups := []wtoken.GroupRaw{
		{
			ID:                 1,
			AllowedAPIs:        "/api/user",
			TokenExpire:        "3600",
			AllowMultipleLogin: 0,
		},
	}

	// 3. 创建token管理器
	tokenManager := wtoken.InitTM[UserInfo](nil, groups, nil)

	// 4. 测试无效的token
	_, err := tokenManager.GetData("invalid_token")
	if err == wtoken.ErrSuccess {
		t.Error("期望获取无效token数据失败，但成功了")
	}

	// 5. 测试使用无效的用户组ID
	_, errData := tokenManager.AddToken(1001, 999, "192.168.1.100")
	if errData == wtoken.ErrSuccess {
		t.Error("期望使用无效的用户组ID失败，但成功了")
	}

	// 6. 测试使用无效的IP地址
	_, errData = tokenManager.AddToken(1001, 1, "")
	if errData == wtoken.ErrSuccess {
		t.Error("期望使用空IP地址失败，但成功了")
	}

	// 7. 测试使用无效的用户ID
	_, errData = tokenManager.AddToken(0, 1, "192.168.1.100")
	if errData == wtoken.ErrSuccess {
		t.Error("期望使用无效的用户ID失败，但成功了")
	}

	// 清理测试文件
	os.Remove(wtoken.DefaultConfigRaw.CacheFilePath)
}
