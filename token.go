package wtoken

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// GetToken 获取token数据
func (tm *Manager[T]) GetToken(key string) (*Token[T], ErrorCode) {
	tm.rLock()
	defer tm.rUnlock()
	t := tm.tokens[key]
	if t == nil {
		return nil, (ErrTokenNotFound)
	}
	if t.IsExpired() {
		// token已过期，则删除该token
		delete(tm.tokens, key)
		return nil, (ErrTokenExpired)
	}
	return t, ErrSuccess
}

// AddToken 新增token，通过它申请token，不存储用户数据，存储用户数据另外用SaveData
func (tm *Manager[T]) AddToken(userID uint, groupID uint, ip string) (string, ErrorCode) {
	if ip == "" {
		return "", ErrInvalidIP
	}
	if userID == 0 {
		return "", ErrInvalidUserID
	}
	if groupID == 0 {
		return "", ErrInvalidGroupID
	}

	// 首先检查用户组是否存在
	tm.rLock()
	g := tm.groups[groupID]
	if g == nil {
		tm.rUnlock()
		return "", ErrGroupNotFound
	}
	expireSeconds := g.ExpireSeconds
	allowMultipleLogin := g.AllowMultipleLogin
	tm.rUnlock()

	// 生成token
	tokenKey, er := tm.GenerateToken()
	if er != nil {
		return "", ErrAddToken
	}

	// 创建用户tokens数据
	now := time.Now()
	var zero T
	tokenData := Token[T]{
		UserID:         userID,
		GroupID:        groupID,
		LoginTime:      now,
		LastAccessTime: now,
		ExpireSeconds:     expireSeconds,
		UserData:       zero,
		IP:             ip,
	}

	// 如果配置了最大token数量，先清理过期token
	if tm.config.MaxTokens > 0 {
		tm.CleanExpiredTokens()
	}

	// 获取写锁进行token操作
	tm.lock()
	defer tm.unlock()

	// 如果不允许多设备登录，则清理该用户在其他设备上的token
	if !allowMultipleLogin {
		for t, ut := range tm.tokens {
			if ut.UserID == userID && ut.IP != ip {
				delete(tm.tokens, t)
			}
		}
	}

	// 如果仍然超过最大token数量，删除最旧的token
	if tm.config.MaxTokens > 0 && len(tm.tokens) >= tm.config.MaxTokens {
		var oldestKey string
		var oldestTime time.Time = time.Now()
		for t, ut := range tm.tokens {
			if ut.LastAccessTime.Before(oldestTime) {
				oldestKey = t
				oldestTime = ut.LastAccessTime
			}
		}
		if oldestKey != "" {
			delete(tm.tokens, oldestKey)
		}
	}

	// 存储token
	tm.tokens[tokenKey] = &tokenData
	tm.updateStatsCount(1, true)
	// 保存到缓存文件
	go tm.saveToFile() // 异步保存到缓存文件
	return tokenKey, ErrSuccess
}

// DelToken 删除指定的token
func (tm *Manager[T]) DelToken(key string) ErrorCode {
	tm.lock()
	defer tm.unlock()
	if _, exists := tm.tokens[key]; !exists {
		return (ErrTokenNotFound)
	}
	delete(tm.tokens, key)
	tm.updateStatsCount(-1, true)
	// 保存到缓存文件
	go tm.saveToFile()
	return ErrSuccess
}

// DelTokensByUserID 删除指定用户的所有token
func (tm *Manager[T]) DelTokensByUserID(userID uint) ErrorCode {
	if userID == 0 {
		return (ErrInvalidUserID)
	}
	tm.lock()
	defer tm.unlock()
	deleteCount := 0
	for token, ut := range tm.tokens {
		if ut.UserID == userID {
			delete(tm.tokens, token)
			deleteCount++
		}
	}
	if deleteCount > 0 {
		tm.updateStatsCount(-deleteCount, true)
		// 保存到缓存文件
		go tm.saveToFile()
	}
	return ErrSuccess
}

// DelTokensByGroupID 删除指定用户组的所有token
func (tm *Manager[T]) DelTokensByGroupID(groupID uint) ErrorCode {
	if groupID == 0 {
		return (ErrInvalidGroupID)
	}
	tm.lock()
	defer tm.unlock()
	// 检查用户组id是不是存在
	if _, exists := tm.groups[groupID]; !exists {
		return (ErrGroupNotFound)
	}
	deleteCount := 0
	for token, ut := range tm.tokens {
		if ut.GroupID == groupID {
			delete(tm.tokens, token)
			deleteCount++
		}
	}
	if deleteCount > 0 {
		tm.updateStatsCount(-deleteCount, true)
		// 保存到缓存文件
		go tm.saveToFile()
	}
	return ErrSuccess
}

// UpdateToken 更新指定的token
func (tm *Manager[T]) UpdateToken(key string, token *Token[T]) ErrorCode {
	tm.lock()
	defer tm.unlock()
	if _, exists := tm.tokens[key]; !exists {
		return (ErrTokenNotFound)
	}
	if token == nil {
		return (ErrInvalidToken)
	}
	token.LastAccessTime = time.Now()
	tm.tokens[key] = token
	// 保存到缓存文件
	go tm.saveToFile()
	return ErrSuccess
}

// CheckToken 验证token是否有效
func (tm *Manager[T]) CheckToken(key string) ErrorCode {
	tm.rLock()
	defer tm.rUnlock()

	if key == "" {
		return (ErrInvalidToken)
	}
	t, exists := tm.tokens[key]
	if !exists {
		return (ErrTokenNotFound)
	}
	if t == nil {
		// token的key存在但值为nil，则删除该token
		delete(tm.tokens, key)
		return (ErrInvalidToken)
	}
	if t.IsExpired() {
		// token已过期，则删除该token
		delete(tm.tokens, key)
		return (ErrTokenExpired)
	}

	return ErrSuccess
}

// CleanExpiredTokens 清理过期token并更新缓存文件
func (tm *Manager[T]) CleanExpiredTokens() {
	tm.lock()
	defer tm.unlock()

	count := 0
	for key, token := range tm.tokens {
		if token == nil {
			delete(tm.tokens, key)
			count++
		} else if token.IsExpired() {
			delete(tm.tokens, key)
			count++
		}
	}

	if count > 0 {
		tm.updateStatsCount(-count, true)
		// 保存到缓存文件
		go tm.saveToFile()
	}
}

// GenerateToken：生成随机token
func (tm *Manager[T]) GenerateToken() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	now := time.Now().UnixNano()
	nowBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		nowBytes[i] = byte(now >> uint(i*8))
	}
	b = append(nowBytes, b...)
	return base64.URLEncoding.EncodeToString(b), nil
}
