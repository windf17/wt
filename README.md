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
    if errData != wtoken.ErrSuccess {
        fmt.Printf("Failed to generate token: %v\n", errData.Error())
        return
    }

    // Authenticate API access
    errData = manager.Authenticate(tokenKey, "/api/user", "127.0.0.1")
    if errData == wtoken.ErrSuccess {
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

## Token Structure

```go
type Token[T any] struct {
    UserID         uint      // User ID
    GroupID        uint      // Group ID
    LoginTime      time.Time // Login time
    ExpireSeconds  int64     // Token expiration seconds, 0 means never expire
    LastAccessTime time.Time // Last access time
    UserData       T        // Custom user data
    IP             string   // User's IP address
}
```

## Error Handling System

### Built-in Error Codes

```go
const (
    ErrSuccess              = 0    // Operation successful
    ErrInvalidToken         = 1001 // Invalid token
    ErrTokenNotFound        = 1002 // Token not found
    ErrTokenExpired         = 1003 // Token expired
    ErrInvalidUserID        = 1004 // Invalid user ID
    ErrInvalidGroupID       = 1005 // Invalid group ID
    ErrInvalidIP            = 1006 // Invalid IP address
    ErrInvalidURL           = 1007 // Invalid URL
    ErrAccessDenied         = 1008 // Access denied
    ErrGroupNotFound        = 1009 // Group not found
    ErrAddToken             = 1010 // Failed to add token
    ErrCacheFileLoadFailed  = 1011 // Failed to load cache file
    ErrCacheFileParseFailed = 1012 // Failed to parse cache file
)
```

### Custom Language Support

Example of adding French language support:

```go
// Register new language
fr := wtoken.registerLanguage("fr")

// Define custom error messages
frenchErrorMessages := map[wtoken.ILanguage]map[wtoken.ErrorCode]string{
    fr: {
        wtoken.ErrSuccess:              "Opération réussie",
        wtoken.ErrInvalidToken:         "Token invalide",
        wtoken.ErrTokenNotFound:        "Token introuvable",
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
    if errData != wtoken.ErrSuccess {
        fmt.Printf("生成token失败：%v\n", errData.Error())
        return
    }

    // API鉴权测试
    errData = manager.Authenticate(tokenKey, "/api/user", "127.0.0.1")
    if errData == wtoken.ErrSuccess {
        fmt.Println("鉴权成功")
    }
}
```

## Token结构

```go
type Token[T any] struct {
    UserID         uint      // 用户ID
    GroupID        uint      // 用户组ID
    LoginTime      time.Time // 登录时间
    ExpireSeconds  int64     // token过期秒数，0表示永不过期
    LastAccessTime time.Time // 最后访问时间
    UserData       T        // 用户自定义数据
    IP             string   // 用户IP地址
}
```

## 错误处理系统

### 内置错误码

```go
const (
    ErrSuccess              = 0    // 操作成功
    ErrUnknown              = 9999 // 未知错误

    // token错误码1101开头
    ErrInvalidToken         = 1101 // 无效的token
    ErrTokenExpired         = 1102 // token已过期
    ErrTokenNotFound        = 1103 // token不存在
    ErrTokenLimitExceeded   = 1104 // 超出token数量限制
    ErrAddToken             = 1105 // 生成token错误

    // 用户错误码，1201开头
    ErrInvalidUserID        = 1201 // 无效的用户ID
    ErrUserIDNotFound       = 1202 // 用户ID不存在
    ErrTypeAssertionError   = 1203 // 类型断言错误

    // 用户组错误码，1301开头
    ErrGroupNotFound        = 1301 // 用户组不存在
    ErrInvalidGroupID       = 1302 // 无效的用户组ID

    // IP错误码，1401开头
    ErrInvalidIP            = 1401 // 无效的IP地址
    ErrIPMismatch           = 1402 // IP地址不匹配

    // 配置错误码，1501开头
    ErrInvalidConfig        = 1501 // 无效的配置

    // 缓存错误码，1601开头
    ErrCacheFileLoadFailed  = 1601 // 加载缓存文件失败
    ErrCacheFileParseFailed = 1602 // 缓存文件解析错误

    // API错误码，1700开头
    ErrAccessDenied         = 1701 // 访问被禁止的API
    ErrInvalidURL           = 1702 // 无效的URL
    ErrNoAPIPermission      = 1703 // 该用户的用户组没有定制API访问权限

    // 内部错误，1901开头
    ErrInternalError        = 1901 // 内部错误
)
```

### 自定义语言支持

添加法语支持示例：

```go
// 注册新语言
fr := wtoken.registerLanguage("fr")

// 定义自定义错误信息
frenchErrorMessages := map[wtoken.ILanguage]map[wtoken.ErrorCode]string{
    fr: {
        wtoken.ErrSuccess:              "Opération réussie",
        wtoken.ErrInvalidToken:         "Token invalide",
        wtoken.ErrTokenNotFound:        "Token introuvable",
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
