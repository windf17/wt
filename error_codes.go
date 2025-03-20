package wtoken

type ErrorCode int

func (e ErrorCode) Error() string {
	return defaultRegistry.getErrorMessage(e)
}

func (e ErrorCode) String() string {
	return defaultRegistry.getErrorMessage(e)
}

func (e ErrorCode) Code() int {
	return int(e)
}

const E_Success ErrorCode = 0 // 成功
/******************** 基础错误 (1000-1999) ********************/
const (
	// 系统级错误
	E_Internal       ErrorCode = 1001 // 系统内部错误
	E_System         ErrorCode = 1002 // 系统错误
	E_ConfigInvalid  ErrorCode = 1101 // 配置无效
	E_CacheLoadFail  ErrorCode = 1201 // 缓存加载失败
	E_CacheParseFail ErrorCode = 1202 // 缓存解析失败
	E_InvalidParams  ErrorCode = 1301 // 无效参数
	E_APINotFound    ErrorCode = 1401 // 接口不存在
)

/******************** 认证鉴权 (2000-2999) ********************/
const (
	// 核心认证
	E_Unauthorized ErrorCode = 2001 // 未授权访问
	E_Forbidden    ErrorCode = 2002 // 禁止访问

	// Token管理
	E_InvalidToken  ErrorCode = 2101 // 无效令牌
	E_TokenExpired  ErrorCode = 2102 // 令牌过期
	E_TokenLimit    ErrorCode = 2103 // 令牌数量超限
	E_TokenGenerate ErrorCode = 2104 // 令牌生成失败

	// 安全验证
	E_InvalidAuthHeader ErrorCode = 2105 // 无效认证头
	E_CaptchaInvalid    ErrorCode = 2106 // 验证码错误

	// IP控制
	E_InvalidIP    ErrorCode = 2201 // 无效IP地址
	E_IPMismatch   ErrorCode = 2202 // IP不匹配
	E_IPNotAllowed ErrorCode = 2203 // IP未授权

	// 会话状态
	E_NotLogin   ErrorCode = 2301 // 未登录
	E_UserLogout ErrorCode = 2302 // 用户已登出
)

/******************** 用户管理 (3000-3999) ********************/
const (
	// 用户状态
	E_UserNotFound ErrorCode = 3001 // 用户不存在
	E_UserInvalid  ErrorCode = 3002 // 用户无效
	E_UserExists   ErrorCode = 3003 // 用户已存在
	E_UserDisabled ErrorCode = 3004 // 用户已禁用
	E_UserPending  ErrorCode = 3005 // 用户审核中

	// 用户组管理
	E_GroupNotFound ErrorCode = 3201 // 用户组不存在
	E_GroupInvalid  ErrorCode = 3202 // 用户组无效
	E_GroupExists   ErrorCode = 3203 // 用户组已存在
	E_GroupDisabled ErrorCode = 3204 // 用户组已禁用
	E_GroupPending  ErrorCode = 3205 // 用户组审核中
)

/******************** 密码安全 (4000-4999) ********************/
const (
	E_PasswordInvalid ErrorCode = 4001 // 密码错误
	E_PasswordExpired ErrorCode = 4002 // 密码已过期
	E_PasswordWeak    ErrorCode = 4003 // 密码强度不足
	E_PasswordReuse   ErrorCode = 4004 // 密码重复使用
	E_PasswordLocked  ErrorCode = 4005 // 密码已锁定
)

/******************** 数据库错误 (5000-5999) ********************/
const (
	// 基础操作
	E_DBConnect       ErrorCode = 5001 // 数据库连接失败
	E_DBQuery         ErrorCode = 5010 // 数据库查询失败
	E_DBQueryNotFound ErrorCode = 5011 // 记录不存在
	E_DBInsert        ErrorCode = 5020 // 数据库插入失败
	E_DBUpdate        ErrorCode = 5030 // 数据库更新失败
	E_DBDelete        ErrorCode = 5040 // 数据库删除失败

	// 高级特性
	E_DBTransaction ErrorCode = 5101 // 事务操作失败
	E_DBDeadlock    ErrorCode = 5102 // 数据库死锁
	E_DBTimeout     ErrorCode = 5103 // 操作超时

	// 数据完整性
	E_DBDuplicate  ErrorCode = 5201 // 唯一键冲突
	E_DBForeignKey ErrorCode = 5202 // 外键约束违反

	// 扩展功能
	E_DBBackupFail ErrorCode = 5301 // 数据库备份失败
)

/******************** 其他通用 (9000-9999) ********************/
const (
	E_Unknown ErrorCode = 9999 // 未知错误
)
