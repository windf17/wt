package wtoken

// GetGroup 获取并验证用户组配置
func (tm *Manager[T]) GetGroup(groupID uint) (*Group, ErrorData) {
	tm.rLock()
	defer tm.rUnlock()
	g := tm.groups[groupID]
	if g == nil {
		return nil, tm.NewError(ErrCodeGroupNotFound)
	}

	return g, tm.NewError(ErrCodeSuccess)
}

// AddGroup 新增用户组
func (tm *Manager[T]) AddGroup(raw GroupRaw) ErrorData {
	if raw.ID == 0 {
		return tm.NewError(ErrCodeInvalidGroupID)
	}
	tm.lock()
	defer tm.unlock()
	group := ConvGroup(raw, tm.config.Delimiter)
	tm.groups[raw.ID] = group
	return tm.NewError(ErrCodeSuccess)
}

// DeleteGroup 删除指定用户组的所有token
func (tm *Manager[T]) DelGroup(groupID uint) ErrorData {
	if groupID == 0 {
		return tm.NewError(ErrCodeGroupNotFound)
	}
	tm.lock()
	deleted := false
	for token, ut := range tm.tokens {
		if ut.GroupID == groupID {
			delete(tm.tokens, token)
			deleted = true
		}
	}
	tm.unlock()
	if deleted {
		tm.UpdateStats()
		// 保存到缓存文件
		go tm.saveToFile()
	}

	return tm.NewError(ErrCodeSuccess)
}

// UpdateGroup 更新用户组
func (tm *Manager[T]) UpdateGroup(raw GroupRaw) ErrorData {
	if raw.ID == 0 {
		return tm.NewError(ErrCodeInvalidGroupID)
	}
	tm.lock()
	defer tm.unlock()
	_, exists := tm.groups[raw.ID]
	if !exists {
		return tm.NewError(ErrCodeGroupNotFound)
	}
	group := ConvGroup(raw, tm.config.Delimiter)
	tm.groups[raw.ID] = group
	return tm.NewError(ErrCodeSuccess)
}
