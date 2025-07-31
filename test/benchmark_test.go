package test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/windf17/wt"
	"github.com/windf17/wt/models"
)

/**
 * BenchmarkTokenOperations 基准测试token操作
 */
func BenchmarkTokenOperations(b *testing.B) {
	tm, err := wt.InitTM[any](models.ConfigRaw{
		MaxTokens:      1000,
		Delimiter:      "|",
		TokenRenewTime: "3600s",
		Language:       "zh",
	}, []models.GroupRaw{})
	if err != nil {
		b.Fatalf("Failed to initialize token manager: %v", err)
	}

	// 确保用户组存在
	group := &models.GroupRaw{
		ID:                 1,
		Name:               "test_group",
		AllowedAPIs:        "/api/*",
		DeniedAPIs:         "",
		TokenExpire:        "1h",
		AllowMultipleLogin: 1,
	}
	err2 := tm.AddGroup(group)
	if err2 != nil {
		b.Fatalf("Failed to add group: %v", err2)
	}

	// 预热
	for i := 0; i < 100; i++ {
		tokenKey, _ := tm.AddToken(uint(i+1), 1, "192.168.1.1")
		_, _ = tm.GetToken(tokenKey)
		tm.DelToken(tokenKey)
	}

	b.Run("AddToken", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := uint(i%1000 + 1) // 循环使用1-1000的用户ID
			ip := fmt.Sprintf("192.168.%d.%d", (i%254)+1, (i%254)+1)
			_, err := tm.AddToken(userID, 1, ip)
			if err != nil {
				b.Errorf("AddToken failed: %v", err)
			}
		}
	})

	b.Run("GetToken", func(b *testing.B) {
		// 预先添加一些token
		tokens := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			tokenKey, _ := tm.AddToken(uint(i+1), 1, "192.168.1.1")
			tokens[i] = tokenKey
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tokenKey := tokens[i%1000]
			_, _ = tm.GetToken(tokenKey)
		}
	})

	b.Run("Auth", func(b *testing.B) {
		// 预先添加一些token
		tokens := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			tokenKey, _ := tm.AddToken(uint(i+1), 1, "192.168.1.1")
			tokens[i] = tokenKey
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tokenKey := tokens[i%1000]
			_ = tm.Auth(tokenKey, "192.168.1.1", "/api/test")
		}
	})

	b.Run("DelToken", func(b *testing.B) {
		// 预先创建一些token
		tokens := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			tokenKey, err := tm.AddToken(uint(i+1), 1, "192.168.1.1")
			if err != nil {
				b.Fatalf("Failed to add token: %v", err)
			}
			tokens[i] = tokenKey
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// 循环使用预创建的token
			tokenKey := tokens[i%1000]
			err := tm.DelToken(tokenKey)
			if err != nil {
				b.Errorf("DelToken failed: %v", err)
			}
			// 如果删除成功，重新创建一个token
			if err == nil {
				newt, addErr := tm.AddToken(uint((i%1000)+1), 1, "192.168.1.1")
				if addErr == nil {
					tokens[i%1000] = newt
				}
			}
		}
	})
}

/**
 * BenchmarkBatchOperations 基准测试批量操作
 */
