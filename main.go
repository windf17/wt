// Package wtoken 提供了一个高性能、线程安全的Token管理系统
// 支持多用户组权限管理、Token生命周期管理、并发访问控制等功能
package wtoken

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// manager 实现类设为私有，防止外部直接访问实例字段
type Manager[T any] struct {
	config *Config
	groups map[uint]*Group
	tokens map[string]*Token[T]
	stats  Stats
	mutex  sync.RWMutex
}

// lock 获取写锁
func (tm *Manager[T]) lock() {
	tm.mutex.Lock()
}

// unlock 释放写锁
func (tm *Manager[T]) unlock() {
	tm.mutex.Unlock()
}

// rLock 获取读锁
func (tm *Manager[T]) rLock() {
	tm.mutex.RLock()
}

// rUnlock 释放读锁
func (tm *Manager[T]) rUnlock() {
	tm.mutex.RUnlock()
}

var (
	once sync.Once
)

// InitTM 初始化Token管理器，并返回全局唯一实例
// 参数：
//   - config: Token管理器配置，如果为nil则使用默认配置
//   - groups: 用户组配置列表，不能为空，否则会触发panic
//   - errorMessages: 用户自定义错误信息映射表，key为语言类型，value为错误码和对应的错误信息
//     如果为nil则只使用默认错误信息，否则会在默认错误信息的基础上追加或覆盖自定义错误信息
//
// 返回：
//   - IManager: 返回实现了IManager接口的manager实例
//   - error: 如果发生错误则返回相应的错误信息
//
// 注意：
//   - 该函数保证返回的是全局唯一的manager实例
//   - 如果配置了缓存文件，会自动从文件加载token数据
//   - 首次调用时会初始化实例，后续调用返回相同实例
//   - 会先注册默认错误信息，然后再注册用户自定义错误信息
//   - 支持多语言错误信息配置，可以根据config中的Language设置来返回对应语言的错误信息
func InitTM[T any](config *Config, groups []GroupRaw, errorMessages map[ILanguage]map[ErrorCode]string) (IManager[T]) {
	var instance *Manager[T]

	once.Do(func() {
		// 配置合并逻辑
		mergedConfig := mergeDefaultConfig(config)

		// 构建用户组
		groupMap := buildGroups(groups, mergedConfig.Delimiter)

		// 初始化实例
		instance = &Manager[T]{
			config: mergedConfig,
			groups: groupMap,
			tokens: make(map[string]*Token[T]),
			stats: Stats{
				LastUpdateTime: time.Now(),
				TotalTokens:    0,
				ActiveTokens:   0,
				ExpiredTokens:  0,
			},
		}

		// 注册错误信息
		registerErrorMessages(errorMessages)

		// 加载缓存文件
		if mergedConfig.CacheFilePath != "" {
			if loadErr := instance.loadFromFile(); loadErr != nil {
				// 记录错误但继续初始化
				if mergedConfig.Debug {
					fmt.Printf("Warning: cache load failed: %v\n", loadErr)
				}
			}
			instance.CleanExpiredTokens()
		}
	})

	// 统一错误处理
	return instance
}

// mergeDefaultConfig 配置合并逻辑（独立函数便于测试）
func mergeDefaultConfig(custom *Config) *Config {
	defaultConfig := &Config{
		CacheFilePath: DefaultCacheFilePath,
		Debug:         DefaultDebug,
		MaxTokens:     DefaultMaxTokens,
		Language:      "zh",
		Delimiter:     DefaultDelimiter,
	}

	if custom == nil {
		return defaultConfig
	}

	// 完整的配置合并逻辑
	if custom.CacheFilePath == "" {
		custom.CacheFilePath = defaultConfig.CacheFilePath
	}
	if custom.MaxTokens <= 0 {
		custom.MaxTokens = defaultConfig.MaxTokens
	}
	if custom.Language == "" {
		custom.Language = defaultConfig.Language
	}
	if custom.Delimiter == "" {
		custom.Delimiter = defaultConfig.Delimiter
	}

	return custom
}

