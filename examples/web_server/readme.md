# wt Web 服务器示例

这是一个完整的 wt Web 服务器示例，演示了如何在实际 Web 应用中集成 wt 进行用户认证和权限管理。

## 📋 示例概述

本示例实现了一个完整的用户权限管理系统，包括：

-   **用户认证**: 登录/登出功能
-   **权限控制**: 基于用户组的 API 访问控制
-   **安全防护**: Token 验证、IP 绑定、权限检查
-   **RESTful API**: 标准的 REST API 设计

## 🏗️ 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │───►│  Auth Middleware │───►│   API Handler   │
│   (用户请求)     │    │   (认证中间件)   │    │   (业务处理)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │  wt Manager │
                       │  (Token管理器)   │
                       └─────────────────┘
```

## 👥 用户组配置

### 管理员组 (ID=1)

-   **权限**: 不允许重复登录
-   **API 访问**: 拥有管理员专属 API 权限
-   **安全级别**: 高

### 普通用户组 (ID=2)

-   **权限**: 允许重复登录
-   **API 访问**: 拥有普通用户专属 API 权限
-   **安全级别**: 标准

## 👤 测试用户账号

### 管理员用户

-   `admin1` / `admin123`
-   `admin2` / `admin123`

### 普通用户

-   `user1` / `user123` 到 `user100` / `user123`

## 🔗 API 端点

| 方法 | 端点                   | 描述         | 权限要求     |
| ---- | ---------------------- | ------------ | ------------ |
| POST | `/api/login`           | 用户登录     | 无           |
| POST | `/api/logout`          | 用户登出     | 需要 Token   |
| GET  | `/api/admin/dashboard` | 管理员仪表板 | 管理员权限   |
| GET  | `/api/user/profile`    | 用户个人资料 | 普通用户权限 |

## 🚀 快速开始

### 1. 启动服务器

```bash
# 进入示例目录
cd examples/web_server

# 运行服务器
go run main.go
```

服务器将在 `http://localhost:8081` 启动

### 2. 测试 API

#### 管理员登录

```bash
curl -X POST http://localhost:8081/api/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin1","password":"admin123"}'
```

**响应示例**:

```json
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_at": "2024-12-20T10:30:00Z",
        "user_info": {
            "user_id": 1,
            "username": "admin1",
            "email": "admin1@example.com",
            "role": "管理员"
        }
    }
}
```

#### 普通用户登录

```bash
curl -X POST http://localhost:8081/api/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"user1","password":"user123"}'
```

#### 访问管理员 API

```bash
curl -X GET http://localhost:8081/api/admin/dashboard \
  -H 'Authorization: Bearer YOUR_ADMIN_TOKEN'
```

#### 访问用户 API

```bash
curl -X GET http://localhost:8081/api/user/profile \
  -H 'Authorization: Bearer YOUR_USER_TOKEN'
```

#### 用户登出

```bash
curl -X POST http://localhost:8081/api/logout \
  -H 'Authorization: Bearer YOUR_TOKEN'
```

## 🔒 安全特性

### 1. Token 验证

-   **格式验证**: 严格的 Token 格式检查
-   **过期检查**: 自动检查 Token 是否过期
-   **IP 绑定**: Token 与客户端 IP 绑定，防止盗用

### 2. 权限控制

-   **角色分离**: 管理员和普通用户权限分离
-   **API 级权限**: 细粒度的 API 访问控制
-   **自动拒绝**: 未授权访问自动拒绝

### 3. 并发安全

-   **线程安全**: 支持高并发访问
-   **死锁预防**: 精心设计的锁机制
-   **性能优化**: 读写锁分离

## 📊 性能测试

### 压力测试

本示例包含完整的压力测试工具：

```bash
# 运行压力测试
cd test
go run pressure_testing.go
```

### 测试结果

经过 1000 并发用户测试验证：

-   ✅ **登录成功率**: 100% (1000/1000)
-   ✅ **登出成功率**: 100% (1000/1000)
-   ✅ **API 请求成功率**: 100% (授权请求)
-   ✅ **权限控制**: 100%正确拒绝未授权请求
-   ⚡ **平均响应时间**: <10ms
-   🚀 **吞吐量**: 509.68 请求/秒

