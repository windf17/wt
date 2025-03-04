package utility

import (
	"net/url"
	"strconv"
	"strings"
)

// NormalizeAPIPath：处理url，返回api路径
func NormalizeAPIPath(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	path := parsedURL.Path
	return parseApiString(path)
}

// ParseApiString：处理api路径，如果是空字符串返回空字符串，非空字符串则确保首位都有/
func parseApiString(api string) string {
	api = strings.Trim(api, "/")
	if api == "" {
		return ""
	}
	api = "/" + api + "/"
	return api
}

// HasPermission 检查api权限
func HasPermission(apiPath string, allowedAPIs, deniedAPIs []string) bool {
	if apiPath == "" {
		return false
	}

	apiSegments := splitPathSegments(apiPath)

	maxDeniedSegments := 0
	maxAllowedSegments := 0

	// Check deniedAPIs
	for _, deniedAPI := range deniedAPIs {
		ruleSegments := splitPathSegments(deniedAPI)
		if pathSegmentsMatch(apiSegments, ruleSegments) {
			if len(ruleSegments) == len(apiSegments) {
				return false
			}
			if len(ruleSegments) > maxDeniedSegments {
				maxDeniedSegments = len(ruleSegments)
			}
		}
	}

	// Check allowedAPIs
	for _, allowedAPI := range allowedAPIs {
		ruleSegments := splitPathSegments(allowedAPI)
		if pathSegmentsMatch(apiSegments, ruleSegments) {
			if len(ruleSegments) == len(apiSegments) {
				return true
			}
			if len(ruleSegments) > maxAllowedSegments {
				maxAllowedSegments = len(ruleSegments)
			}
		}
	}

	return maxAllowedSegments > maxDeniedSegments
}

// splitPathSegments splits the path into segments
func splitPathSegments(path string) []string {
	if path == "" {
		return nil
	}
	segments := strings.Split(path, "/")
	if len(segments) <= 1 {
		return nil
	}
	start := 1
	end := len(segments) - 1
	return segments[start:end]
}

// pathSegmentsMatch checks if the path segments match the rule segments
func pathSegmentsMatch(apiSegments, ruleSegments []string) bool {
	if len(ruleSegments) > len(apiSegments) {
		return false
	}
	for i := range ruleSegments {
		if apiSegments[i] != ruleSegments[i] {
			return false
		}
	}
	return true
}

// ParseAPIs：将api权限字符串转换成api字符串数组
func ParseAPIs(apiStr string, delimiter string) []string {
	if apiStr == "" {
		return []string{}
	}

	apis := strings.Split(apiStr, delimiter)
	var newApis []string
	for _, api := range apis {
		a := parseApiString(api)
		if a != "" {
			newApis = append(newApis, a)
		}
	}

	return newApis
}


// ParseDuration 解析时间字符串为秒数
// 支持的格式：
// - 10d 或 10D：10天
// - 5m 或 5M：5分钟
// - 2h 或 2H：2小时
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