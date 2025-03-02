package wtoken

// IManager 定义Token管理器接口
type IManager[T any] interface {
	// GetToken 获取token
	GetToken(key string) (*Token[T], ErrorData)
	// AddToken 新增Token
	AddToken(userID uint, groupID uint, ip string) (string, ErrorData)
	// GenerateToken 生成token
	GenerateToken() (string, error)
	// DelToken 删除token
	DelToken(key string) ErrorData
	// DelTokenByUserId 根据用户ID删除所有对应的token
	DelTokensByUserID(userID uint) ErrorData
	// DelTokensByGroupID 根据用户组ID删除所有对应的token
	DelTokensByGroupID(groupID uint) ErrorData
	// UpdateToken 更新token
	UpdateToken(key string, token *Token[T]) ErrorData
	// CheckToken 检查token是否有效
	CheckToken(key string) ErrorData
	// CleanExpiredTokens 清理过期的token
	CleanExpiredTokens()

	// Authenticate 鉴权
	Authenticate(key string, url string, ip string) ErrorData

	// GetGroup 获取指定用户组信息
	GetGroup(groupID uint) (*Group, ErrorData)
	// AddGroup 新增用户组信息
	AddGroup(raw GroupRaw) ErrorData
	// DeleteGroup 删除用户组信息
	DelGroup(groupID uint) ErrorData
	// UpdateGroup 更新用户组信息
	UpdateGroup(raw GroupRaw) ErrorData

	// GetStats 获取统计信息
	GetStats() Stats

	// 新增错误信息
	NewError(code ErrorCode) ErrorData

	// 存储用户数据
	SaveData(key string, data T) error
	// 获取用户数据
	GetData(key string) (T, error)
}
