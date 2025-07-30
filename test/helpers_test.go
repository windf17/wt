package test

import (
	"testing"

	"github.com/windf17/wt"
	"github.com/windf17/wt/utility"
)

// ========== Utility Functions Tests ==========

/**
 * TestParseURLToPathSegments 测试URL解析为路径段功能
 */
func TestParseURLToPathSegments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "完整URL",
			input:    "https://example.com/api/v1/users",
			expected: []string{"api", "v1", "users"},
		},
		{
			name:     "HTTP URL",
			input:    "http://localhost:8080/auth/login",
			expected: []string{"auth", "login"},
		},
		{
			name:     "只有路径的URL",
			input:    "/api/v2/tokens",
			expected: []string{"api", "v2", "tokens"},
		},
		{
			name:     "根路径",
			input:    "https://example.com/",
			expected: []string{},
		},
		{
			name:     "空字符串",
			input:    "",
			expected: []string{},
		},
		{
			name:     "无效URL",
			input:    "://invalid-url",
			expected: []string{},
		},
		{
			name:     "带查询参数的URL",
			input:    "https://api.example.com/users/123?name=test&age=25",
			expected: []string{"users", "123"},
		},
		{
			name:     "带锚点的URL",
			input:    "https://example.com/docs/guide#section1",
			expected: []string{"docs", "guide"},
		},
		{
			name:     "多级嵌套路径",
			input:    "https://api.example.com/v1/users/123/posts/456/comments",
			expected: []string{"v1", "users", "123", "posts", "456", "comments"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utility.ParseURLToPathSegments(tt.input)
			if !equalStringSlices(result, tt.expected) {
				t.Errorf("ParseURLToPathSegments(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

/**
 * TestParsePathToSegments 测试路径字符串解析为路径段功能
 */
func TestParsePathToSegments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "标准路径",
			input:    "/api/v1/users",
			expected: []string{"api", "v1", "users"},
		},
		{
			name:     "无前导斜杠的路径",
			input:    "api/v1/users",
			expected: []string{"api", "v1", "users"},
		},
		{
			name:     "带尾随斜杠的路径",
			input:    "/api/v1/users/",
			expected: []string{"api", "v1", "users"},
		},
		{
			name:     "根路径",
			input:    "/",
			expected: []string{},
		},
		{
			name:     "空字符串",
			input:    "",
			expected: []string{},
		},
		{
			name:     "连续斜杠",
			input:    "/api//v1///users",
			expected: []string{"api", "v1", "users"},
		},
		{
			name:     "带空白字符的路径段",
			input:    "/api/ v1 /users",
			expected: []string{"api", "v1", "users"},
		},
		{
			name:     "单个路径段",
			input:    "/api",
			expected: []string{"api"},
		},
		{
			name:     "只有空白字符的路径段",
			input:    "/api/   /users",
			expected: []string{"api", "users"},
		},
		{
			name:     "复杂路径",
			input:    "///api//v1/users/123//posts///",
			expected: []string{"api", "v1", "users", "123", "posts"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utility.ParsePathToSegments(tt.input)
			if !equalStringSlices(result, tt.expected) {
				t.Errorf("ParsePathToSegments(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// ========== Validation Functions Tests ==========

// TestValidateIPAddress 测试IP地址验证
func TestValidateIPAddress(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{"Valid IPv4", "192.168.1.1", true},
		{"Valid IPv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", true},
		{"Valid Localhost", "127.0.0.1", true},
		{"Invalid Empty IP", "", false},
		{"Invalid IP Format", "256.256.256.256", false},
		{"Invalid IP String", "not.an.ip.address", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wt.ValidateIPAddress(tt.ip)
			if (err == nil) != tt.expected {
				t.Errorf("ValidateIPAddress(%s) = %v, expected %v", tt.ip, err == nil, tt.expected)
			}
		})
	}
}

// TestValidateTokenExpire 测试Token过期时间验证
func TestValidateTokenExpire(t *testing.T) {
	tests := []struct {
		name     string
		expire   int
		expected bool
	}{
		{"Valid Expire Time", 3600, true},
		{"Valid Min Expire Time", 60, true},
		{"Valid Max Expire Time", 86400, true},
		{"Invalid Zero Expire", 0, false},
		{"Invalid Negative Expire", -1, false},
		{"Invalid Too Small Expire", 30, false},
		{"Invalid Too Large Expire", 90000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wt.ValidateTokenExpire(tt.expire, nil)
			if (err == nil) != tt.expected {
				t.Errorf("ValidateTokenExpire(%d) = %v, expected %v", tt.expire, err == nil, tt.expected)
			}
		})
	}
}

// TestValidateStringNotEmpty 测试字符串非空验证
func TestValidateStringNotEmpty(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		expected  bool
	}{
		{"Valid String", "test", "field", true},
		{"Valid String with Spaces", "  test  ", "field", true},
		{"Invalid Empty String", "", "field", false},
		{"Invalid Only Spaces", "   ", "field", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wt.ValidateStringNotEmpty(tt.value, tt.fieldName)
			if (err == nil) != tt.expected {
				t.Errorf("ValidateStringNotEmpty(%s, %s) = %v, expected %v", tt.value, tt.fieldName, err == nil, tt.expected)
			}
		})
	}
}

// TestValidateStringLength 测试字符串长度验证
func TestValidateStringLength(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		minLen    int
		maxLen    int
		expected  bool
	}{
		{"Valid Length", "test", "field", 3, 10, true},
		{"Valid Min Length", "abc", "field", 3, 10, true},
		{"Valid Max Length", "abcdefghij", "field", 3, 10, true},
		{"Invalid Too Short", "ab", "field", 3, 10, false},
		{"Invalid Too Long", "abcdefghijk", "field", 3, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wt.ValidateStringLength(tt.value, tt.fieldName, tt.minLen, tt.maxLen)
			if (err == nil) != tt.expected {
				t.Errorf("ValidateStringLength(%s, %s, %d, %d) = %v, expected %v",
					tt.value, tt.fieldName, tt.minLen, tt.maxLen, err == nil, tt.expected)
			}
		})
	}
}

