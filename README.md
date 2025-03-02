# token

A lightweight and efficient token management system designed for API authentication and authorization. It provides flexible path-based permission control with support for allowed and denied API lists.

## Core Features

-   Token-based authentication with customizable expiration time
-   Hierarchical API path permission control
-   Fine-grained access control with allow/deny lists
-   High-performance memory cache with persistence support
-   Concurrent operation support
-   Debug logging and monitoring
-   Bilingual error messages (English/Chinese)
-   Token quantity limit control for resource management
-   Generic support for custom user data

## Installation

```bash
go get github.com/windf17/wtoken
```

## Quick Start

```go
package main

import (
    "github.com/windf17/wtoken"
    "fmt"
)

func main() {
    // Initialize token manager with configuration
    config := wtoken.Config{
        CacheFilePath: "token.cache",
        Language:      "en",
        MaxTokens:     1000,
        Debug:         true,
        Delimiter:     " ",
    }

    // Configure user groups
    groups := []wtoken.GroupRaw{
        {
            ID:                 1,
            AllowedAPIs:        "/api/user /api/product",
            DeniedAPIs:         "/api/admin",
            TokenExpire:        "3600",
            AllowMultipleLogin: 0,
        },
    }

    // Initialize token manager with generic type support
    manager, err := wtoken.InitTM[any](&config, groups, nil)
    if err != nil {
        fmt.Printf("Failed to initialize token manager: %v\n", err)
        return
    }

    // Generate token
    tokenKey, errData := manager.AddToken(1, 1, "127.0.0.1")
    if errData.Code != wtoken.ErrCodeSuccess {
        fmt.Printf("Failed to generate token: %v\n", errData.Error())
        return
    }

    // Authenticate API access
    errData = manager.Authenticate(tokenKey, "/api/user", "127.0.0.1")
    if errData.Code == wtoken.ErrCodeSuccess {
        fmt.Println("Authentication successful")
    }
}
```

## Configuration

```go
type Config struct {
    CacheFilePath string    // Cache file path
    Language    string      // "en" for English or "zh" for Chinese error messages
    MaxTokens    int        // Maximum concurrent tokens, no limit if <= 0
    Debug       bool        // Enable debug mode
    Delimiter   string      // API path delimiter, default is space
}
```

## API Reference

### Token Structure

```go
type Token[T any] struct {
    UserID         uint        // User ID
    GroupID        uint        // Group ID
    LoginTime      time.Time   // Initial login time
    ExpireTime     time.Time   // Expiration time (zero value means never expire)
    LastAccessTime time.Time   // Last access time
    UserData       T          // Custom user data with generic type support
    IP             string      // Token user's IP address
}
```

### Group Configuration

```go
type GroupRaw struct {
    ID                 uint   `json:"id"`              // Group ID
    AllowedAPIs        string `json:"allowedApis"`     // Space-separated list of allowed APIs
    DeniedAPIs         string `json:"deniedApis"`      // Space-separated list of denied APIs
    TokenExpire        string `json:"tokenExpire"`     // Token expiration time in seconds, 0 means never expire
    AllowMultipleLogin int    `json:"allowMultipleLogin"` // 1: allow multiple device login, others: single device only
}
```

### Permission System

1. **Priority Order**

    - Whitelist check takes precedence over blacklist
    - When present in both lists, blacklist takes precedence

2. **Path Matching Rules**

    - Uses standardized URL paths
    - Follows longest prefix matching principle

3. **IP Check Policy**
    - AllowMultipleLogin=false: requires consistent login IP
    - AllowMultipleLogin=true: allows different IP logins

### Debug Mode Output

When debug mode is enabled, the system logs detailed permission check process:

```
Verifying /api/product:
  Checking AllowedAPIs:
    - Match found: /api/product (length: 13)
  Checking DeniedAPIs:
    - No match found
Result: Access allowed (matching allow rule: /api/product)
```

---

# token

一个轻量高效的 Token 管理系统，专为 API 认证和授权设计。它提供了基于路径的灵活权限控制，支持允许和拒绝 API 列表。

## 核心功能

-   基于 Token 的身份验证，支持自定义过期时间
-   层级式的 API 路径权限控制
-   细粒度的访问控制，支持允许/拒绝列表
-   高性能内存缓存，支持持久化
-   并发操作支持
-   调试日志和监控
-   双语错误信息（中文/英文）
-   Token 数量限制控制
-   泛型支持的自定义用户数据

## 安装

```bash
go get github.com/windf17/wtoken
```

## 快速开始

```go
package main

import (
    "github.com/windf17/wtoken"
    "fmt"
)

func main() {
    // 初始化 token 管理器配置
    config := wtoken.Config{
        CacheFilePath: "token.cache", // token缓存文件路径
        Language:      "zh",          // 错误信息语言
        MaxTokens:     1000,          // 最大token数量
        Debug:         true,          // 是否开启调试模式
        Delimiter:     " ",           // API分隔符
    }

    // 配置用户组
    groups := []wtoken.GroupRaw{
        {
            ID:                 1,
            AllowedAPIs:        "/api/user /api/product",
            DeniedAPIs:         "/api/admin",
            TokenExpire:        "3600",
            AllowMultipleLogin: 0,
        },
    }

    // 使用泛型支持初始化 token 管理器
    manager, err := wtoken.InitTM[any](&config, groups, nil)
    if err != nil {
        fmt.Printf("初始化token管理器失败：%v\n", err)
        return
    }

    // 生成用户token
    tokenKey, errData := manager.AddToken(1, 1, "127.0.0.1")
    if errData.Code != wtoken.ErrCodeSuccess {
        fmt.Printf("生成token失败：%v\n", errData.Error())
        return
    }

    // API鉴权测试
    errData = manager.Authenticate(tokenKey, "/api/user", "127.0.0.1")
    if errData.Code == wtoken.ErrCodeSuccess {
        fmt.Println("鉴权成功")
    }
}
```
