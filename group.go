package wtoken

// GetGroup 获取并验证用户组配置
func (tm *Manager[T]) GetGroup(groupID uint) (*Group, ErrorCode) {
	tm.rLock()
	defer tm.rUnlock()
	g := tm.groups[groupID]
	if g == nil {
		return nil, (ErrGroupNotFound)
	}

	return g, ErrSuccess
}

// AddGroup 新增用户组
func (tm *Manager[T]) AddGroup(raw GroupRaw) ErrorCode {
	if raw.ID == 0 {
		return (ErrInvalidGroupID)
	}
	tm.lock()
	defer tm.unlock()
	group := ConvGroup(raw, tm.config.Delimiter)
	tm.groups[raw.ID] = group
	return ErrSuccess
}

// DeleteGroup 删除指定用户组的所有token
func (tm *Manager[T]) DelGroup(groupID uint) ErrorCode {
	if groupID == 0 {
		return (ErrGroupNotFound)
	}
	tm.lock()
	deleteCount := 0
	for token, ut := range tm.tokens {
		if ut.GroupID == groupID {
			delete(tm.tokens, token)
			deleteCount++
		}
	}
	tm.unlock()
	if deleteCount > 0 {
		tm.updateStatsCount(-deleteCount, true)
		// 保存到缓存文件
		go tm.saveToFile()
	}

	return ErrSuccess
}

// UpdateGroup 更新用户组
func (tm *Manager[T]) UpdateGroup(raw GroupRaw) ErrorCode {
	if raw.ID == 0 {
		return (ErrInvalidGroupID)
	}
	tm.lock()
	defer tm.unlock()
	_, exists := tm.groups[raw.ID]
	if !exists {
		return (ErrGroupNotFound)
	}
	group := ConvGroup(raw, tm.config.Delimiter)
	tm.groups[raw.ID] = group
	return ErrSuccess
}