func BenchmarkBatchOperations(b *testing.B) {
	tm, err := wt.InitTM[any](models.ConfigRaw{
		MaxTokens:      1000,
		Delimiter:      "|",
		TokenRenewTime: "3600s",
		Language:       "zh",
	}, []models.GroupRaw{})
	if err != nil {
		b.Fatalf("Failed to initialize token manager: %v", err)
	}

	b.Run("BatchDeleteTokensByUserIDs", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			// 预先添加一些token
			userIDs := make([]uint, 10)
			for j := 0; j < 10; j++ {
				userID := uint(i*10 + j + 1)
				userIDs[j] = userID
				tm.AddToken(userID, 1, "192.168.1.1")
			}
			b.StartTimer()

			err := tm.BatchDeleteTokensByUserIDs(userIDs)
			if err != nil {
				b.Errorf("BatchDeleteTokensByUserIDs failed: %v", err)
			}
		}
	})

	b.Run("BatchDeleteTokensByGroupIDs", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			// 预先添加一些token
			groupIDs := make([]uint, 5)
			for j := 0; j < 5; j++ {
				groupID := uint(j + 1)
				groupIDs[j] = groupID
				for k := 0; k < 10; k++ {
					tm.AddToken(uint(i*50+j*10+k+1), groupID, "192.168.1.1")
				}
			}
			b.StartTimer()

			err := tm.BatchDeleteTokensByGroupIDs(groupIDs)
			if err != nil {
				b.Errorf("BatchDeleteTokensByGroupIDs failed: %v", err)
			}
		}
	})

	b.Run("GetTokensByUserID", func(b *testing.B) {
		// 预先添加一些token
		userIDs := make([]uint, 100)
		for i := 0; i < 100; i++ {
			userID := uint(i + 1)
			userIDs[i] = userID
			for j := 0; j < 5; j++ {
				tm.AddToken(userID, 1, "192.168.1.1")
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := userIDs[i%100]
			_ = tm.GetTokensByUserID(userID)
		}
	})

	b.Run("GetTokensByGroupID", func(b *testing.B) {
		// 预先添加一些token
		groupIDs := make([]uint, 10)
		for i := 0; i < 10; i++ {
			groupID := uint(i + 1)
			groupIDs[i] = groupID
			for j := 0; j < 50; j++ {
				tm.AddToken(uint(i*50+j+1), groupID, "192.168.1.1")
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			groupID := groupIDs[i%10]
			_ = tm.GetTokensByGroupID(groupID)
		}
	})
}

/**
 * BenchmarkConcurrentOperations 基准测试并发操作
 */
func BenchmarkConcurrentOperations(b *testing.B) {
	tm, err := wt.InitTM[any](models.ConfigRaw{
		MaxTokens:      1000,
		Delimiter:      "|",
		TokenRenewTime: "3600s",
		Language:       "zh",
	}, []models.GroupRaw{})
	if err != nil {
		b.Fatalf("Failed to initialize token manager: %v", err)
	}

	// 确保用户组存在
	err2 := tm.AddGroup(&models.GroupRaw{
		ID:                 1,
		Name:               "test_group",
		AllowedAPIs:        "",
		DeniedAPIs:         "",
		TokenExpire:        "3600",
		AllowMultipleLogin: 1,
	})
	if err2 != nil {
		b.Fatalf("Failed to add group: %v", err2)
	}

	b.Run("ConcurrentAddToken", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				userID := uint(i%1000 + 1)
				ip := fmt.Sprintf("192.168.%d.%d", (i%254)+1, (i%254)+1)
				_, err := tm.AddToken(userID, 1, ip)
				if err != nil {
					b.Errorf("ConcurrentAddToken failed: %v", err)
				}
				i++
			}
		})
	})

	b.Run("ConcurrentGetToken", func(b *testing.B) {
		// 预先添加一些token
		tokens := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			tokenKey, _ := tm.AddToken(uint(i+1), 1, "192.168.1.1")
			tokens[i] = tokenKey
		}

		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				tokenKey := tokens[i%1000]
				_, _ = tm.GetToken(tokenKey)
				i++
			}
		})
	})

	b.Run("ConcurrentMixedOperations", func(b *testing.B) {
		// 预先创建一些token供测试使用
		tokens := make([]string, 100)
		for i := 0; i < 100; i++ {
			tokenKey, err := tm.AddToken(uint(i+1), 1, "192.168.1.1")
			if err == nil {
				tokens[i] = tokenKey
			}
		}

		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				switch i % 4 {
				case 0: // AddToken
					userID := uint(i%1000 + 1)
					ip := fmt.Sprintf("192.168.%d.%d", (i%254)+1, (i%254)+1)
					tm.AddToken(userID, 1, ip)
				case 1: // GetToken
					if len(tokens) > 0 {
						tokenKey := tokens[i%len(tokens)]
						if tokenKey != "" {
							_, _ = tm.GetToken(tokenKey)
						}
					}
				case 2: // Auth
					if len(tokens) > 0 {
						tokenKey := tokens[i%len(tokens)]
						if tokenKey != "" {
							tm.Auth(tokenKey, "192.168.1.1", "/api/test")
						}
					}
				case 3: // DelToken (偶尔删除，但不要太频繁)
					if i%20 == 0 && len(tokens) > 0 {
						tokenKey := tokens[i%len(tokens)]
						if tokenKey != "" {
							tm.DelToken(tokenKey)
							// 重新创建一个token
							newt, err := tm.AddToken(uint((i%100)+1), 1, "192.168.1.1")
						if err == nil {
								tokens[i%len(tokens)] = newt
							}
						}
					}
				}
				i++
			}
		})
	})
}

/**
 * BenchmarkSecurityOperations 基准测试安全操作
 */
