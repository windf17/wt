# token

A lightweight and efficient token management system designed for API authentication and authorization. It provides flexible path-based permission control with support for allowed and denied API lists.

## Core Features

-   Token-based authentication with customizable expiration time
-   Hierarchical API path permission control
-   Fine-grained access control with allow/deny lists
-   High-performance memory cache with persistence support
-   Concurrent operation support
-   Debug logging and monitoring
-   Multilingual error messages (English/Chinese/Custom)
-   Token quantity limit control for resource management
-   Generic support for custom user data
-   Customizable error handling system
-   Multiple device login control
-   Real-time token statistics

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
    Language    string      // Language for error messages ("en", "zh", or custom)
    MaxTokens    int        // Maximum concurrent tokens, no limit if <= 0
    Debug       bool        // Enable debug mode
    Delimiter   string      // API path delimiter, default is space
}
```

## Error Handling System

### Built-in Error Codes

```go
const (
    ErrCodeSuccess              = 0    // Operation successful
    ErrCodeInvalidToken         = 1001 // Invalid token
    ErrCodeTokenNotFound        = 1002 // Token not found
    ErrCodeTokenExpired         = 1003 // Token expired
    ErrCodeInvalidUserID        = 1004 // Invalid user ID
    ErrCodeInvalidGroupID       = 1005 // Invalid group ID
    ErrCodeInvalidIP            = 1006 // Invalid IP address
    ErrCodeInvalidURL           = 1007 // Invalid URL
    ErrCodeAccessDenied         = 1008 // Access denied
    ErrCodeGroupNotFound        = 1009 // Group not found
    ErrCodeAddToken             = 1010 // Failed to add token
    ErrCodeCacheFileLoadFailed  = 1011 // Failed to load cache file
    ErrCodeCacheFileParseFailed = 1012 // Failed to parse cache file
)
```

### Custom Language Support

Example of adding French language support:

```go
// Register new language
fr := wtoken.RegisterLanguage("fr")

// Define custom error messages
frenchErrorMessages := map[wtoken.ILanguage]map[wtoken.ErrorCode]string{
    fr: {
        wtoken.ErrCodeSuccess:              "Opération réussie",
        wtoken.ErrCodeInvalidToken:         "Token invalide",
        wtoken.ErrCodeTokenNotFound:        "Token introuvable",
        // ... more error messages
    },
}

// Initialize with custom error messages
manager, err := wtoken.InitTM[any](&config, groups, frenchErrorMessages)
```

## Advanced Features

### Token Statistics

```go
// Get token statistics
stats := manager.GetStats()
fmt.Printf("Total tokens: %d\nActive tokens: %d\n", stats.TotalTokens, stats.ActiveTokens)
```

### Custom User Data

```go
// Save custom data
userInfo := "custom data"
if err = manager.SaveData(tokenKey, userInfo); err == nil {
    fmt.Println("Data saved successfully")
}

// Retrieve custom data
if data, err := manager.GetData(tokenKey); err == nil {
    fmt.Printf("User data: %v\n", data)
}
```

### Multiple Device Login Control

```go
groups := []wtoken.GroupRaw{
    {
        ID:                 1,
        AllowMultipleLogin: 0, // 0: single device only
    },
    {
        ID:                 2,
        AllowMultipleLogin: 1, // 1: allow multiple devices
    },
}
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
-   多语言错误信息（中文/英文/自定义）
-   Token 数量限制控制
-   泛型支持的自定义用户数据
-   可自定义的错误处理系统
-   多设备登录控制
-   实时 Token 统计

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

## 错误处理系统

### 内置错误码

```go
const (
    ErrCodeSuccess              = 0    // 操作成功
    ErrCodeInvalidToken         = 1001 // 无效的token
    ErrCodeTokenNotFound        = 1002 // token不存在
    ErrCodeTokenExpired         = 1003 // token已过期
    ErrCodeInvalidUserID        = 1004 // 无效的用户ID
    ErrCodeInvalidGroupID       = 1005 // 无效的用户组ID
    ErrCodeInvalidIP            = 1006 // 无效的IP地址
    ErrCodeInvalidURL           = 1007 // 无效的URL
    ErrCodeAccessDenied         = 1008 // 访问被拒绝
    ErrCodeGroupNotFound        = 1009 // 用户组不存在
    ErrCodeAddToken             = 1010 // 添加token失败
    ErrCodeCacheFileLoadFailed  = 1011 // 加载缓存文件失败
    ErrCodeCacheFileParseFailed = 1012 // 解析缓存文件失败
)
```

### 自定义语言支持

添加法语支持示例：

```go
// 注册新语言
fr := wtoken.RegisterLanguage("fr")

// 定义自定义错误信息
frenchErrorMessages := map[wtoken.ILanguage]map[wtoken.ErrorCode]string{
    fr: {
        wtoken.ErrCodeSuccess:              "Opération réussie",
        wtoken.ErrCodeInvalidToken:         "Token invalide",
        wtoken.ErrCodeTokenNotFound:        "Token introuvable",
        // ... 更多错误信息
    },
}

// 使用自定义错误信息初始化
manager, err := wtoken.InitTM[any](&config, groups, frenchErrorMessages)
```

## 高级功能

### Token 统计

```go
// 获取token统计信息
stats := manager.GetStats()
fmt.Printf("总token数：%d\n活跃token数：%d\n", stats.TotalTokens, stats.ActiveTokens)
```

### 自定义用户数据

```go
// 保存自定义数据
userInfo := "自定义数据"
if err = manager.SaveData(tokenKey, userInfo); err == nil {
    fmt.Println("数据保存成功")
}

// 获取自定义数据
if data, err := manager.GetData(tokenKey); err == nil {
    fmt.Printf("用户数据：%v\n", data)
}
```

### 多设备登录控制

```go
groups := []wtoken.GroupRaw{
    {
        ID:                 1,
        AllowMultipleLogin: 0, // 0: 仅允许单设备登录
    },
    {
        ID:                 2,
        AllowMultipleLogin: 1, // 1: 允许多设备登录
    },
}
```
