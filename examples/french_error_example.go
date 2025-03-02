package main

import (
	"fmt"

	wtoken "github.com/windf17/wtoken"
)

func TestFR() {
	// 1. 注册法语为有效的语言类型
	fr := wtoken.RegisterLanguage("fr")

	// 2. 定义法语错误提示信息
	frenchErrorMessages := map[wtoken.ILanguage]map[wtoken.ErrorCode]string{
		fr: {
			wtoken.ErrCodeSuccess:              "Opération réussie",
			wtoken.ErrCodeInvalidToken:         "Token invalide",
			wtoken.ErrCodeTokenNotFound:        "Token introuvable",
			wtoken.ErrCodeTokenExpired:         "Token expiré",
			wtoken.ErrCodeInvalidUserID:        "ID utilisateur invalide",
			wtoken.ErrCodeInvalidGroupID:       "ID groupe invalide",
			wtoken.ErrCodeInvalidIP:            "Adresse IP invalide",
			wtoken.ErrCodeInvalidURL:           "URL invalide",
			wtoken.ErrCodeAccessDenied:         "Accès refusé",
			wtoken.ErrCodeGroupNotFound:        "Groupe introuvable",
			wtoken.ErrCodeAddToken:             "Échec de l'ajout du token",
			wtoken.ErrCodeCacheFileLoadFailed:  "Échec du chargement du fichier cache",
			wtoken.ErrCodeCacheFileParseFailed: "Échec de l'analyse du fichier cache",
		},
	}

	// 3. 初始化token管理器配置
	config := wtoken.Config{
		CacheFilePath: "token.cache", // token缓存文件路径
		Language:      "fr",          // 设置错误信息语言为法语
		MaxTokens:     1000,          // 最大token数量
		Debug:         true,          // 开启调试模式
		Delimiter:     " ",           // API分隔符
	}

	// 4. 配置用户组
	groups := []wtoken.GroupRaw{
		{
			ID:                 1,
			AllowedAPIs:        "/api/user /api/product",
			DeniedAPIs:         "/api/admin",
			TokenExpire:        "3600",
			AllowMultipleLogin: 0,
		},
	}
	// 5. 创建token管理器（使用自定义的法语错误信息）
	tokenManager := wtoken.InitTM[any](&config, groups, frenchErrorMessages)
	// 6. 测试错误信息（使用无效的用户ID生成token）
	_, errData := tokenManager.AddToken(0, 1, "192.168.1.100")
	// 应该返回法语的错误信息
	if errData.Code != wtoken.ErrCodeSuccess {
		fmt.Printf("Erreur: %v\n", errData.Error())
	}

	// 7. 尝试使用不存在的token进行认证
	errData = tokenManager.Authenticate("invalid_token", "/api/user", "192.168.1.100")
	// 应该返回法语的错误信息
	if errData.Code != wtoken.ErrCodeSuccess {
		fmt.Printf("Erreur d'authentification: %v\n", errData.Error())
	}

	// 8. 正确生成token
	tokenKey, errData := tokenManager.AddToken(1001, 1, "192.168.1.100")
	if errData.Code == wtoken.ErrCodeSuccess {
		fmt.Printf("Token généré avec succès: %s\n", tokenKey)
	}

	// 9. 使用正确的token进行认证
	errData = tokenManager.Authenticate(tokenKey, "/api/user", "192.168.1.100")
	if errData.Code == wtoken.ErrCodeSuccess {
		fmt.Println("Authentification réussie")
	}

	// 10. 尝试访问未授权的API
	errData = tokenManager.Authenticate(tokenKey, "/api/admin", "192.168.1.100")
	if errData.Code != wtoken.ErrCodeSuccess {
		fmt.Printf("Erreur d'accès: %v\n", errData.Error())
	}

	// 11. 删除token
	errData = tokenManager.DelToken(tokenKey)
	if errData.Code == wtoken.ErrCodeSuccess {
		fmt.Println("Token supprimé avec succès")
	}
}
