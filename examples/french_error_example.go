package main

import (
	"fmt"

	wtoken "github.com/windf17/wtoken"
)

func TestFR() {
	// 1. 注册法语为有效的语言类型
	fr := wtoken.Language("fr")

	// 2. 定义法语错误提示信息
	frenchErrorMessages := map[wtoken.Language]map[wtoken.ErrorCode]string{
		fr: {
			wtoken.ErrSuccess:              "Opération réussie",
			wtoken.ErrInvalidToken:         "Token invalide",
			wtoken.ErrTokenNotFound:        "Token introuvable",
			wtoken.ErrTokenExpired:         "Token expiré",
			wtoken.ErrInvalidUserID:        "ID utilisateur invalide",
			wtoken.ErrInvalidGroupID:       "ID groupe invalide",
			wtoken.ErrInvalidIP:            "Adresse IP invalide",
			wtoken.ErrInvalidURL:           "URL invalide",
			wtoken.ErrAccessDenied:         "Accès refusé",
			wtoken.ErrGroupNotFound:        "Groupe introuvable",
			wtoken.ErrTokenLimitExceeded:   "Limite de token dépassée",
			wtoken.ErrCacheFileLoadFailed:  "Échec du chargement du fichier cache",
			wtoken.ErrCacheFileParseFailed: "Échec de l'analyse du fichier cache",
		},
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
	tokenManager := wtoken.InitTM[any](nil, groups, frenchErrorMessages)
	// 6. 测试错误信息（使用无效的用户ID生成token）
	_, errData := tokenManager.AddToken(0, 1, "192.168.1.100")
	// 应该返回法语的错误信息
	if errData != wtoken.ErrSuccess {
		fmt.Printf("Erreur: %v\n", errData.Error())
	}

	// 7. 尝试使用不存在的token进行认证
	errData = tokenManager.Authenticate("invalid_token", "/api/user", "192.168.1.100")
	// 应该返回法语的错误信息
	if errData != wtoken.ErrSuccess {
		fmt.Printf("Erreur d'authentification: %v\n", errData.Error())
	}

	// 8. 正确生成token
	tokenKey, errData := tokenManager.AddToken(1001, 1, "192.168.1.100")
	if errData == wtoken.ErrSuccess {
		fmt.Printf("Token généré avec succès: %s\n", tokenKey)
	}

	// 9. 使用正确的token进行认证
	errData = tokenManager.Authenticate(tokenKey, "/api/user", "192.168.1.100")
	if errData == wtoken.ErrSuccess {
		fmt.Println("Authentification réussie")
	}

	// 10. 尝试访问未授权的API
	errData = tokenManager.Authenticate(tokenKey, "/api/admin", "192.168.1.100")
	if errData != wtoken.ErrSuccess {
		fmt.Printf("Erreur d'accès: %v\n", errData.Error())
	}

	// 11. 保存和获取用户数据
	userData := "données utilisateur"
	if err := tokenManager.SaveData(tokenKey, userData); err == wtoken.ErrSuccess {
		fmt.Println("Données utilisateur enregistrées avec succès")
	}

	if loadedData, err := tokenManager.GetData(tokenKey); err == wtoken.ErrSuccess {
		fmt.Printf("Données utilisateur chargées: %v\n", loadedData)
	}

	// 12. 获取token统计信息
	stats := tokenManager.GetStats()
	fmt.Printf("Statistiques des tokens: Total=%d, Actif=%d\n", stats.TotalTokens, stats.ActiveTokens)

	// 13. 删除token
	errData = tokenManager.DelToken(tokenKey)
	if errData == wtoken.ErrSuccess {
		fmt.Println("Token supprimé avec succès")
	}
}