// buildGroups 用户组构建（分离复杂逻辑）
func buildGroups(rawGroups []GroupRaw, delimiter string) map[uint]*Group {
	groups := make(map[uint]*Group)
	for _, raw := range rawGroups {
		group := ConvGroup(raw, delimiter)
		groups[raw.ID] = group
	}
	return groups
}

// saveToFile 将tokens数据异步保存到缓存文件
// 如果配置中未设置缓存文件路径，则不执行保存操作
func (tm *Manager[T]) saveToFile() {
	if tm.config.CacheFilePath == "" {
		return
	}

	// 确保目录存在
	dir := filepath.Dir(tm.config.CacheFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		if tm.config.Debug {
			fmt.Printf("Failed to create directory for cache file: %v\n", err)
		}
		return
	}

	// 准备数据
	data := struct {
		Tokens map[string]*Token[T] `json:"tokens"`
		Stats  Stats                `json:"stats"`
	}{
		Tokens: tm.tokens,
		Stats:  tm.stats,
	}

	// 序列化数据
	fileData, err := json.Marshal(data)
	if err != nil {
		if tm.config.Debug {
			fmt.Printf("Failed to marshal data for cache file: %v\n", err)
		}
		return
	}

	// 写入文件
	if err := os.WriteFile(tm.config.CacheFilePath, fileData, 0644); err != nil {
		if tm.config.Debug {
			fmt.Printf("Failed to write cache file: %v\n", err)
		}
	}
}

// loadFromFile 从缓存文件加载tokens数据
// 如果文件不存在则返回nil
// 如果文件存在但无法读取或解析则返回错误
func (tm *Manager[T]) loadFromFile() error {
	fileData, err := os.ReadFile(tm.config.CacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		if tm.config.Debug {
			fmt.Printf("Failed to read cache file: %v\n", err)
		}
		return tm.NewError(ErrCodeCacheFileLoadFailed)
	}

	var data struct {
		Tokens map[string]*Token[T] `json:"tokens"`
		Stats  Stats                `json:"stats"`
	}

	if err := json.Unmarshal(fileData, &data); err != nil {
		if tm.config.Debug {
			fmt.Printf("Failed to parse cache file: %v\n", err)
		}
		return tm.NewError(ErrCodeCacheFileParseFailed)
	}

	tm.lock()
	tm.tokens = data.Tokens
	tm.stats = data.Stats
	tm.unlock()

	return nil
}

// NewError 创建一个新的错误数据对象
// 根据配置的语言类型返回对应的错误信息
// 注册错误信息
func registerErrorMessages(errorMessages map[ILanguage]map[ErrorCode]string) {
	// 先注册默认错误信息
	defaultRegistry.RegisterDefaultMessages()
	// 如果有用户自定义错误信息，则进行注册
	for lang, messages := range errorMessages {
		for code, message := range messages {
			defaultRegistry.RegisterErrorMessage(lang, code, message)
		}
	}
}

func (m *Manager[T]) NewError(code ErrorCode) ErrorData {
	return ErrorData{
		Lan:  m.config.Language,
		Code: code,
	}
}

// GetData 获取用户数据并进行类型转换
// 支持泛型，可以直接转换为指定类型
// 如果token不存在或已过期则返回错误
func (m *Manager[T]) GetData(key string) (T, error) {
	m.rLock()
	defer m.rUnlock()

	t := m.tokens[key]
	if t == nil {
		var zero T
		return zero, m.NewError(ErrCodeTokenNotFound)
	}
	if t.IsExpired() {
		delete(m.tokens, key)
		var zero T
		return zero, m.NewError(ErrCodeTokenExpired)
	}
	return t.UserData, nil
}

// SaveData 保存用户数据
// 支持泛型，可以保存任意类型的数据
// 如果token不存在或已过期则返回错误
func (m *Manager[T]) SaveData(key string, data T) error {
	m.lock()
	defer m.unlock()

	t := m.tokens[key]
	if t == nil {
		return m.NewError(ErrCodeTokenNotFound)
	}
	if t.IsExpired() {
		delete(m.tokens, key)
		return m.NewError(ErrCodeTokenExpired)
	}
	// 更新数据和访问时间
	t.UserData = data
	t.LastAccessTime = time.Now()
	m.tokens[key] = t
	return nil
}
