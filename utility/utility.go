package utility

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/windf17/wt/models"
)

/**
 * ParseURLToPathSegments 解析URL并返回路径段数组
 * 将完整的URL解析为路径段，去除空段和根路径
 * @param {string} urlStr 完整的URL字符串（如："https://example.com/api/v1/users"）
 * @returns {[]string} 非空路径段数组，例如：["api", "v1", "users"]
 *                    如果URL解析失败或路径为空，返回空数组
 */
func ParseURLToPathSegments(urlStr string) []string {
	// 输入验证
	if urlStr == "" {
		return []string{}
	}

	// 解析URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return []string{}
	}

	// 提取并解析路径部分
	return ParsePathToSegments(parsedURL.Path)
}

/**
 * ParsePathToSegments 解析路径字符串为路径段数组
 * 将路径字符串分割为非空的路径段数组
 * @param {string} path 路径字符串（如："/api/v1/users" 或 "api/v1/users"）
 * @returns {[]string} 非空路径段数组，例如：["api", "v1", "users"]
 *                    如果路径为空或只包含根路径，返回空数组
 */
func ParsePathToSegments(path string) []string {
	// 输入验证：空字符串或只有根路径
	if path == "" || path == "/" {
		return []string{}
	}

	// 移除开头和结尾的斜杠，标准化路径
	path = strings.Trim(path, "/")

	// 再次检查处理后的路径是否为空
	if path == "" {
		return []string{}
	}

	// 按斜杠分割路径
	parts := strings.Split(path, "/")

	// 过滤掉空字符串段（处理连续斜杠的情况）
	var result []string
	for _, part := range parts {
		// 去除空白字符并检查是否为空
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}

/**
 * HasPermission 检查API权限 - wt核心权限验证算法
 *
 * 这是wt系统的核心权限验证函数，实现了基于路径前缀匹配的权限控制算法。
 * 该算法采用"最长匹配优先"原则，确保权限控制的精确性和安全性。
 *
 * @param {string} urlStr 请求的URL字符串（支持完整URL，会自动解析路径部分）
 * @param {[]models.ApiRule} apiRules API规则数组（已按优先级排序）
 *
 * 核心算法原理：
 *
 * 1. 路径解析阶段：
 *    - 将URL解析为标准化的路径段数组
 *    - 自动处理查询参数、锚点、URL编码等
 *    - 过滤空段，确保路径标准化
 *    - 示例："/api/v1/users?id=123" → ["api", "v1", "users"]
 *
 * 2. 前缀匹配算法：
 *    - 采用从左到右的精确字符串匹配
 *    - 规则路径必须是请求路径的完整前缀
 *    - 不支持通配符或正则表达式匹配
 *    - 匹配过程：逐段比较，遇到不匹配立即停止
 *
 * 3. 最长匹配优先原则：
 *    - 在所有匹配的规则中，选择路径段数最多的规则
 *    - 确保更具体的规则优先于更通用的规则
 *    - 避免权限泄露和误判
 *
 * 4. 匹配示例详解：
 *    请求路径：/api/v1/users/profile
 *
 *    规则集合：
 *    - Rule1: ["api"] (Rule: true)           → 匹配1段
 *    - Rule2: ["api", "v1"] (Rule: true)     → 匹配2段
 *    - Rule3: ["api", "v1", "users"] (Rule: false) → 匹配3段 ★最长匹配
 *    - Rule4: ["api", "v2"] (Rule: true)     → 匹配0段（v2≠v1）
 *
 *    结果：选择Rule3，返回false（拒绝访问）
 *
 * 5. 边界情况处理：
 *    - 空路径：直接返回false
 *    - 无匹配规则：默认拒绝访问（安全优先）
 *    - 规则路径长于请求路径：不匹配
 *    - URL编码：自动解码处理
 *
 * 6. 性能特性：
 *    - 时间复杂度：O(n*m)，n为规则数，m为平均路径长度
 *    - 空间复杂度：O(1)，原地匹配
 *    - 早期终止：遇到不匹配立即停止
 *
 * 7. 安全保证：
 *    - 默认拒绝策略：无匹配规则时拒绝访问
 *    - 精确匹配：防止路径遍历攻击
 *    - 最长匹配：防止权限泄露
 *
 * @returns {bool} 权限验证结果（true=允许访问，false=拒绝访问）
 *
 * @example
 * // 基本使用
 * rules := []models.ApiRule{
 *     {Path: ["api", "user"], Rule: true},
 *     {Path: ["api", "admin"], Rule: false},
 * }
 *
 * hasPermission := HasPermission("/api/user/profile", rules) // true
 * hasPermission = HasPermission("/api/admin/delete", rules)   // false
 *
 * @see ParseURLToPathSegments 路径解析函数
 * @see models.ApiRule API规则结构定义
 */
func HasPermission(urlStr string, apiRules []models.ApiRule) bool {
	// 解析请求路径为路径段数组
	apiPath := ParseURLToPathSegments(urlStr)
	if len(apiPath) == 0 {
		return false
	}

	// 找到有效匹配的规则中路径最长的规则
	var bestRule *models.ApiRule
	maxRuleLength := 0

	for _, rule := range apiRules {
		// 计算从左向右能匹配多少段
		// 从左向右逐段比较，遇到不匹配就停止
		matchedSegments := 0
		minLen := len(apiPath)
		if len(rule.Path) < minLen {
			minLen = len(rule.Path)
		}

		for i := range minLen {
			if apiPath[i] == rule.Path[i] {
				matchedSegments++
			} else {
				// 遇到不匹配就停止
				break
			}
		}

		// 只有当匹配段数等于规则路径长度时，该规则才有效（规则路径是请求路径的前缀）
		if matchedSegments == len(rule.Path) && len(rule.Path) > maxRuleLength {
			maxRuleLength = len(rule.Path)
			bestRule = &rule
		}
	}

	// 如果找到有效匹配的规则，返回该规则的权限设置
	if bestRule != nil {
		return bestRule.Rule
	}

	// 没有找到任何匹配的规则，默认拒绝访问
	return false
}

// ParseDuration 解析时间字符串为秒数
// 支持的格式：
// - 10d 或 10D：10天
// - 2h 或 2H：2小时
// - 5m 或 5M：5分钟
// - 100 或 100s 或 100S：100秒
// 转换失败时返回0
func ParseDuration(duration string) int64 {
	if duration == "" {
		return 0
	}

	// 如果是纯数字，直接按秒处理
	if value, err := strconv.Atoi(duration); err == nil {
		return int64(value)
	}

	// 获取最后一个字符作为单位
	unit := duration[len(duration)-1]
	// 如果最后一个字符是's'或'S'，去掉它再尝试转换
	if unit == 's' || unit == 'S' {
		if value, err := strconv.Atoi(duration[:len(duration)-1]); err == nil {
			return int64(value)
		}
		return 0
	}

	// 解析数值部分
	valueStr := duration[:len(duration)-1]
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0
	}

	// 根据单位转换为秒
	switch unit {
	case 'h', 'H': // 小时
		return int64(value * 3600)
	case 'm', 'M': // 分钟
		return int64(value * 60)
	case 'd', 'D': // 天
		return int64(value * 86400)
	default:
		return 0
	}
}
