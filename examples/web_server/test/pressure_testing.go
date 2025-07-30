package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

/**
 * API响应结构体
 */
type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

/**
 * 登录响应数据结构体
 */
type LoginResponseData struct {
	Token     string   `json:"token"`
	ExpiresAt string   `json:"expires_at"`
	UserInfo  UserInfo `json:"user_info"`
}

/**
 * 用户信息结构体
 */
type UserInfo struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

/**
 * 用户会话结构体
 */
type UserSession struct {
	UserID   int
	Token    string
	Username string
}

/**
 * API测试结果结构体
 */
type APITestResult struct {
	Endpoint     string
	Success      bool
	ResponseTime time.Duration
	StatusCode   int
	Error        string
}

/**
 * 压力测试结果结构体
 */
type PressureTestResult struct {
	UserID        int
	Username      string
	LoginSuccess  bool
	LoginTime     time.Duration
	APIResults    []APITestResult
	LogoutSuccess bool
	LogoutTime    time.Duration
	TotalTime     time.Duration
}

/**
 * 测试统计结构体
 */
type TestStats struct {
	TotalUsers        int
	SuccessfulLogins  int
	FailedLogins      int
	SuccessfulLogouts int
	FailedLogouts     int
	// 可访问的API统计
	AllowedAPIRequests    int
	SuccessfulAllowedAPIs int
	FailedAllowedAPIs     int
	AverageAllowedAPITime time.Duration
	// 不可访问的API统计
	DeniedAPIRequests    int
	SuccessfulDeniedAPIs int
	FailedDeniedAPIs     int
	AverageDeniedAPITime time.Duration
	// 总体统计
	TotalAPIRequests  int
	SuccessfulAPIs    int
	FailedAPIs        int
	AverageLoginTime  time.Duration
	AverageLogoutTime time.Duration
	AverageAPITime    time.Duration
	TotalTestTime     time.Duration
}

const (
	baseURL      = "http://localhost:8081"
	testUsers    = 1000 // 测试用户数量
	testDuration = 120  // 测试持续时间（秒）
)

// API端点列表
var apiEndpoints = []string{
	"/api/user/profile",
	"/api/admin/dashboard",
}

/**
 * 执行用户登录
 * @param {int} userID 用户ID
 * @returns {UserSession, error} 用户会话和错误信息
 */
func performLogin(userID int) (*UserSession, time.Duration, error) {
	start := time.Now()

	// 循环使用1-100的用户账号，因为服务器只配置了100个用户
	actualUserID := ((userID - 1) % 100) + 1
	username := fmt.Sprintf("user%d", actualUserID)
	password := "user123"

	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, _ := json.Marshal(loginData)
	resp, err := http.Post(baseURL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, time.Since(start), fmt.Errorf("登录请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, time.Since(start), fmt.Errorf("读取响应失败: %v", err)
	}

	var apiResp APIResponse
	if err = json.Unmarshal(body, &apiResp); err != nil {
		return nil, time.Since(start), fmt.Errorf("解析响应失败: %v", err)
	}

	if apiResp.Code != 200 {
		return nil, time.Since(start), fmt.Errorf("登录失败: %s", apiResp.Message)
	}

	// 解析登录响应数据
	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, time.Since(start), fmt.Errorf("解析登录数据失败: %v", err)
	}

	var loginResp LoginResponseData
	if err := json.Unmarshal(dataBytes, &loginResp); err != nil {
		return nil, time.Since(start), fmt.Errorf("解析登录数据失败: %v", err)
	}

	if loginResp.Token == "" {
		return nil, time.Since(start), fmt.Errorf("登录失败: Token为空")
	}

	session := &UserSession{
		UserID:   userID,
		Token:    loginResp.Token,
		Username: username,
	}

	return session, time.Since(start), nil
}

/**
 * 执行API请求
 * @param {*UserSession} session 用户会话
 * @param {string} endpoint API端点
 * @returns {APITestResult} API测试结果
 */
