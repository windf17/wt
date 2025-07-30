# wt v1.0.0 - 企业级高性能 Token 管理系统

## 🚀 核心特性

### 🔐 安全认证

-   **JWT Token 管理**：基于行业标准的 JWT 实现
-   **多层权限控制**：支持用户组权限管理
-   **安全加密**：采用高强度加密算法保护 Token 安全
-   **IP 绑定验证**：可选的 IP 地址绑定增强安全性

### ⚡ 高性能设计

-   **内存优化**：高效的内存管理和缓存机制
-   **并发安全**：完全的并发安全设计
-   **批量操作**：支持批量 Token 创建和管理
-   **快速验证**：毫秒级 Token 验证响应

### 🛠️ 企业级功能

-   **持久化存储**：支持文件缓存和数据持久化
-   **统计监控**：实时 Token 使用统计和监控
-   **错误处理**：完善的错误处理和国际化支持
-   **配置灵活**：丰富的配置选项满足不同需求

## 📦 安装指南

### Go 模块安装

```bash
go get github.com/windf17/wt@v1.0.0
```

### 基本使用

```go
package main

import (
    "fmt"
    "time"
    "github.com/windf17/wt"
)

func main() {
    // 创建配置
    config := &wt.Config{
        Language: wt.LangChinese,
        Debug: false,
    }

    // 创建管理器
    manager, err := wt.GetTM(config, nil)
    if err != nil {
        panic(err)
    }

    // 创建Token
    token, err := manager.CreateToken("user123", time.Hour*24)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Token创建成功: %s\n", token)
}
```

## 📊 性能指标

-   **Token 创建速度**：> 10,000 tokens/秒
-   **Token 验证速度**：> 50,000 验证/秒
-   **内存占用**：< 10MB (10 万 Token)
-   **并发支持**：1000+ 并发连接

## 🎯 快速开始

### 1. 基础 Token 管理

```go
// 创建Token
token, err := manager.CreateToken(userID, expiration)

// 验证Token
userID, err := manager.ValidateToken(token)

// 撤销Token
err = manager.RevokeToken(token)
```

### 2. 权限组管理

```go
// 定义权限组
groups := []wt.GroupRaw{
    {
        ID: 1,
        AllowedAPIs: "/api/user /api/profile",
        TokenExpire: "86400", // 24小时
    },
}

// 创建带权限的管理器
manager, err := wt.GetTM(config, groups)
```

### 3. 批量操作

```go
// 批量创建Token
userIDs := []string{"user1", "user2", "user3"}
tokens, err := manager.CreateBatchTokens(userIDs, time.Hour*12)
```

## 📚 文档链接

-   [完整 API 文档](https://pkg.go.dev/github.com/windf17/wt)
-   [使用示例](https://github.com/windf17/wt/tree/main/examples)
-   [性能测试报告](https://github.com/windf17/wt/tree/main/test)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进 wt！

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

---

**wt v1.0.0** - 为现代应用提供安全、高效的 Token 管理解决方案！

## 版本说明

**v1.0.0** - 首个正式发布版本，提供完整的企业级 Token 管理功能
