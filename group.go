package wtoken

import (
	"sort"
	"strings"

	"github.com/windf17/wtoken/models"
	"github.com/windf17/wtoken/utility"
)

// GetGroup 获取并验证用户组配置
func (tm *Manager[T]) GetGroup(groupID uint) (*models.Group, ErrorCode) {
	tm.rLock()
	defer tm.rUnlock()
	g := tm.groups[groupID]
	if g == nil {
		return nil, (E_GroupNotFound)
	}

	return g, E_Success
}

// AddGroup 新增用户组
func (tm *Manager[T]) AddGroup(raw *models.GroupRaw) ErrorCode {
	if raw.ID == 0 {
		return (E_GroupInvalid)
	}
	tm.lock()
	defer tm.unlock()
	group := ConvGroup(*raw, tm.config.Delimiter)
	tm.groups[raw.ID] = group

	return E_Success
}

// DelGroup 删除指定用户组及其所有token
func (tm *Manager[T]) DelGroup(groupID uint) ErrorCode {
	if groupID == 0 {
		return (E_GroupNotFound)
	}
	tm.lock()
	defer tm.unlock()
	
	// 检查用户组是否存在
	if _, exists := tm.groups[groupID]; !exists {
		return (E_GroupNotFound)
	}
	
	// 删除该用户组的所有token
	deleteCount := 0
	for token, ut := range tm.tokens {
		if ut.GroupID == groupID {
			delete(tm.tokens, token)
			deleteCount++
		}
	}
	
	// 删除用户组本身
	delete(tm.groups, groupID)
	
	if deleteCount > 0 {
		tm.updateStatsCount(-deleteCount, true)
	}


	return E_Success
}

// UpdateGroup 更新用户组
func (tm *Manager[T]) UpdateGroup(groupID uint, raw *models.GroupRaw) ErrorCode {

	if raw.ID == 0 {
		return (E_GroupInvalid)
	}
	tm.lock()
	defer tm.unlock()
	_, exists := tm.groups[groupID]
	if !exists {
		return (E_GroupNotFound)
	}
	group := ConvGroup(*raw, tm.config.Delimiter)
	tm.groups[raw.ID] = group

	return E_Success
}

/**
 * UpdateAllGroup 批量更新所有用户组
 * @param {[]models.GroupRaw} groups 用户组原始数据列表
 * @returns {ErrorCode} 操作结果错误码
 */
func (tm *Manager[T]) UpdateAllGroup(groups []models.GroupRaw) ErrorCode {
	// 验证所有用户组配置
	for _, group := range groups {
		if group.ID == 0 {
			return E_GroupInvalid
		}
	}

	tm.lock()
	defer tm.unlock()

	// 清空现有的用户组
	tm.groups = make(map[uint]*models.Group)

	// 添加新的用户组
	for _, raw := range groups {
		group := ConvGroup(raw, tm.config.Delimiter)
		tm.groups[raw.ID] = group
	}


	return E_Success
}



/**
 * ConvGroup 将GroupRaw转换为Group
 * @param {GroupRaw} raw 原始用户组数据
 * @param {string} delimiter API分隔符
 * @returns {*Group} 转换后的用户组对象
 */
func ConvGroup(raw models.GroupRaw, delimiter string) *models.Group {
	g := models.Group{}

	// 处理 AllowMultipleLogin
	if raw.AllowMultipleLogin == 1 {
		g.AllowMultipleLogin = true
	} else {
		g.AllowMultipleLogin = false
	}
	// 处理 Name
	g.Name = raw.Name
	// 处理 TokenExpire
	g.ExpireSeconds = utility.ParseDuration(raw.TokenExpire)
	rules := []models.ApiRule{}
	
	// 处理拒绝的规则
	for _, api := range strings.Split(raw.DeniedAPIs, delimiter) {
		api = strings.TrimSpace(api)
		if api != "" {
			apis := utility.ParsePathToSegments(api)
			if len(apis) > 0 {
				rules = append(rules, models.ApiRule{
					Path: apis,
					Rule: false,
				})
			}
		}
	}
	
	// 处理允许的规则
	for _, api := range strings.Split(raw.AllowedAPIs, delimiter) {
		api = strings.TrimSpace(api)
		if api != "" {
			apis := utility.ParsePathToSegments(api)
			if len(apis) > 0 {
				rules = append(rules, models.ApiRule{
					Path: apis,
					Rule: true,
				})
			}
		}
	}
	
	// 对规则进行复杂排序
	sort.Slice(rules, func(i, j int) bool {
		return compareApiRules(rules[i].Path, rules[j].Path)
	})
	
	g.ApiRules = rules
	return &g
}

/**
 * compareApiRules 比较两个API规则路径的优先级
 * 排序规则：
 * 1. 路径段数组长度越长排越前
 * 2. 长度相同时逐个比较字符串长度，长的排前
 * 3. 字符串长度也相同时，先遇到的排前（稳定排序）
 * 4. 特例：*号优先级最低排最后
 * @param {[]string} pathA 第一个路径段数组
 * @param {[]string} pathB 第二个路径段数组
 * @returns {bool} 如果pathA应该排在pathB前面返回true
 */
func compareApiRules(pathA, pathB []string) bool {
	// 1. 首先比较路径段数组长度，长度越长优先级越高
	if len(pathA) != len(pathB) {
		return len(pathA) > len(pathB)
	}
	
	// 2. 长度相同时，逐个比较每个路径段
	for i := 0; i < len(pathA); i++ {
		segmentA := pathA[i]
		segmentB := pathB[i]
		
		// 特例处理：*号优先级最低
		if segmentA == "*" && segmentB != "*" {
			return false // A是*，B不是*，B优先级更高
		}
		if segmentA != "*" && segmentB == "*" {
			return true // A不是*，B是*，A优先级更高
		}
		if segmentA == "*" && segmentB == "*" {
			continue // 都是*，继续比较下一个段
		}
		
		// 比较字符串长度，长度越长优先级越高
		if len(segmentA) != len(segmentB) {
			return len(segmentA) > len(segmentB)
		}
		
		// 长度相同时，字典序比较（保证稳定排序）
		if segmentA != segmentB {
			return segmentA < segmentB
		}
	}
	
	// 所有段都相同，保持原有顺序（稳定排序）
	return false
}

