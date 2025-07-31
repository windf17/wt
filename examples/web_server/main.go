/*
wt Web服务器完整示例

本示例演示了一个完整的用户权限管理系统，包括：

========== 用户组配置 ==========
1. 管理员组 (ID=1): 不允许重复登录，拥有管理员专属API权限
2. 普通用户组 (ID=2): 允许重复登录，拥有普通用户专属API权限

==========// 用户账号 ==========
- 2个管理员用户: admin1/admin123, admin2/admin123
- 100个普通用户: user1/user123 到 user100/user123

========== API端点 ==========
1. POST /api/login - 登录接口（所有用户可访问）
2. POST /api/logout - 登出接口（所有用户可访问）
3. GET /api/admin/dashboard - 管理员专属接口
4. GET /api/user/profile - 普通用户专属接口
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/windf17/wt"
	"github.com/windf17/wt/models"
)

// UserInfo 用户信息结构
type UserInfo struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// APIResponse 统一API响应格式
type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserInfo  UserInfo  `json:"user_info"`
}

var tokenManager *wt.Manager[UserInfo]

/**
 * initTokenManager 初始化Token管理器
 * 演示wt软件包的各项参数设置使用方法
 */
func initTokenManager() {
	// ========== 系统配置参数说明 ==========
	// wt.ConfigRaw 结构体包含了所有可配置的参数
	config := &wt.ConfigRaw{
		// 基础配置
		MaxTokens:      10000,          // 最大Token数量限制
		Delimiter:      ",",            // API路径分隔符
		TokenRenewTime: "24h",          // Token自动续期时间
		Language:       wt.LangChinese, // 错误消息语言

		// ========== Token验证参数配置 ==========
		// 这些参数可以在运行时动态配置，提供更灵活的验证策略
		MinTokenExpire: 60,    // Token最小过期时间（秒），默认60秒
		MaxTokenExpire: 86400, // Token最大过期时间（秒），默认24小时
	}

	// ========== 用户组配置 ==========
	// 配置普通用户组，允许访问所有API
	groups := []models.GroupRaw{
		{
			// 普通用户组 - 可以访问所有API
			ID:                 2,
			Name:               "普通用户",
			TokenExpire:        "4h",
			AllowMultipleLogin: 1,                       // 允许重复登录
			AllowedAPIs:        "/api/user,/api/logout", // 允许访问用户API和登出API
			DeniedAPIs:         "/api/admin",            // 禁止访问管理员接口
		},
	}

	// ========== Token管理器初始化 ==========
	// InitTM函数接受三个参数：配置、用户组、自定义验证函数
	manager := wt.InitTM[UserInfo](config, groups, nil)
	tokenManager = manager.(*wt.Manager[UserInfo])

	// ========== 其他可选配置示例 ==========
	// 以下是一些运行时可以调用的配置方法示例：

	// 1. 设置自定义Token生成器（可选）
	// tokenManager.SetTokenGenerator(customTokenGenerator)

	// 2. 设置自定义验证函数（可选）
	// tokenManager.SetCustomValidator(customValidator)

	// 3. 动态添加用户组（可选）
	// newGroup := wt.GroupRaw{
	//     ID: 4,
	//     Name: "临时用户",
	//     TokenExpire: "30m",
	//     AllowMultipleLogin: 0,
	//     AllowedAPIs: "/api/temp",
	//     DeniedAPIs: "/api/admin,/api/user",
	// }
	// tokenManager.AddGroup(newGroup)

	fmt.Println("Token管理器初始化完成，配置参数已生效")
	fmt.Printf("最大Token数量: %d\n", config.MaxTokens)
	fmt.Printf("Token续期时间: %s\n", config.TokenRenewTime)
	fmt.Printf("用户组数量: %d\n", len(groups))
}

