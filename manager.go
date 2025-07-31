package wt

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/windf17/wt/models"
)

// Manager Token管理器结构体
type Manager[T any] struct {
	// tokens 存储所有token
	tokens map[string]*models.Token[T]
	// groups 存储所有用户组
	groups map[uint]*models.Group
	// config 配置信息
	config *models.Config
	// mu 读写锁
	mu sync.RWMutex
	// stats 统计信息
	stats models.Stats
}

/**
 * InitTM 初始化Token管理器
 * @param {*ConfigRaw} config 配置信息
 * @param {[]models.GroupRaw} groups 用户组配置
 * @returns {IManager[T]} Token管理器实例
 */
func InitTM[T any](config models.ConfigRaw, groups []models.GroupRaw) (models.IManager[T], error) {
	// 验证配置
	if err := ValidateConfig(config); err != nil {
		// 配置无效，返回nil
		return nil, err
	}

	// 转换配置
	cfg := &models.Config{
		Language:       config.Language,
		MaxTokens:      config.MaxTokens,
		Delimiter:      config.Delimiter,
		TokenRenewTime: parseTokenRenewTime(config.TokenRenewTime),
	}

	// 创建管理器实例
	tm := &Manager[T]{
		tokens: make(map[string]*models.Token[T]),
		groups: make(map[uint]*models.Group),
		config: cfg,
		stats:  models.Stats{LastUpdateTime: time.Now()},
	}

	// 添加用户组（如果提供了groups）
	if len(groups) > 0 {
		for _, group := range groups {
			if err := ValidateGroupRaw(group); err != nil {
				// 用户组验证失败，返回错误
				return nil, err
			}
			tm.AddGroup(&group)
		}
	}
	return tm, nil
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
	tm.mu.Lock()
}

func (tm *Manager[T]) unlock() {
	tm.mu.Unlock()
}

func (tm *Manager[T]) rLock() {
	tm.mu.RLock()
}

func (tm *Manager[T]) rUnlock() {
	tm.mu.RUnlock()
}

// SetUserData 设置用户数据
func (tm *Manager[T]) SetUserData(key string, data T) error {
	// 输入参数验证
	if key == "" {
		return errors.New(getErrorMessage(tm.config.Language, "invalid_token"))
	}

	tm.lock()
	defer tm.unlock()

	// 检查token是否存在
	token, exists := tm.tokens[key]
	if !exists {
		return errors.New(getErrorMessage(tm.config.Language, "invalid_token"))
	}

	// 检查token是否过期
	if token.IsExpired() {
		delete(tm.tokens, key)
		return errors.New(getErrorMessage(tm.config.Language, "token_expired"))
	}

	// 设置用户数据
	token.UserData = data

	// 更新访问时间
	token.LastAccessTime = time.Now()

	return nil
}

// GetUserData 获取用户数据
func (tm *Manager[T]) GetUserData(key string) (T, error) {
	// 先用读锁检查token
	tm.rLock()
	var zeroValue T

	// 检查token是否存在
	token, exists := tm.tokens[key]
	if !exists {
		tm.rUnlock()
		return zeroValue, errors.New(getErrorMessage(tm.config.Language, "invalid_token"))
	}

	// 检查token是否过期
	if token.IsExpired() {
		tm.rUnlock()
		// 使用写锁删除过期token
		tm.lock()
		delete(tm.tokens, key)
		tm.unlock()
		return zeroValue, errors.New(getErrorMessage(tm.config.Language, "token_expired"))
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

	return userData, nil
}
