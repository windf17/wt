package main

import (
	"fmt"

	"github.com/windf17/wtoken"
)

func main() {
	// 1. 初始化token管理器配置
	config := wtoken.Config{
		CacheFilePath: "token.cache", // token缓存文件路径
		Language:      "zh",          // 错误信息语言
		MaxTokens:     1000,          // 最大token数量
		Debug:         true,          // 是否开启调试模式
		Delimiter:     " ",           // API分隔符
	}

	// 2. 配置用户组
	groups := []wtoken.GroupRaw{
		{
			ID:                 1,
			AllowedAPIs:        "/api/user /api/product",
			DeniedAPIs:         "/api/admin",
			TokenExpire:        "3600",
			AllowMultipleLogin: 0,
		},
		{
			ID:                 2,
			AllowedAPIs:        "/api/admin /api/user /api/product",
			DeniedAPIs:         "",
			TokenExpire:        "7200",
			AllowMultipleLogin: 1,
		},
	}

	// 3. 创建token管理器
	tokenManager, err := wtoken.InitTM[any](&config, groups, nil)
	if err != nil {
		fmt.Printf("初始化token管理器失败：%v\n", err)
		return
	}

	// 4. 生成用户token
	// userData := map[string]interface{}{
	// 	"username": "张三",
	// 	"role":     "user",
	// }
	tokenKey, errData := tokenManager.AddToken(1001, 1, "192.168.1.100")
	if errData.Code != wtoken.ErrCodeSuccess {
		fmt.Printf("生成token失败：%v\n", errData.Error())
		return
	}
	fmt.Printf("生成token成功：%s\n", tokenKey)

	// 5. API鉴权测试
	// 5.1 允许访问的API
	errData = tokenManager.Authenticate(tokenKey, "/api/user", "192.168.1.100")
	if errData.Code == wtoken.ErrCodeSuccess {
		fmt.Println("访问/api/user鉴权成功")
	} else {
		fmt.Printf("访问/api/user鉴权失败：%v\n", errData.Error())
	}

	// 5.2 禁止访问的API
	errData = tokenManager.Authenticate(tokenKey, "/api/admin", "192.168.1.100")
	if errData.Code != wtoken.ErrCodeSuccess {
		fmt.Printf("访问/api/admin鉴权失败（预期结果）：%v\n", errData.Error())
	}

	// 6. 获取token信息
	token, errData := tokenManager.GetToken(tokenKey)
	if errData.Code == wtoken.ErrCodeSuccess {
		fmt.Printf("token信息：%+v\n", token)
	}

	// 7. 获取用户数据
	userInfo, err := tokenManager.GetData(tokenKey)
	if err == nil {
		fmt.Printf("用户数据：%+v\n", userInfo)
	}

	// 8. 更新用户数据
	userInfo = "admin"
	if err = tokenManager.SaveData(tokenKey, userInfo); err == nil {
		fmt.Println("更新用户数据成功")
	}

	// 9. 删除token
	errData = tokenManager.DelToken(tokenKey)
	if errData.Code == wtoken.ErrCodeSuccess {
		fmt.Println("删除token成功")
	}

	// 10. 获取统计信息
	stats := tokenManager.GetStats()
	fmt.Printf("统计信息：%+v\n", stats)
}
