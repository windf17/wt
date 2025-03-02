package test

import (
	"testing"

	wtoken "github.com/windf17/wtoken"
)

func TestFrenchErrorMessages(t *testing.T) {
	// 1. 注册法语
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

	// 3. 初始化配置
	config := wtoken.Config{
		CacheFilePath: "test_token_fr.cache",
		Language:      "fr",
		MaxTokens:     1000,
		Debug:         true,
		Delimiter:     " ",
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

	// 5. 创建token管理器
	tokenManager := wtoken.InitTM[any](&config, groups, frenchErrorMessages)
	// 6. 测试无效用户ID场景
	_, errData := tokenManager.AddToken(0, 1, "192.168.1.100")
	if errData.Code != wtoken.ErrCodeInvalidUserID {
		t.Errorf("期望无效用户ID错误，但得到：%v", errData.Code)
	}
	if errData.Error() != "ID utilisateur invalide" {
		t.Errorf("期望法语错误消息，但得到：%v", errData.Error())
	}

	// 7. 测试无效token认证场景
	errData = tokenManager.Authenticate("invalid_token", "/api/user", "192.168.1.100")
	if errData.Code != wtoken.ErrCodeTokenNotFound {
		t.Errorf("期望token未找到错误，但得到：%v", errData.Code)
	}
	if errData.Error() != "Token introuvable" {
		t.Errorf("期望法语错误消息，但得到：%v", errData.Error())
	}

	// 8. 测试访问未授权API场景
	tokenKey, errData := tokenManager.AddToken(1001, 1, "192.168.1.100")
	if errData.Code != wtoken.ErrCodeSuccess {
		t.Fatalf("生成token失败：%v", errData.Error())
	}

	errData = tokenManager.Authenticate(tokenKey, "/api/admin", "192.168.1.100")
	if errData.Code != wtoken.ErrCodeAccessDenied {
		t.Errorf("期望访问拒绝错误，但得到：%v", errData.Code)
	}
	if errData.Error() != "Accès refusé" {
		t.Errorf("期望法语错误消息，但得到：%v", errData.Error())
	}

	// 9. 测试无效用户组场景
	_, errData = tokenManager.AddToken(1001, 999, "192.168.1.100")
	if errData.Code != wtoken.ErrCodeGroupNotFound {
		t.Errorf("期望用户组未找到错误，但得到：%v", errData.Code)
	}
	if errData.Error() != "Groupe introuvable" {
		t.Errorf("期望法语错误消息，但得到：%v", errData.Error())
	}

	// 10. 测试无效IP地址场景
	_, errData = tokenManager.AddToken(1001, 1, "")
	if errData.Code != wtoken.ErrCodeInvalidIP {
		t.Errorf("期望无效IP地址错误，但得到：%v", errData.Code)
	}
	if errData.Error() != "Adresse IP invalide" {
		t.Errorf("期望法语错误消息，但得到：%v", errData.Error())
	}
}