func performAPIRequest(session *UserSession, endpoint string) APITestResult {
	start := time.Now()

	req, err := http.NewRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		return APITestResult{
			Endpoint:     endpoint,
			Success:      false,
			ResponseTime: time.Since(start),
			StatusCode:   0,
			Error:        fmt.Sprintf("创建请求失败: %v", err),
		}
	}

	req.Header.Set("Authorization", "Bearer "+session.Token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return APITestResult{
			Endpoint:     endpoint,
			Success:      false,
			ResponseTime: time.Since(start),
			StatusCode:   0,
			Error:        fmt.Sprintf("请求失败: %v", err),
		}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResp APIResponse
	json.Unmarshal(body, &apiResp)

	success := resp.StatusCode == 200 && apiResp.Code == 200
	errorMsg := ""
	if !success {
		errorMsg = fmt.Sprintf("状态码: %d, 消息: %s", resp.StatusCode, apiResp.Message)
	}

	return APITestResult{
		Endpoint:     endpoint,
		Success:      success,
		ResponseTime: time.Since(start),
		StatusCode:   resp.StatusCode,
		Error:        errorMsg,
	}
}

/**
 * 执行用户登出
 * @param {*UserSession} session 用户会话
 * @returns {time.Duration, error} 登出时间和错误信息
 */
