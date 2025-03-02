package utility

import (
	"net/url"
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
