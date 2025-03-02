package test

import (
	"sort"
	"strings"
	"testing"

	"github.com/windf17/wtoken/utility"
)

// TestHasPermission 测试 HasPermission 函数
func TestHasPermission(t *testing.T) {
	type hasPermissionTest struct {
		name        string
		apiPath     string
		allowedAPIs []string
		deniedAPIs  []string
		expected    bool
	}

	tests := []hasPermissionTest{

		// 案例1：完全匹配 DeniedAPIs
		{
			name:        "完全匹配 DeniedAPIs",
			apiPath:     "/api/v1/product/delete/",
			allowedAPIs: []string{},
			deniedAPIs:  []string{"/api/v1/product/delete/"},
			expected:    false,
		},

		// 案例2：完全匹配 AllowedAPIs
		{
			name:        "完全匹配 AllowedAPIs",
			apiPath:     "/api/v1/product/list/",
			allowedAPIs: []string{"/api/v1/product/list/"},
			deniedAPIs:  []string{"/api/v1/product/delete/"},
			expected:    true,
		},

		// 案例3：分段匹配 DeniedAPIs
		{
			name:        "分段匹配 DeniedAPIs",
			apiPath:     "/api/v1/product/status",
			allowedAPIs: []string{},
			deniedAPIs:  []string{"/api/v1/product/"},
			expected:    false,
		},

		// 案例4：分段匹配 AllowedAPIs
		{
			name:        "分段匹配 AllowedAPIs",
			apiPath:     "/api/v1/product/status",
			allowedAPIs: []string{"/api/v1/product/"},
			deniedAPIs:  []string{},
			expected:    true,
		},

		// 案例5：DeniedAPIs 优先于 AllowedAPIs
		{
			name:        "DeniedAPIs 优先于 AllowedAPIs",
			apiPath:     "/api/v1/user/info/",
			allowedAPIs: []string{"/api/v1/user/"},
			deniedAPIs:  []string{"/api/v1/user/info/"},
			expected:    false,
		},

		// 案例6：多个 DeniedAPIs 和 AllowedAPIs
		{
			name:        "多个 DeniedAPIs 和 AllowedAPIs",
			apiPath:     "/api/v1/product/list/0/",
			allowedAPIs: []string{"/api/v1/product/list/"},
			deniedAPIs:  []string{"/api/v1/product/delete/"},
			expected:    true,
		},

		// 案例7：路径标准化（以 / 结尾）
		{
			name:        "路径标准化（以 / 结尾）",
			apiPath:     "/api/v1/product/list",
			allowedAPIs: []string{"/api/v1/product/list/"},
			deniedAPIs:  []string{},
			expected:    true,
		},

		// 案例8：路径标准化（没有 / 结尾）
		{
			name:        "路径标准化（没有 / 结尾）",
			apiPath:     "/api/v1/product/list/",
			allowedAPIs: []string{"/api/v1/product/list"},
			deniedAPIs:  []string{},
			expected:    true,
		},

		// 案例9：无匹配路径
		{
			name:        "无匹配路径",
			apiPath:     "/api/v2/product/status",
			allowedAPIs: []string{"/api/v1/product/list/"},
			deniedAPIs:  []string{"/api/v1/product/delete/"},
			expected:    false,
		},

		// 案例10：最长前缀匹配 DeniedAPIs
		{
			name:        "最长前缀匹配 DeniedAPIs",
			apiPath:     "/api/v3/product/status/detail",
			deniedAPIs:  []string{"/api/v3/product/status/"},
			allowedAPIs: []string{"/api/v3/product/"},
			expected:    false,
		},

		// 案例11：最长前缀匹配 AllowedAPIs
		{
			name:        "最长前缀匹配 AllowedAPIs",
			apiPath:     "/api/v3/product/status/detail",
			allowedAPIs: []string{"/api/v3/product/status/"},
			deniedAPIs:  []string{"/api/v3/product/"},
			expected:    true,
		},

		// 案例12：AllowedAPIs 为空
		{
			name:        "AllowedAPIs 为空",
			apiPath:     "/api/v1/product/list/",
			allowedAPIs: []string{},
			deniedAPIs:  []string{},
			expected:    false,
		},

		// 案例13：DeniedAPIs 为空
		{
			name:        "DeniedAPIs 为空",
			apiPath:     "/api/v1/product/list/",
			allowedAPIs: []string{"/api/v1/product/list/"},
			deniedAPIs:  []string{},
			expected:    true,
		},

		// 案例14：Allowed 和 Denied 长度相同
		{
			name:        "Allowed 和 Denied 长度相同",
			apiPath:     "/api/v1/product/info/",
			allowedAPIs: []string{"/api/v1/product/"},
			deniedAPIs:  []string{"/api/v1/product/"},
			expected:    false,
		},

		// 案例15：Denied 包含更短的路径
		{
			name:        "Denied 包含更短的路径",
			apiPath:     "/api/v1/product/list/0/",
			allowedAPIs: []string{"/api/v1/product/list/"},
			deniedAPIs:  []string{"/api/v1/product/"},
			expected:    true,
		},

		// 案例16：Allowed 包含更长的路径
		{
			name:        "Allowed 包含更长的路径",
			apiPath:     "/api/v1/user/info/0/",
			allowedAPIs: []string{"/api/v1/user/info/"},
			deniedAPIs:  []string{"/api/v1/user/"},
			expected:    true,
		},

		// 案例17：Allowed 和 Denied 都包含更长的路径
		{
			name:        "Allowed 和 Denied 都包含更长的路径",
			apiPath:     "/api/v1/product/list/0/",
			allowedAPIs: []string{"/api/v1/product/list/"},
			deniedAPIs:  []string{"/api/v1/product/list/0/"},
			expected:    false,
		},

		// 案例18：Allowed 包含更长的路径， Denied 包含更短的路径
		{
			name:        "Allowed 包含更长的路径， Denied 包含更短的路径",
			apiPath:     "/api/v1/product/list/0/",
			allowedAPIs: []string{"/api/v1/product/list/"},
			deniedAPIs:  []string{"/api/v1/product/"},
			expected:    true,
		},

		// 案例19：Allowed 包含完整匹配
		{
			name:        "Allowed 包含完整匹配",
			apiPath:     "/api/v1/user/info/",
			allowedAPIs: []string{"/api/v1/user/info/"},
			deniedAPIs:  []string{"/api/v1/user/"},
			expected:    true,
		},

		// 案例20：Denied 包含完整匹配
		{
			name:        "Denied 包含完整匹配",
			apiPath:     "/api/v1/user/info/",
			allowedAPIs: []string{"/api/v1/user/"},
			deniedAPIs:  []string{"/api/v1/user/info/"},
			expected:    false,
		},

		// 案例21：分段匹配 DeniedAPIs 中的不同位置
		{
			name:        "分段匹配 DeniedAPIs 中的不同位置",
			apiPath:     "/api/v1/product/delete/0/",
			deniedAPIs:  []string{"/api/v1/product/delete/"},
			allowedAPIs: []string{},
			expected:    false,
		},

		// 案例22：分段匹配 AllowedAPIs 中的不同位置
		{
			name:        "分段匹配 AllowedAPIs 中的不同位置",
			apiPath:     "/api/v1/product/list/0/",
			allowedAPIs: []string{"/api/v1/product/list/"},
			deniedAPIs:  []string{},
			expected:    true,
		},

		// 案例23：Allowed 和 Denied 都包含更长的路径
		{
			name:        "Allowed 和 Denied 都包含更长的路径",
			apiPath:     "/api/v1/user/info/0/",
			allowedAPIs: []string{"/api/v1/user/info/"},
			deniedAPIs:  []string{"/api/v1/user/info/0/"},
			expected:    false,
		},

		// 案例24：Allowed 包含更长的路径，Denied 包含更短的路径
		{
			name:        "Allowed 包含更长的路径，Denied 包含更短的路径",
			apiPath:     "/api/v1/user/info/0/",
			allowedAPIs: []string{"/api/v1/user/info/"},
			deniedAPIs:  []string{"/api/v1/user/"},
			expected:    true,
		},

		// 案例25：Denied 包含更长的路径，Allowed 包含更短的路径
		{
			name:        "Denied 包含更长的路径，Allowed 包含更短的路径",
			apiPath:     "/api/v1/user/info/0/",
			deniedAPIs:  []string{"/api/v1/user/info/0/"},
			allowedAPIs: []string{"/api/v1/user/info/"},
			expected:    false,
		},

		// 案例26：Allowed 包含更长的路径，Denied 包含更短的路径（路径标准化测试）
		{
			name:        "Allowed 包含更长的路径，Denied 包含更短的路径（路径标准化测试）",
			apiPath:     "/api/v1/user/info",
			allowedAPIs: []string{"/api/v1/user/info/"},
			deniedAPIs:  []string{"/api/v1/user/"},
			expected:    true,
		},

		// 案例27：Denied 包含更长的路径，Allowed 包含更短的路径（路径标准化测试）
		{
			name:        "Denied 包含更长的路径，Allowed 包含更短的路径（路径标准化测试）",
			apiPath:     "/api/v1/user/info/0",
			deniedAPIs:  []string{"/api/v1/user/info/0/"},
			allowedAPIs: []string{"/api/v1/user/info/"},
			expected:    false,
		},

		// 案例28：测试空 APIpath
		{
			name:        "测试空 APIpath",
			apiPath:     "",
			allowedAPIs: []string{"/api/v1/product/list/"},
			deniedAPIs:  []string{},
			expected:    false,
		},

		// 案例29：测试 APIpath 不以 / 开头
		{
			name:        "测试 APIpath 不以 / 开头",
			apiPath:     "api/v1/product/list/",
			allowedAPIs: []string{"/api/v1/product/list/"},
			deniedAPIs:  []string{},
			expected:    true,
		},

		// 案例30：测试 APIpath 中包含特殊字符，不处理特殊字符，直接匹配
		// {
		// 	name:        "测试 APIpath 中包含特殊字符",
		// 	apiPath:     "/api/v1/product%20list/",
		// 	allowedAPIs: []string{"/api/v1/product%20list/"},
		// 	deniedAPIs:  []string{},
		// 	expected:    true,
		// },
	}

	// 预处理 allowedAPIs 和 deniedAPIs
	for i, test := range tests {
		apiPath := utility.NormalizeAPIPath(test.apiPath)
		allowedStr := strings.Join(test.allowedAPIs, " ")
		deniedStr := strings.Join(test.deniedAPIs, " ")

		// 将 allowedAPIs 和 deniedAPIs 转换为切片
		parsedAllowed := utility.ParseAPIs(allowedStr, " ")
		parsedDenied := utility.ParseAPIs(deniedStr, " ")

		// 确保 allowedAPIs 和 deniedAPIs 按路径长度排序
		if len(parsedAllowed) > 0 {
			sort.Slice(parsedAllowed, func(i, j int) bool {
				return len(parsedAllowed[i]) > len(parsedAllowed[j])
			})
		}
		if len(parsedDenied) > 0 {
			sort.Slice(parsedDenied, func(i, j int) bool {
				return len(parsedDenied[i]) > len(parsedDenied[j])
			})
		}

		tests[i].apiPath = apiPath
		tests[i].allowedAPIs = parsedAllowed
		tests[i].deniedAPIs = parsedDenied
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utility.HasPermission(tt.apiPath, tt.allowedAPIs, tt.deniedAPIs)
			if result != tt.expected {
				t.Errorf("HasPermission(%s, %v, %v) = %v; want %v", tt.apiPath, tt.allowedAPIs, tt.deniedAPIs, result, tt.expected)
			}
		})
	}
}
