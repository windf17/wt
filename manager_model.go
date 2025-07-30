package wtoken

import "github.com/windf17/wtoken/models"

// IManager token管理器接口
type IManager[T any] interface {
	// token管理
	GetToken(key string) (*Token[T], ErrorCode)
	AddToken(userID uint, groupID uint, clientIp string) (string, ErrorCode)
	GenerateToken() (string, error)
	DelToken(key string) ErrorCode
	DelTokensByUserID(userID uint) ErrorCode
	DelTokensByGroupID(groupID uint) ErrorCode
	UpdateToken(key string, token *Token[T]) ErrorCode
	CleanExpiredTokens()

	// 批量操作
	BatchDeleteTokensByUserIDs(userIDs []uint) ErrorCode
	BatchDeleteTokensByGroupIDs(groupIDs []uint) ErrorCode
	BatchDeleteExpiredTokens() ErrorCode
	GetTokensByUserID(userID uint) []*Token[T]
	GetTokensByGroupID(groupID uint) []*Token[T]

	// 身份验证
	Auth(key string, clientIp string, api string) ErrorCode
	BatchAuth(key string, clientIp string, apis []string) []bool

	// 用户组管理
	GetGroup(groupID uint) (*models.Group, ErrorCode)
	AddGroup(group *models.GroupRaw) ErrorCode
	DelGroup(groupID uint) ErrorCode
	UpdateGroup(groupID uint, group *models.GroupRaw) ErrorCode
	UpdateAllGroup(groups []models.GroupRaw) ErrorCode

	// 统计信息
	GetStats() Stats

	// 用户数据存储
	SetUserData(key string, data T) ErrorCode
	GetUserData(key string) (T, ErrorCode)

	// 资源管理
	Close()
}