func BenchmarkSecurityOperations(b *testing.B) {
	sm := wt.NewSecurityManager("test_password")

	b.Run("EncryptToken", func(b *testing.B) {
		testToken := "test_token_for_encryption_benchmark_" + strconv.Itoa(rand.Int())
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := sm.EncryptToken(testToken)
			if err != nil {
				b.Errorf("EncryptToken failed: %v", err)
			}
		}
	})

	b.Run("DecryptToken", func(b *testing.B) {
		testToken := "test_token_for_decryption_benchmark_" + strconv.Itoa(rand.Int())
		encryptedToken, err := sm.EncryptToken(testToken)
		if err != nil {
			b.Fatalf("Failed to encrypt token for benchmark: %v", err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := sm.DecryptToken(encryptedToken)
			if err != nil {
				b.Errorf("DecryptToken failed: %v", err)
			}
		}
	})

	b.Run("HashData", func(b *testing.B) {
		testData := "test_data_for_hashing_benchmark_" + strconv.Itoa(rand.Int())
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = sm.HashSensitiveData(testData)
		}
	})

	b.Run("ValidateTokenFormat", func(b *testing.B) {
		testToken := "valid_token_format_test_" + strconv.Itoa(rand.Int())
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = wt.ValidateTokenFormat(testToken)
		}
	})

	b.Run("SanitizeInput", func(b *testing.B) {
		testInput := "<script>alert('test');</script>" + strconv.Itoa(rand.Int())
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = wt.SanitizeInput(testInput)
		}
	})
}

/**
 * BenchmarkConfigValidation 基准测试配置验证
 */
func BenchmarkConfigValidation(b *testing.B) {
	config := &models.ConfigRaw{
		MaxTokens:      1000,
		Delimiter:      "|",
		TokenRenewTime: "24h",
		Language:       "en",
	}

	b.Run("ValidateConfig", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err := wt.ValidateConfig(*config)
			if err != nil {
				b.Errorf("ValidateConfig failed: %v", err)
			}
		}
	})

	groupRaw := models.GroupRaw{
		ID:                 1,
		Name:               "test",
		TokenExpire:        "3600",
		AllowMultipleLogin: 1,
	}

	b.Run("ValidateGroupRaw", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err := wt.ValidateGroupRaw(groupRaw)
			if err != nil {
				b.Errorf("ValidateGroupRaw failed: %v", err)
			}
		}
	})
}

/**
 * BenchmarkCompleteWorkflow 基准测试完整工作流程
 */
func BenchmarkCompleteWorkflow(b *testing.B) {
	tm, err := wt.InitTM[string](models.ConfigRaw{
		MaxTokens:      1000,
		Delimiter:      "|",
		TokenRenewTime: "3600s",
		Language:       "zh",
	}, []models.GroupRaw{})
	if err != nil {
		b.Fatalf("Failed to initialize token manager: %v", err)
	}

	// 确保用户组存在
	err = tm.AddGroup(&models.GroupRaw{
		ID:                 1,
		Name:               "test_group",
		AllowedAPIs:        "/api/*",
		DeniedAPIs:         "",
		TokenExpire:        "1h",
		AllowMultipleLogin: 1,
	})
	if err != nil {
		b.Fatalf("Failed to add group: %v", err)
	}

	b.Run("CompleteTokenLifecycle", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := uint(i%1000 + 1)
			ip := fmt.Sprintf("192.168.%d.%d", (i%254)+1, (i%254)+1)

			// 添加token
			tokenKey, err := tm.AddToken(userID, 1, ip)
			if err != nil {
				b.Errorf("AddToken failed: %v", err)
				continue
			}

			// 获取token
			token, _ := tm.GetToken(tokenKey)
			if token == nil {
				b.Errorf("GetToken failed for key: %s", tokenKey)
				continue
			}

			// 检查token
			err = tm.Auth(tokenKey, ip, "/api/test")
			if err != nil {
				b.Errorf("Auth failed for key: %s", tokenKey)
			}

			// 设置用户数据
			err = tm.SetUserData(tokenKey, "benchmark_data_"+strconv.Itoa(i))
			if err != nil {
				b.Errorf("SetUserData failed: %v", err)
			}

			// 获取用户数据
			_, err = tm.GetUserData(tokenKey)
			if err != nil {
				b.Errorf("GetUserData failed: %v", err)
			}

			// 删除token
			err = tm.DelToken(tokenKey)
			if err != nil {
				b.Errorf("DelToken failed: %v", err)
			}
		}
	})
}