/**
 * authMiddleware Token认证中间件
 * 演示wt软件包的Token验证和权限检查功能
 * @param {http.Handler} next 下一个处理器
 * @returns {http.Handler} 包装后的处理器
 */
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ========== 白名单路径配置 ==========
		// 跳过不需要Token验证的接口
		if r.URL.Path == "/api/login" {
			next.ServeHTTP(w, r)
			return
		}

		// ========== Token提取 ==========
		// 从HTTP Authorization头中提取Bearer Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			responseJSON(w, http.StatusUnauthorized, "缺少Authorization头", nil)
			return
		}

		// 解析Bearer Token格式："Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			responseJSON(w, http.StatusUnauthorized, "无效的Authorization格式", nil)
			return
		}

		token := parts[1]

		// ========== Token有效性验证 ==========
		// 获取客户端IP地址（与登录时保持一致的逻辑）
		clientIP := r.Header.Get("X-Real-IP")
		if clientIP == "" {
			clientIP = r.Header.Get("X-Forwarded-For")
		}
		if clientIP == "" {
			// 处理IPv6地址格式
			if strings.Contains(r.RemoteAddr, "[") {
				// IPv6格式: [::1]:port
				parts := strings.Split(r.RemoteAddr, "]")
				if len(parts) > 0 {
					clientIP = strings.TrimPrefix(parts[0], "[")
				}
			} else {
				// IPv4格式: 127.0.0.1:port
				clientIP = strings.Split(r.RemoteAddr, ":")[0]
			}
		}
		// 如果IP为空或无效，使用默认IP
		if clientIP == "" || clientIP == "::1" {
			clientIP = "127.0.0.1"
		}
		log.Printf("认证中间件 - 客户端IP: %s, 访问路径: %s", clientIP, r.URL.Path)
		log.Printf("Token验证开始 - Token: %s", token)

		// ========== Token验证和API权限验证 ==========
		// Auth函数包含完整的鉴权流程：Token验证、IP验证和API权限验证
		// 参数说明：
		// 1. token: 要验证的Token字符串
		// 2. clientIP: 客户端IP地址
		// 3. apiPath: 要访问的API路径（如：/api/user/profile）
		authResult := tokenManager.Auth(token, clientIP, r.URL.Path)
		log.Printf("Token验证结果: %v", authResult)
		if authResult != wt.E_Success {
			if authResult == wt.E_Unauthorized {
				responseJSON(w, http.StatusForbidden, fmt.Sprintf("权限不足: %v", authResult), nil)
			} else {
				responseJSON(w, http.StatusUnauthorized, fmt.Sprintf("认证失败: %v", authResult), nil)
			}
			return
		}

		// ========== 用户信息获取 ==========
		// GetUserData函数获取Token绑定的用户信息
		// 返回之前通过SetUserData设置的自定义用户数据
		userInfo, err := tokenManager.GetUserData(token)
		if err == wt.E_Success {
			// 将用户信息添加到请求头，供后续处理器使用
			r.Header.Set("X-User-ID", fmt.Sprintf("%d", userInfo.UserID))
			r.Header.Set("X-Username", userInfo.Username)
			r.Header.Set("X-User-Role", userInfo.Role)
			log.Printf("Token验证成功 - 用户: %s, 角色: %s, 访问路径: %s", userInfo.Username, userInfo.Role, r.URL.Path)
		}

		next.ServeHTTP(w, r)
	})
}

/**
 * responseJSON 统一JSON响应
 * @param {http.ResponseWriter} w 响应写入器
 * @param {int} statusCode HTTP状态码
 * @param {string} message 响应消息
 * @param {any} data 响应数据
 */
func responseJSON(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
	}

	log.Printf("发送响应: 状态码=%d, 消息=%s", statusCode, message)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("JSON编码错误: %v", err)
	}
	log.Printf("响应发送完成")
}

