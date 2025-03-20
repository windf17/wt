package wtoken

// Language 定义语言类型
type Language string

const (
	// LangChinese 中文
	LangChinese Language = "zh"
	// LangEnglish 英文
	LangEnglish Language = "en"
)

// 多语言错误信息配置
var ErrorI18nMessages = map[Language]map[ErrorCode]string{
	LangChinese: {
		// 基础错误
		E_Internal:       "系统内部错误",
		E_System:         "系统错误",
		E_ConfigInvalid:  "配置无效",
		E_CacheLoadFail:  "缓存加载失败",
		E_CacheParseFail: "缓存解析失败",
		E_InvalidParams:  "无效参数",
		E_APINotFound:    "接口不存在",

		// 认证鉴权
		E_Unauthorized:      "未授权访问",
		E_Forbidden:         "禁止访问",
		E_InvalidToken:      "无效令牌",
		E_TokenExpired:      "令牌过期",
		E_TokenLimit:        "令牌数量超限",
		E_TokenGenerate:     "令牌生成失败",
		E_InvalidAuthHeader: "请求头认证格式错误",
		E_CaptchaInvalid:    "验证码错误",
		E_InvalidIP:         "无效IP地址",
		E_IPMismatch:        "IP地址不匹配",
		E_IPNotAllowed:      "IP未授权",
		E_NotLogin:          "用户未登录",
		E_UserLogout:        "用户已登出",

		// 用户管理
		E_UserNotFound:  "用户不存在",
		E_UserInvalid:   "用户无效",
		E_UserExists:    "用户已存在",
		E_UserDisabled:  "用户已禁用",
		E_UserPending:   "用户审核中",
		E_GroupNotFound: "用户组不存在",
		E_GroupInvalid:  "用户组无效",
		E_GroupExists:   "用户组已存在",
		E_GroupDisabled: "用户组已禁用",
		E_GroupPending:  "用户组审核中",

		// 密码安全
		E_PasswordInvalid: "密码错误",
		E_PasswordExpired: "密码已过期",
		E_PasswordWeak:    "密码强度不足",
		E_PasswordReuse:   "密码重复使用",
		E_PasswordLocked:  "密码已锁定",

		// 数据库
		E_DBConnect:       "数据库连接失败",
		E_DBQuery:         "数据库查询失败",
		E_DBQueryNotFound: "记录不存在",
		E_DBInsert:        "数据库插入失败",
		E_DBUpdate:        "数据库更新失败",
		E_DBDelete:        "数据库删除失败",
		E_DBTransaction:   "事务操作失败",
		E_DBDeadlock:      "数据库死锁",
		E_DBTimeout:       "数据库操作超时",
		E_DBDuplicate:     "唯一键冲突",
		E_DBForeignKey:    "外键约束违反",
		E_DBBackupFail:    "数据库备份失败",

		E_Unknown: "未知错误",
	},
	LangEnglish: {
		E_Internal:       "Internal system error",
		E_System:         "System error",
		E_ConfigInvalid:  "Invalid configuration",
		E_CacheLoadFail:  "Cache load failed",
		E_CacheParseFail: "Cache parse failed",
		E_InvalidParams:  "Invalid parameters",
		E_APINotFound:    "API not found",

		E_Unauthorized:      "Unauthorized access",
		E_Forbidden:         "Access forbidden",
		E_InvalidToken:      "Invalid token",
		E_TokenExpired:      "Token expired",
		E_TokenLimit:        "Token limit exceeded",
		E_TokenGenerate:     "Token generation failed",
		E_InvalidAuthHeader: "Invalid authorization header",
		E_CaptchaInvalid:    "Invalid captcha",
		E_InvalidIP:         "Invalid IP address",
		E_IPMismatch:        "IP address mismatch",
		E_IPNotAllowed:      "IP not allowed",
		E_NotLogin:          "User not logged in",
		E_UserLogout:        "User logged out",

		E_UserNotFound:  "User not found",
		E_UserInvalid:   "Invalid user",
		E_UserExists:    "User already exists",
		E_UserDisabled:  "User disabled",
		E_UserPending:   "User verification pending",
		E_GroupNotFound: "User group not found",
		E_GroupInvalid:  "Invalid user group",
		E_GroupExists:   "User group already exists",
		E_GroupDisabled: "User group disabled",
		E_GroupPending:  "User group verification pending",

		E_PasswordInvalid: "Invalid password",
		E_PasswordExpired: "Password expired",
		E_PasswordWeak:    "Password strength insufficient",
		E_PasswordReuse:   "Password reused",
		E_PasswordLocked:  "Password locked",

		E_DBConnect:       "Database connection failed",
		E_DBQuery:         "Database query failed",
		E_DBInsert:        "Database insert failed",
		E_DBQueryNotFound: "Record not found",
		E_DBUpdate:        "Database update failed",
		E_DBDelete:        "Database delete failed",
		E_DBTransaction:   "Transaction failed",
		E_DBDeadlock:      "Database deadlock",
		E_DBTimeout:       "Database operation timeout",
		E_DBDuplicate:     "Duplicate key violation",
		E_DBForeignKey:    "Foreign key violation",
		E_DBBackupFail:    "Database backup failed",

		E_Unknown: "Unknown error",
	},
}
