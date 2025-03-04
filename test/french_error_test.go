package test

import (
	"testing"

	wtoken "github.com/windf17/wtoken"
)

func TestFrenchErrorMessages(t *testing.T) {
	// 清理测试环境
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("测试过程中发生panic：%v", r)
		}
	}()
	// 1. 注册法语
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
	config := wtoken.DefaultConfigRaw
	config.Language = fr
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
	tokenManager := wtoken.InitTM[any](config, groups, frenchErrorMessages)
	// 6. 测试无效用户ID场景
	_, errData := tokenManager.AddToken(0, 1, "192.168.1.100")
	if errData != wtoken.ErrInvalidUserID {
		t.Errorf("期望无效用户ID错误，但得到：%v", errData.Error())
	}
	if errData.Error() != "ID utilisateur invalide" {
		t.Errorf("期望法语错误消息，但得到：%v", errData.Error())
	}

	// 7. 测试无效token认证场景
	errData = tokenManager.Authenticate("invalid_token", "/api/user", "192.168.1.100")
	if errData != wtoken.ErrTokenNotFound {
		t.Errorf("期望token未找到错误，但得到：%v", errData)
	}
	if errData.Error() != "Token introuvable" {
		t.Errorf("期望法语错误消息，但得到：%v", errData.Error())
	}

	// 8. 测试访问未授权API场景
	tokenKey, errData := tokenManager.AddToken(1001, 1, "192.168.1.100")
	if errData != wtoken.ErrSuccess {
		t.Fatalf("生成token失败：%v", errData.Error())
	}

	errData = tokenManager.Authenticate(tokenKey, "/api/admin", "192.168.1.100")
	if errData != wtoken.ErrAccessDenied {
		t.Errorf("期望访问拒绝错误，但得到：%v", errData)
	}
	if errData.Error() != "Accès refusé" {
		t.Errorf("期望法语错误消息，但得到：%v", errData.Error())
	}

	// 9. 测试无效用户组场景
	_, errData = tokenManager.AddToken(1001, 999, "192.168.1.100")
	if errData != wtoken.ErrGroupNotFound {
		t.Errorf("期望用户组未找到错误，但得到：%v", errData)
	}
	if errData.Error() != "Groupe introuvable" {
		t.Errorf("期望法语错误消息，但得到：%v", errData.Error())
	}

	// 10. 测试无效IP地址场景
	_, errData = tokenManager.AddToken(1001, 1, "")
	if errData != wtoken.ErrInvalidIP {
		t.Errorf("期望无效IP地址错误，但得到：%v", errData)
	}
	if errData.Error() != "Adresse IP invalide" {
		t.Errorf("期望法语错误消息，但得到：%v", errData.Error())
	}

	// 11. 测试token统计信息
	stats := tokenManager.GetStats()
	if stats.TotalTokens == 0 {
		t.Error("期望统计信息中有token记录，但得到0")
	}

	// 12. 测试用户数据操作
	userData := "test_data"
	if err := tokenManager.SaveData(tokenKey, userData); err != wtoken.ErrSuccess {
		t.Errorf("保存用户数据失败：%v", err)
	}

	loadedData, err := tokenManager.GetData(tokenKey)
	if err != wtoken.ErrSuccess {
		t.Errorf("获取用户数据失败：%v", err)
	}
	if loadedData != userData {
		t.Errorf("期望用户数据为%v，但得到：%v", userData, loadedData)
	}
}