## 🧪 测试用例

### 功能测试

1. **用户认证测试**

    - 正确用户名密码登录
    - 错误用户名密码登录
    - Token 过期处理

2. **权限控制测试**

    - 管理员访问管理员 API
    - 普通用户访问用户 API
    - 跨权限访问拒绝

3. **安全测试**
    - Token 盗用检测
    - IP 验证
    - 并发登录控制

### 压力测试

```go
// 压力测试配置
const (
    baseURL      = "http://localhost:8081"
    testUsers    = 1000  // 并发用户数
    testDuration = 120   // 测试时长(秒)
)
```

## 📁 文件结构

```
examples/web_server/
├── main.go              # 主服务器文件
├── go.mod              # Go模块依赖
├── README.md           # 本文档
└── test/
    └── pressure_testing.go  # 压力测试工具
```

## 🔧 配置说明

### Token 管理器配置

```go
config := &wt.ConfigRaw{
    MaxTokens:      10000,           // 最大Token数量
    Delimiter:      ",",             // API权限分隔符
    TokenRenewTime: "24h",           // Token续期时间
    Language:       wt.Chinese,  // 错误信息语言
}
```

### 用户组配置

```go
groups := []wt.GroupRaw{
    {
        ID:                 1,
        Name:               "管理员",
        ExpireSeconds:      3600,        // 1小时过期
        AllowMultipleLogin: 0,           // 不允许多设备登录
        AllowedAPIs:        "/api/admin,/api/user", // 允许的API
    },
    {
        ID:                 2,
        Name:               "普通用户",
        ExpireSeconds:      7200,        // 2小时过期
        AllowMultipleLogin: 1,           // 允许多设备登录
        AllowedAPIs:        "/api/user",    // 允许的API
    },
}
```

## 🚨 错误处理

### 常见错误码

| 错误码 | HTTP 状态码 | 描述       | 解决方案            |
| ------ | ----------- | ---------- | ------------------- |
| 2001   | 401         | 未授权访问 | 检查 Token 是否有效 |
| 2002   | 403         | 权限不足   | 检查用户权限        |
| 2101   | 401         | 无效 Token | 重新登录获取 Token  |
| 2102   | 401         | Token 过期 | 重新登录            |

### 错误响应格式

```json
{
    "code": 401,
    "message": "认证失败: 无效令牌",
    "data": null
}
```

## 🔍 调试技巧

### 1. 启用详细日志

服务器启动时会显示详细的初始化信息和 API 端点。

### 2. 检查 Token 状态

```bash
# 使用无效Token测试
curl -X GET http://localhost:8081/api/user/profile \
  -H 'Authorization: Bearer invalid_token'
```

### 3. 监控服务器日志

服务器会输出详细的请求处理日志，包括：

-   客户端 IP
-   访问路径
-   Token 验证结果
-   权限检查结果

## 🎯 最佳实践

### 1. 安全建议

-   使用 HTTPS 传输 Token
-   定期轮换 Token
-   实施 IP 白名单
-   监控异常访问

### 2. 性能优化

-   合理设置 Token 过期时间
-   使用连接池
-   启用缓存
-   监控系统指标

### 3. 错误处理

-   统一错误响应格式
-   详细的错误日志
-   用户友好的错误信息
-   异常情况的降级处理

## 📈 扩展功能

### 可扩展的功能点

1. **数据库集成**: 用户信息持久化存储
2. **缓存优化**: Redis 缓存 Token
3. **监控告警**: 集成 Prometheus 监控
4. **负载均衡**: 多实例部署
5. **API 网关**: 集成到 API 网关

### 集成示例

```go
// 集成数据库
func initDatabase() {
    // 数据库连接和初始化
}

// 集成缓存
func initRedis() {
    // Redis连接和配置
}

// 集成监控
func initMetrics() {
    // Prometheus指标注册
}
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个示例！

## 📞 支持

如果您在使用过程中遇到问题，请：

1. 查看本文档的常见问题
2. 检查服务器日志
3. 提交 Issue 到项目仓库

---

**这个示例展示了 wt 在实际 Web 应用中的强大功能和易用性！** 🚀
