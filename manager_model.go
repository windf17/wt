package wtoken

// IManager 定义Token管理器接口
type IManager[T any] interface {
	// GetToken 获取token
	GetToken(key string) (*Token[T], ErrorCode)
	// AddToken 新增Token
	AddToken(userID uint, groupID uint, ip string) (string, ErrorCode)
	// GenerateToken 生成token
	GenerateToken() (string, error)
	// DelToken 删除token
	DelToken(key string) ErrorCode
	// DelTokenByUserId 根据用户ID删除所有对应的token
	DelTokensByUserID(userID uint) ErrorCode
	// DelTokensByGroupID 根据用户组ID删除所有对应的token
	DelTokensByGroupID(groupID uint) ErrorCode
	// UpdateToken 更新token
	UpdateToken(key string, token *Token[T]) ErrorCode
	// CheckToken 检查token是否有效
	CheckToken(key string) ErrorCode
	// CleanExpiredTokens 清理过期的token
	CleanExpiredTokens()

	// Authenticate 鉴权
	Authenticate(key string, url string, ip string) ErrorCode

	// GetGroup 获取指定用户组信息
	GetGroup(groupID uint) (*Group, ErrorCode)
	// AddGroup 新增用户组信息
	AddGroup(raw GroupRaw) ErrorCode
	// DeleteGroup 删除用户组信息
	DelGroup(groupID uint) ErrorCode
	// UpdateGroup 更新用户组信息
	UpdateGroup(raw GroupRaw) ErrorCode

	// GetStats 获取统计信息
	GetStats() Stats

	// 存储用户数据
	SaveData(key string, data T) ErrorCode
	// 获取用户数据
	GetData(key string) (T, ErrorCode)
}