/**
 * loginHandler 用户登录处理器
 * @param {http.ResponseWriter} w 响应写入器
 * @param {*http.Request} r 请求对象
 */
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responseJSON(w, http.StatusMethodNotAllowed, "只支持POST方法", nil)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseJSON(w, http.StatusBadRequest, "无效的JSON格式", nil)
		return
	}

	// 用户验证 - 2个管理员用户 + 28个普通用户
	var userInfo UserInfo
	var groupID uint

	// 验证密码
	if req.Password != "admin123" && req.Password != "user123" {
		responseJSON(w, http.StatusUnauthorized, "用户名或密码错误", nil)
		return
	}

	// 管理员用户验证
	if req.Username == "admin1" || req.Username == "admin2" {
		if req.Password != "admin123" {
			responseJSON(w, http.StatusUnauthorized, "用户名或密码错误", nil)
			return
		}
		userID := uint(1001)
		if req.Username == "admin2" {
			userID = 1002
		}
		userInfo = UserInfo{
			UserID:   userID,
			Username: req.Username,
			Email:    req.Username + "@company.com",
			Role:     "administrator",
		}
		groupID = 2 // 所有用户都使用同一个用户组
	} else {
		// 普通用户验证 (user1 到 user100)
		var userID uint
		for i := 1; i <= 100; i++ {
			if req.Username == fmt.Sprintf("user%d", i) {
				if req.Password != "user123" {
					responseJSON(w, http.StatusUnauthorized, "用户名或密码错误", nil)
					return
				}
				userID = uint(2000 + i)
				break
			}
		}
		if userID == 0 {
			responseJSON(w, http.StatusUnauthorized, "用户名或密码错误", nil)
			return
		}
		userInfo = UserInfo{
			UserID:   userID,
			Username: req.Username,
			Email:    req.Username + "@company.com",
			Role:     "user",
		}
		groupID = 2 // 所有用户都使用同一个用户组
	}

	// 获取客户端IP
	clientIP := r.Header.Get("X-Real-IP")
	if clientIP == "" {
		clientIP = r.Header.Get("X-Forwarded-For")
	}
	if clientIP == "" {
		// 处理IPv6地址格式
		if strings.Contains(r.RemoteAddr, "[") {
			// IPv6格式: [::1]:port
			parts := strings.Split(r.RemoteAddr, "]")
			if len(parts) > 0 {
				clientIP = strings.TrimPrefix(parts[0], "[")
			}
		} else {
			// IPv4格式: 127.0.0.1:port
			clientIP = strings.Split(r.RemoteAddr, ":")[0]
		}
	}
	// 如果IP为空或无效，使用默认IP
	if clientIP == "" || clientIP == "::1" {
		clientIP = "127.0.0.1"
	}
	log.Printf("客户端IP: %s", clientIP)

	// ========== Token创建参数说明 ==========
	// AddToken函数参数详解：
	// 1. userID: 用户唯一标识符
	// 2. groupID: 用户组ID（决定权限和Token过期时间）
	// 3. clientIP: 客户端IP地址（用于安全验证）
	token, err := tokenManager.AddToken(userInfo.UserID, groupID, clientIP)
	if err != wt.E_Success {
		responseJSON(w, http.StatusInternalServerError, fmt.Sprintf("Token创建失败: %v", err), nil)
		return
	}

	// Token创建成功后的日志记录
	log.Printf("Token创建成功 - 用户ID: %d, 组ID: %d, IP: %s", userInfo.UserID, groupID, clientIP)

	// ========== 用户数据设置说明 ==========
	// SetUserData函数用于将自定义用户信息绑定到Token
	// 这样在后续的Token验证中可以直接获取用户信息，无需再次查询数据库
	log.Printf("开始设置用户数据")
	go func() {
		// 在goroutine中异步设置用户数据，避免阻塞响应
		// 参数说明：
		// 1. token: 要绑定数据的Token字符串
		// 2. userInfo: 要绑定的用户信息（泛型类型，本例中为UserInfo结构体）
		if err := tokenManager.SetUserData(token, userInfo); err != wt.E_Success {
			log.Printf("设置用户数据失败: %v", err)
		} else {
			log.Printf("用户数据设置成功 - 用户: %s, 邮箱: %s", userInfo.Username, userInfo.Email)
		}
	}()

	// 获取Token信息以确定过期时间
	log.Printf("开始获取Token信息")
	tokenInfo, _ := tokenManager.GetToken(token)
	log.Printf("Token信息获取完成")
	var expiresAt time.Time
	if tokenInfo != nil {
		expiresAt = tokenInfo.LoginTime.Add(time.Duration(tokenInfo.ExpireSeconds) * time.Second)
		log.Printf("过期时间计算完成: %v", expiresAt)
	}

	response := LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		UserInfo:  userInfo,
	}

	responseJSON(w, http.StatusOK, "登录成功", response)
}

