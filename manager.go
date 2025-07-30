package wt

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/windf17/wt/models"
)

// Manager Token管理器结构体
type Manager[T any] struct {
	// tokens 存储所有token的映射
	tokens map[string]*Token[T]
	// groups 存储所有用户组的映射
	groups map[uint]*models.Group
	// config 配置信息
	config *Config
	// mutex 读写锁，保证并发安全
	mutex sync.RWMutex
	// stats 统计信息
	stats Stats

	// errorHandler 错误处理器
	errorHandler *ErrorHandler
}

// 全局单例实例
var (
	instance     any
	instanceLock sync.Once
)

// GetInstance 获取Token管理器单例实例
func GetInstance[T any]() IManager[T] {
	instanceLock.Do(func() {
		if instance == nil {
			// 使用默认配置初始化
			instance = InitTM[T](DefaultConfigRaw, []models.GroupRaw{}, nil)
		}
	})
	return instance.(IManager[T])
}

// InitTM 初始化Token管理器
func InitTM[T any](config *ConfigRaw, groups []models.GroupRaw, errorMessages map[Language]map[ErrorCode]string) IManager[T] {
	// 验证配置
	if err := ValidateConfig(config); err != nil {
		// 使用错误处理器记录错误并返回nil
		eh := NewErrorHandler()
		eh.HandleError(E_ConfigInvalid, fmt.Sprintf("Invalid config: %v", err), ERROR_LEVEL_FATAL, nil)
		return nil
	}

	// 转换配置
	cfg := convertConfig(config)

	// 创建管理器实例
	tm := &Manager[T]{
		tokens: make(map[string]*Token[T]),
		groups: make(map[uint]*models.Group),
		config: cfg,
		stats:  Stats{LastUpdateTime: time.Now()},
	}

	// 初始化组件
	tm.errorHandler = NewErrorHandler()

	// 设置错误消息
	if errorMessages != nil {
		SetErrorMessages(errorMessages)
		// 设置语言为配置中指定的语言
		if cfg.Language != "" {
			SetLanguage(cfg.Language)
		}
	}

	// 添加用户组（如果提供了groups）
	if len(groups) > 0 {
		for _, group := range groups {
			if err := ValidateGroupRaw(group); err != nil {
				// 使用错误处理器记录错误并返回nil
				tm.errorHandler.HandleError(E_GroupInvalid, fmt.Sprintf("Invalid group config: %v", err), ERROR_LEVEL_FATAL, nil)
				return nil
			}
			tm.AddGroup(&group)
		}
	}
	// 注意：当groups为空时，鉴权功能将被禁用，所有请求都会通过

	return tm
}

// convertConfig 转换配置
func convertConfig(raw *ConfigRaw) *Config {
	if raw == nil {
		raw = DefaultConfigRaw
	}
	return &Config{
		Language:       raw.Language,
		MaxTokens:      raw.MaxTokens,
		Delimiter:      raw.Delimiter,
		TokenRenewTime: parseTokenRenewTime(raw.TokenRenewTime),
		MinTokenExpire: raw.MinTokenExpire,
		MaxTokenExpire: raw.MaxTokenExpire,
	}
}

// parseTokenRenewTime 解析Token续期时间
func parseTokenRenewTime(renewTime string) int64 {
	if renewTime == "" {
		return 600 // 默认10分钟
	}

	// 获取单位
	unit := renewTime[len(renewTime)-1]
	valueStr := renewTime[:len(renewTime)-1]

	var value int64
	fmt.Sscanf(valueStr, "%d", &value)

	switch unit {
	case 'h', 'H': // 小时
		return value * SECONDS_PER_HOUR
	case 'm', 'M': // 分钟
		return value * SECONDS_PER_MINUTE
	case 'd', 'D': // 天
		return value * SECONDS_PER_DAY
	case 's', 'S': // 秒
		return value
	default:
		return value
	}
}

// 锁操作方法
func (tm *Manager[T]) lock() {
	tm.mutex.Lock()
}

func (tm *Manager[T]) unlock() {
	tm.mutex.Unlock()
}

func (tm *Manager[T]) rLock() {
	tm.mutex.RLock()
}

func (tm *Manager[T]) rUnlock() {
	tm.mutex.RUnlock()
}

// SetErrorMessages 设置全局错误消息
func SetErrorMessages(errorMessages map[Language]map[ErrorCode]string) {
	defaultRegistry.setErrorMessages(errorMessages)
}

// SetLanguage 设置全局默认语言
func SetLanguage(lang Language) {
	defaultRegistry.setLanguage(lang)
}

// RotateLog 轮转日志文件（保留接口兼容性）
func (tm *Manager[T]) RotateLog() error {
	// 简化实现，不支持日志轮转
	return nil
}

/**
 * EnableLogging 启用或禁用日志记录
 * @param {bool} enabled 是否启用日志
 */
func (tm *Manager[T]) EnableLogging(enabled bool) {
	// 简化实现，保持接口兼容性
	// 实际的日志控制通过LogPath和LogLevel配置
}

/**
 * SetLogFile 设置日志文件路径
 * @param {string} filePath 日志文件路径
 * @returns {error} 错误信息
 */
func (tm *Manager[T]) SetLogFile(filePath string) error {
	// 简化实现，保持接口兼容性
	// 实际的日志文件路径通过配置文件设置
	return nil
}

// SetUserData 设置用户数据
func (tm *Manager[T]) SetUserData(key string, data T) ErrorCode {
	// 输入参数验证
	if strings.TrimSpace(key) == "" {
		return E_InvalidToken
	}

	tm.lock()
	defer tm.unlock()

	// 检查token是否存在
	token, exists := tm.tokens[key]
	if !exists {
		return E_InvalidToken
	}

	// 检查token是否过期
	if token.IsExpired() {
		delete(tm.tokens, key)
		return E_TokenExpired
	}

	// 设置用户数据
	token.UserData = data

	// 更新访问时间
	token.LastAccessTime = time.Now()

	return E_Success
}

// GetUserData 获取用户数据
func (tm *Manager[T]) GetUserData(key string) (T, ErrorCode) {
	// 先用读锁检查token
	tm.rLock()
	var zeroValue T

	// 检查token是否存在
	token, exists := tm.tokens[key]
	if !exists {
		tm.rUnlock()
		return zeroValue, E_InvalidToken
	}

	// 检查token是否过期
	if token.IsExpired() {
		tm.rUnlock()
		// 使用写锁删除过期token
		tm.lock()
		delete(tm.tokens, key)
		tm.unlock()
		return zeroValue, E_TokenExpired
	}

	// 获取用户数据
	userData := token.UserData
	tm.rUnlock()

	// 使用写锁更新访问时间
	tm.lock()
	// 再次检查token是否仍然存在（避免竞态条件）
	if token, exists := tm.tokens[key]; exists && !token.IsExpired() {
		token.LastAccessTime = time.Now()
	}
	tm.unlock()

	return userData, E_Success
}

/**
 * Close 关闭Token管理器，清理资源
 */
func (tm *Manager[T]) Close() {
	// 清理资源
	// 标准log.Logger不需要显式关闭
	// 日志记录器会自动处理资源释放
}
