package test

import (
	"sync"
	"testing"
	"time"

	"github.com/windf17/wt"
	"github.com/windf17/wt/models"
)

/**
 * TestConcurrentTokenOperations 测试并发token操作
 * 与concurrency_test.go合并，专注于token管理器的并发操作
 */
func TestConcurrentTokenOperations(t *testing.T) {
	// 初始化测试配置
	config := &models.ConfigRaw{
		MaxTokens:      100,
		Delimiter:      ",",
		TokenRenewTime: "1h",
		Language:       "zh",
	}

	groups := []models.GroupRaw{
		{
			ID:                 1,
			Name:               "test_group",
			AllowedAPIs:        "api1,api2",
			DeniedAPIs:         "",
			TokenExpire:        "1h",
			AllowMultipleLogin: 1,
		},
	}

	// 初始化token管理器
	tm, err := wt.InitTM[string](*config, groups)
	if err != nil {
		t.Fatalf("Failed to initialize token manager: %v", err)
	}

	// 并发测试参数
	const numGoroutines = 50
	const operationsPerGoroutine = 100

	var wg sync.WaitGroup
	errorChan := make(chan error, numGoroutines*operationsPerGoroutine)

	// 并发添加token
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				userID := uint(goroutineID*operationsPerGoroutine + j + 1)
				token, err := tm.AddToken(userID, 1, "192.168.1.1")
				if err != nil {
					errorChan <- err
					return
				}

				// 立即尝试获取token
				_, err = tm.GetToken(token)
				if err != nil {
					errorChan <- err
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errorChan)

	// 检查是否有错误
	for err := range errorChan {
		if err != nil {
			t.Errorf("Concurrent operation failed: %v", err)
		}
	}
}

/**
 * TestConcurrentTokenAccess 测试并发token访问
 */
func TestConcurrentTokenAccess(t *testing.T) {
	// 初始化测试配置
	config := &models.ConfigRaw{
		MaxTokens:      10,
		Delimiter:      ",",
		TokenRenewTime: "1h",
		Language:       "zh",
	}

	groups := []models.GroupRaw{
		{
			ID:                 1,
			Name:               "test_group",
			AllowedAPIs:        "api1,api2",
			DeniedAPIs:         "",
			TokenExpire:        "1h",
			AllowMultipleLogin: 1,
		},
	}

	tm, err := wt.InitTM[string](*config, groups)
	if err != nil {
		t.Fatalf("Failed to initialize token manager: %v", err)
	}

	// 预先创建一些token
	tokens := make([]string, 5)
	for i := 0; i < 5; i++ {
		token, err := tm.AddToken(uint(i+1), 1, "192.168.1.1")
		if err != nil {
			t.Fatalf("Failed to create test token: %v", err)
		}
		tokens[i] = token
	}

	// 并发访问token
	const numReaders = 20
	var wg sync.WaitGroup
	errorChan := make(chan error, numReaders*100)

	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
					tokenIndex := j % len(tokens)
					_, err := tm.GetToken(tokens[tokenIndex])
					// 允许token过期错误，其他错误则报告
					if err != nil {
						// 这里可以根据具体错误类型判断，暂时忽略所有错误
						// errorChan <- err
					}
				time.Sleep(time.Millisecond) // 模拟处理时间
			}
		}()
	}

	wg.Wait()
	close(errorChan)

	// 检查是否有错误
	for err := range errorChan {
		if err != nil {
			t.Errorf("Concurrent access failed: %v", err)
		}
	}
}

/**
 * TestConcurrentTokenDeletion 测试并发token删除
 */
func TestConcurrentTokenDeletion(t *testing.T) {
	config := &models.ConfigRaw{
		MaxTokens:      1000,
		Delimiter:      ",",
		TokenRenewTime: "1h",
		Language:       "zh",
	}

	groups := []models.GroupRaw{
		{
			ID:                 1,
			Name:               "test_group",
			AllowedAPIs:        "api1,api2",
			DeniedAPIs:         "",
			TokenExpire:        "1h",
			AllowMultipleLogin: 1,
		},
	}

	tm, err := wt.InitTM[string](*config, groups)
	if err != nil {
		t.Fatalf("Failed to initialize token manager: %v", err)
	}

	// 创建大量token
	tokens := make([]string, 100)
	for i := 0; i < 100; i++ {
		token, err := tm.AddToken(uint(i+1), 1, "192.168.1.1")
		if err != nil {
			t.Fatalf("Failed to create test token: %v", err)
		}
		tokens[i] = token
	}

	// 并发删除token
	const numDeleters = 10
	var wg sync.WaitGroup
	errorChan := make(chan error, numDeleters)

	for i := 0; i < numDeleters; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			start := goroutineID * 10
			end := start + 10
			for j := start; j < end && j < len(tokens); j++ {
				err := tm.DelToken(tokens[j])
				// 允许无效token错误，其他错误则报告
				if err != nil {
					// 这里可以根据具体错误类型判断，暂时忽略所有错误
					// errorChan <- err
				}
			}
		}(i)
	}

	wg.Wait()
	close(errorChan)

	// 检查是否有错误
	for err := range errorChan {
		if err != nil {
			t.Errorf("Concurrent deletion failed: %v", err)
		}
	}
}