/**
 * logoutHandler 用户登出处理器
 * 演示wt软件包的Token删除功能
 * @param {http.ResponseWriter} w 响应写入器
 * @param {*http.Request} r 请求对象
 */
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responseJSON(w, http.StatusMethodNotAllowed, "只支持POST方法", nil)
		return
	}

	// ========== Token提取 ==========
	// 从Authorization头获取要删除的Token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		responseJSON(w, http.StatusBadRequest, "缺少Authorization头", nil)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		responseJSON(w, http.StatusBadRequest, "无效的Authorization格式", nil)
		return
	}

	token := parts[1]

	// ========== Token删除说明 ==========
	// DelToken函数用于删除指定的Token，实现用户登出功能
	// 参数说明：
	// 1. token: 要删除的Token字符串
	// 删除后该Token将立即失效，无法再用于API访问
	// 同时会清理相关的用户数据和缓存信息
	if err := tokenManager.DelToken(token); err != wt.E_Success {
		responseJSON(w, http.StatusInternalServerError, fmt.Sprintf("登出失败: %v", err), nil)
		return
	}

	log.Printf("用户登出成功，Token已删除: %s", token[:10]+"...")
	responseJSON(w, http.StatusOK, "登出成功", nil)
}

/**
 * profileHandler 获取用户资料处理器
 * @param {http.ResponseWriter} w 响应写入器
 * @param {*http.Request} r 请求对象
 */
func profileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responseJSON(w, http.StatusMethodNotAllowed, "只支持GET方法", nil)
		return
	}

	userID := r.Header.Get("X-User-ID")
	username := r.Header.Get("X-Username")

	responseJSON(w, http.StatusOK, "获取用户资料成功", map[string]string{
		"user_id":  userID,
		"username": username,
		"message":  "这是用户资料页面",
	})
}

/**
 * adminHandler 管理员接口处理器
 * @param {http.ResponseWriter} w 响应写入器
 * @param {*http.Request} r 请求对象
 */
func adminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responseJSON(w, http.StatusMethodNotAllowed, "只支持GET方法", nil)
		return
	}

	// 获取Token统计信息
	stats := tokenManager.GetStats()

	responseJSON(w, http.StatusOK, "管理员数据获取成功", map[string]any{
		"message":       "这是管理员页面",
		"token_stats":   stats,
		"current_time":  time.Now().Format("2006-01-02 15:04:05"),
		"server_status": "running",
	})
}

/**
 * WebServerExample Web服务器示例主函数
 */
func WebServerExample() {
	fmt.Println("=== wt Web服务器完整示例 ===")

	// 初始化Token管理器
	initTokenManager()

	// 创建HTTP路由
	mux := http.NewServeMux()

	// 注册4个API端点
	mux.HandleFunc("/api/login", loginHandler)           // 登录接口（所有用户可访问）
	mux.HandleFunc("/api/logout", logoutHandler)         // 登出接口（所有用户可访问）
	mux.HandleFunc("/api/user/profile", profileHandler)  // 普通用户专属接口
	mux.HandleFunc("/api/admin/dashboard", adminHandler) // 管理员专属接口

	// 应用认证中间件
	handler := authMiddleware(mux)

	// 启动服务器
	port := ":8081"
	fmt.Printf("服务器启动在端口 %s\n", port)
	fmt.Println("\n========== API端点 ==========")
	fmt.Println("POST /api/login - 登录接口（所有用户可访问）")
	fmt.Println("POST /api/logout - 登出接口（所有用户可访问）")
	fmt.Println("GET  /api/user/profile - 普通用户专属接口")
	fmt.Println("GET  /api/admin/dashboard - 管理员专属接口")
	fmt.Println("\n========== 测试用户 ==========")
	fmt.Println("管理员用户（不可重复登录）:")
	fmt.Println("  admin1/admin123")
	fmt.Println("  admin2/admin123")
	fmt.Println("普通用户（可重复登录）:")
	fmt.Println("  user1/user123 到 user100/user123")
	fmt.Println("\n========== 使用示例 ==========")
	fmt.Println("# 管理员登录")
	fmt.Println("curl -X POST http://localhost:8081/api/login -H 'Content-Type: application/json' -d '{\"username\":\"admin1\",\"password\":\"admin123\"}'")
	fmt.Println("# 普通用户登录")
	fmt.Println("curl -X POST http://localhost:8081/api/login -H 'Content-Type: application/json' -d '{\"username\":\"user1\",\"password\":\"user123\"}'")

	log.Fatal(http.ListenAndServe(port, handler))
}

func main() {
	// 运行Web服务器
	WebServerExample()
}