// TestValidatePositiveInt 测试正整数验证
func TestValidatePositiveInt(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		fieldName string
		expected  bool
	}{
		{"Valid Positive Int", 1, "field", true},
		{"Valid Large Positive Int", 999999, "field", true},
		{"Invalid Zero", 0, "field", false},
		{"Invalid Negative", -1, "field", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wt.ValidatePositiveInt(tt.value, tt.fieldName)
			if (err == nil) != tt.expected {
				t.Errorf("ValidatePositiveInt(%d, %s) = %v, expected %v", tt.value, tt.fieldName, err == nil, tt.expected)
			}
		})
	}
}

// TestValidateIntRange 测试整数范围验证
func TestValidateIntRange(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		fieldName string
		min       int
		max       int
		expected  bool
	}{
		{"Valid Range", 5, "field", 1, 10, true},
		{"Valid Min Value", 1, "field", 1, 10, true},
		{"Valid Max Value", 10, "field", 1, 10, true},
		{"Invalid Below Min", 0, "field", 1, 10, false},
		{"Invalid Above Max", 11, "field", 1, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wt.ValidateIntRange(tt.value, tt.fieldName, tt.min, tt.max)
			if (err == nil) != tt.expected {
				t.Errorf("ValidateIntRange(%d, %s, %d, %d) = %v, expected %v",
					tt.value, tt.fieldName, tt.min, tt.max, err == nil, tt.expected)
			}
		})
	}
}

// TestValidateSliceNotEmpty 测试切片非空验证
func TestValidateSliceNotEmpty(t *testing.T) {
	tests := []struct {
		name      string
		slice     []string
		fieldName string
		expected  bool
	}{
		{"Valid Non-Empty Slice", []string{"a", "b"}, "field", true},
		{"Valid Single Element", []string{"a"}, "field", true},
		{"Invalid Empty Slice", []string{}, "field", false},
		{"Invalid Nil Slice", nil, "field", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert []string to []any
			var interfaceSlice []any
			for _, v := range tt.slice {
				interfaceSlice = append(interfaceSlice, v)
			}
			err := wt.ValidateSliceNotEmpty(interfaceSlice, tt.fieldName)
			if (err == nil) != tt.expected {
				t.Errorf("ValidateSliceNotEmpty(%v, %s) = %v, expected %v", tt.slice, tt.fieldName, err == nil, tt.expected)
			}
		})
	}
}

// TestValidateMapNotEmpty 测试映射非空验证
func TestValidateMapNotEmpty(t *testing.T) {
	tests := []struct {
		name      string
		mapValue  map[string]string
		fieldName string
		expected  bool
	}{
		{"Valid Non-Empty Map", map[string]string{"a": "b"}, "field", true},
		{"Valid Multiple Elements", map[string]string{"a": "b", "c": "d"}, "field", true},
		{"Invalid Empty Map", map[string]string{}, "field", false},
		{"Invalid Nil Map", nil, "field", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert map[string]string to map[string]any
			interfaceMap := make(map[string]any)
			for k, v := range tt.mapValue {
				interfaceMap[k] = v
			}
			err := wt.ValidateMapNotEmpty(interfaceMap, tt.fieldName)
			if (err == nil) != tt.expected {
				t.Errorf("ValidateMapNotEmpty(%v, %s) = %v, expected %v", tt.mapValue, tt.fieldName, err == nil, tt.expected)
			}
		})
	}
}

// ========== Benchmark Tests ==========

/**
 * BenchmarkParseURLToPathSegments 性能基准测试
 */
func BenchmarkParseURLToPathSegments(b *testing.B) {
	testURL := "https://api.example.com/v1/users/123/posts/456/comments"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utility.ParseURLToPathSegments(testURL)
	}
}

/**
 * BenchmarkParsePathToSegments 性能基准测试
 */
func BenchmarkParsePathToSegments(b *testing.B) {
	testPath := "/v1/users/123/posts/456/comments"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utility.ParsePathToSegments(testPath)
	}
}

// BenchmarkValidateIPAddress 基准测试IP地址验证
func BenchmarkValidateIPAddress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wt.ValidateIPAddress("192.168.1.1")
	}
}