func performLogout(session *UserSession) (time.Duration, error) {
	start := time.Now()

	req, err := http.NewRequest("POST", baseURL+"/api/logout", nil)
	if err != nil {
		return time.Since(start), fmt.Errorf("创建登出请求失败: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+session.Token)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return time.Since(start), fmt.Errorf("登出请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return time.Since(start), fmt.Errorf("登出失败，状态码: %d", resp.StatusCode)
	}

	return time.Since(start), nil
}

/**
 * 模拟单个用户的完整测试流程
 * @param {int} userID 用户ID
 * @param {chan PressureTestResult} results 结果通道
 * @param {*sync.WaitGroup} wg 等待组
 */
func simulateUser(userID int, results chan<- PressureTestResult, wg *sync.WaitGroup) {
	defer wg.Done()

	testStart := time.Now()
	result := PressureTestResult{
		UserID:   userID,
		Username: fmt.Sprintf("user%d", userID),
	}

	// 1. 登录
	session, loginTime, err := performLogin(userID)
	result.LoginTime = loginTime
	if err != nil {
		result.LoginSuccess = false
		result.TotalTime = time.Since(testStart)
		results <- result
		return
	}
	result.LoginSuccess = true

	// 2. 随机访问API
	numRequests := rand.Intn(10) + 5 // 5-14个请求
	for i := 0; i < numRequests; i++ {
		endpoint := apiEndpoints[rand.Intn(len(apiEndpoints))]
		apiResult := performAPIRequest(session, endpoint)
		result.APIResults = append(result.APIResults, apiResult)

		// 随机等待时间，模拟真实用户行为
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}

	// 3. 登出
	logoutTime, err := performLogout(session)
	result.LogoutTime = logoutTime
	if err != nil {
		result.LogoutSuccess = false
	} else {
		result.LogoutSuccess = true
	}

	result.TotalTime = time.Since(testStart)
	results <- result
}

/**
 * 计算测试统计信息
 * @param {[]PressureTestResult} results 测试结果列表
 * @returns {TestStats} 统计信息
 */
func calculateStats(results []PressureTestResult) TestStats {
	stats := TestStats{
		TotalUsers: len(results),
	}

	var totalLoginTime, totalLogoutTime, totalAPITime time.Duration
	var totalAllowedAPITime, totalDeniedAPITime time.Duration
	apiCount := 0
	allowedAPICount := 0
	deniedAPICount := 0

	for _, result := range results {
		if result.LoginSuccess {
			stats.SuccessfulLogins++
			totalLoginTime += result.LoginTime
		} else {
			stats.FailedLogins++
		}

		// 统计登出结果
		if result.LogoutSuccess {
			stats.SuccessfulLogouts++
			totalLogoutTime += result.LogoutTime
		} else {
			stats.FailedLogouts++
		}

		for _, apiResult := range result.APIResults {
			stats.TotalAPIRequests++
			apiCount++
			totalAPITime += apiResult.ResponseTime

			// 判断是否为可访问的API（普通用户可以访问/api/user/profile，不能访问/api/admin/dashboard）
			isAllowedAPI := apiResult.Endpoint == "/api/user/profile"

			if isAllowedAPI {
				stats.AllowedAPIRequests++
				allowedAPICount++
				totalAllowedAPITime += apiResult.ResponseTime
				if apiResult.Success {
					stats.SuccessfulAllowedAPIs++
				} else {
					stats.FailedAllowedAPIs++
				}
			} else {
				stats.DeniedAPIRequests++
				deniedAPICount++
				totalDeniedAPITime += apiResult.ResponseTime
				if apiResult.Success {
					stats.SuccessfulDeniedAPIs++
				} else {
					stats.FailedDeniedAPIs++
				}
			}

			if apiResult.Success {
				stats.SuccessfulAPIs++
			} else {
				stats.FailedAPIs++
			}
		}
	}

	if stats.SuccessfulLogins > 0 {
		stats.AverageLoginTime = totalLoginTime / time.Duration(stats.SuccessfulLogins)
	}

	if stats.SuccessfulLogouts > 0 {
		stats.AverageLogoutTime = totalLogoutTime / time.Duration(stats.SuccessfulLogouts)
	}

	if allowedAPICount > 0 {
		stats.AverageAllowedAPITime = totalAllowedAPITime / time.Duration(allowedAPICount)
	}

	if deniedAPICount > 0 {
		stats.AverageDeniedAPITime = totalDeniedAPITime / time.Duration(deniedAPICount)
	}

	if apiCount > 0 {
		stats.AverageAPITime = totalAPITime / time.Duration(apiCount)
	}

	return stats
}

/**
 * 打印测试统计报告
 * @param {TestStats} stats 统计信息
 * @param {time.Duration} totalTime 总测试时间
 */
func printTestReport(stats TestStats, totalTime time.Duration) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("                           压力测试报告")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Printf("测试开始时间: %s\n", time.Now().Add(-totalTime).Format("2006-01-02 15:04:05"))
	fmt.Printf("测试结束时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("总测试时间: %v\n", totalTime)

	fmt.Println("\n📊 用户登录统计:")
	fmt.Printf("   总用户数: %d\n", stats.TotalUsers)
	fmt.Printf("   成功登录: %d\n", stats.SuccessfulLogins)
	fmt.Printf("   失败登录: %d\n", stats.FailedLogins)
	fmt.Printf("   登录成功率: %.2f%%\n", float64(stats.SuccessfulLogins)/float64(stats.TotalUsers)*100)
	fmt.Printf("   平均登录时间: %v\n", stats.AverageLoginTime)

	fmt.Println("\n🚪 用户登出统计:")
	fmt.Printf("   成功登出: %d\n", stats.SuccessfulLogouts)
	fmt.Printf("   失败登出: %d\n", stats.FailedLogouts)
	fmt.Printf("   登出成功率: %.2f%%\n", float64(stats.SuccessfulLogouts)/float64(stats.TotalUsers)*100)
	fmt.Printf("   平均登出时间: %v\n", stats.AverageLogoutTime)

	fmt.Println("\n✅ 可访问的API请求:")
	fmt.Printf("   总请求数: %d\n", stats.AllowedAPIRequests)
	fmt.Printf("   成功请求: %d\n", stats.SuccessfulAllowedAPIs)
	fmt.Printf("   失败请求: %d\n", stats.FailedAllowedAPIs)
	if stats.AllowedAPIRequests > 0 {
		fmt.Printf("   成功率: %.2f%%\n", float64(stats.SuccessfulAllowedAPIs)/float64(stats.AllowedAPIRequests)*100)
	} else {
		fmt.Printf("   成功率: 0.00%%\n")
	}
	fmt.Printf("   平均响应时间: %v\n", stats.AverageAllowedAPITime)

	fmt.Println("\n❌ 不可访问的API请求:")
	fmt.Printf("   总请求数: %d\n", stats.DeniedAPIRequests)
	fmt.Printf("   成功请求: %d\n", stats.SuccessfulDeniedAPIs)
	fmt.Printf("   失败请求: %d\n", stats.FailedDeniedAPIs)
	if stats.DeniedAPIRequests > 0 {
		fmt.Printf("   成功率: %.2f%%\n", float64(stats.SuccessfulDeniedAPIs)/float64(stats.DeniedAPIRequests)*100)
	} else {
		fmt.Printf("   成功率: 0.00%%\n")
	}
	fmt.Printf("   平均响应时间: %v\n", stats.AverageDeniedAPITime)

	fmt.Println("\n⚡ 性能指标:")
	fmt.Printf("   并发用户数: %d\n", testUsers)
	fmt.Printf("   平均每用户API请求: %.1f\n", float64(stats.TotalAPIRequests)/float64(stats.TotalUsers))
	fmt.Printf("   总吞吐量: %.2f 请求/秒\n", float64(stats.TotalAPIRequests)/totalTime.Seconds())

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("测试完成时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(strings.Repeat("=", 80))
}

/**
 * 主函数 - 执行压力测试
 */
func main() {
	fmt.Println("🚀 启动Web服务器压力测试...")
	fmt.Printf("测试配置: %d个并发用户，持续%d秒\n", testUsers, testDuration)
	fmt.Println("目标服务器:", baseURL)
	fmt.Println()

	// 检查服务器是否可用
	resp, err := http.Get(baseURL + "/api/login")
	if err != nil {
		fmt.Printf("❌ 无法连接到服务器 %s: %v\n", baseURL, err)
		fmt.Println("请确保Web服务器正在运行")
		return
	}
	resp.Body.Close()

	if resp.StatusCode != 405 && resp.StatusCode != 200 {
		fmt.Printf("⚠️  服务器响应异常，状态码: %d\n", resp.StatusCode)
	}

	fmt.Println("✅ 服务器连接正常，开始压力测试...")
	fmt.Println()

	testStart := time.Now()

	// 创建结果通道和等待组
	results := make(chan PressureTestResult, testUsers)
	var wg sync.WaitGroup

	// 启动所有用户协程
	for i := 1; i <= testUsers; i++ {
		wg.Add(1)
		go simulateUser(i, results, &wg)

		// 避免同时发起太多连接
		if i%10 == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	// 等待所有用户完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	var allResults []PressureTestResult
	for result := range results {
		allResults = append(allResults, result)
		fmt.Printf("\r用户 %s 测试完成 [%d/%d]", result.Username, len(allResults), testUsers)
	}

	totalTime := time.Since(testStart)
	fmt.Println("\n\n✅ 所有用户测试完成!")

	// 计算并显示统计信息
	stats := calculateStats(allResults)
	stats.TotalTestTime = totalTime
	printTestReport(stats, totalTime)

	// 保存详细结果到文件（可选）
	if len(allResults) > 0 {
		fmt.Println("\n💾 保存详细测试结果到 pressure_test_results.json")
		if data, err := json.MarshalIndent(allResults, "", "  "); err == nil {
			if err := os.WriteFile("pressure_test_results.json", data, 0644); err == nil {
				fmt.Println("✅ 结果已保存")
			} else {
				fmt.Printf("❌ 保存失败: %v\n", err)
			}
		}
	}
}
