package models

// IManager token管理器接口
type IManager[T any] interface {
	// token管理
	GetToken(key string) (*Token[T], error)
	AddToken(userID uint, groupID uint, clientIp string) (string, error)
	GenerateToken() (string, error)
	DelToken(key string) error
	DelTokensByUserID(userID uint) error
	DelTokensByGroupID(groupID uint) error
	UpdateToken(key string, token *Token[T]) error
	CleanExpiredTokens()

	// 批量操作
	BatchDeleteTokensByUserIDs(userIDs []uint) error
	BatchDeleteTokensByGroupIDs(groupIDs []uint) error
	BatchDeleteExpiredTokens() error
	GetTokensByUserID(userID uint) []*Token[T]
	GetTokensByGroupID(groupID uint) []*Token[T]

	// 身份验证
	Auth(key string, clientIp string, api string) error
	BatchAuth(key string, clientIp string, apis []string) []bool

	// 用户组管理
	GetGroup(groupID uint) (*Group, error)
	AddGroup(group *GroupRaw) error
	DelGroup(groupID uint) error
	UpdateGroup(groupID uint, group *GroupRaw) error
	UpdateAllGroup(groups []GroupRaw) error

	// 统计信息
	GetStats() Stats

	// 用户数据管理
	SetUserData(key string, data T) error
	GetUserData(key string) (T, error)
}