package main

import (
	"fmt"

	wtoken "github.com/windf17/wtoken"
)

// UserInfo 定义用户信息结构体
type UserInfo struct {
	Username string
	Role     string
	Age      int
}

func Test() {
	// 1. 初始化token管理器配置
	config := wtoken.Config{
		CacheFilePath: "token_user.cache",
		Language:      "zh",
		MaxTokens:     1000,
		Debug:         true,
		Delimiter:     " ",
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
	tokenManager:= wtoken.InitTM[UserInfo](&config, groups, nil)

	// 4. 生成用户token
	tokenKey, errData := tokenManager.AddToken(1001, 1, "192.168.1.100")
	if errData.Code != wtoken.ErrCodeSuccess {
		fmt.Printf("生成token失败：%v\n", errData.Error())
		return
	}
	fmt.Printf("生成token成功：%s\n", tokenKey)

	// 5. 保存用户数据（使用UserInfo类型）
	userData := UserInfo{
		Username: "张三",
		Role:     "user",
		Age:      25,
	}
	if err := tokenManager.SaveData(tokenKey, userData); err == nil {
		fmt.Println("保存用户数据成功")
	}

	// 6. 获取用户数据（自动转换为UserInfo类型）
	if userData, err := tokenManager.GetData(tokenKey); err == nil {
		fmt.Printf("用户数据：用户名=%s, 角色=%s, 年龄=%d\n",
			userData.Username, userData.Role, userData.Age)
	}

	// 7. 更新用户数据
	userData.Role = "admin"
	userData.Age = 26
	if errSave := tokenManager.SaveData(tokenKey, userData); errSave == nil {
		fmt.Println("更新用户数据成功")
	}

	// 8. 再次获取更新后的数据
	if userData, err := tokenManager.GetData(tokenKey); err == nil {
		fmt.Printf("更新后的用户数据：用户名=%s, 角色=%s, 年龄=%d\n",
			userData.Username, userData.Role, userData.Age)
	}

	// 9. 删除token
	errData = tokenManager.DelToken(tokenKey)
	if errData.Code == wtoken.ErrCodeSuccess {
		fmt.Println("删除token成功")
	}
}
