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
