package wt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"golang.org/x/crypto/pbkdf2"
)

/**
 * SecurityManager 安全管理器
 */
type SecurityManager struct {
	key  []byte // 加密密钥
	salt []byte // 盐值
}

/**
 * NewSecurityManager 创建一个新的安全管理器
 * @param {string} password 密码
 * @returns {*SecurityManager} 安全管理器实例
 */
func NewSecurityManager(password string) *SecurityManager {
	return NewSecurityManagerWithSalt(password, nil)
}

/*
NewSecurityManagerWithSalt 创建一个带盐值的安全管理器
*/
func NewSecurityManagerWithSalt(password string, salt []byte) *SecurityManager {
	if len(salt) == 0 {
		salt = make([]byte, 32)
		rand.Read(salt)
	}
	
	// 使用PBKDF2增强密钥安全性
	key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)
	
	return &SecurityManager{
		key:  key,
		salt: salt,
	}
}

/*
RotateKey 轮换密钥
*/
func (sm *SecurityManager) RotateKey(newPassword string) {
	// 生成新的盐值
	newSalt := make([]byte, 32)
	rand.Read(newSalt)
	
	// 生成新的密钥
	newKey := pbkdf2.Key([]byte(newPassword), newSalt, 10000, 32, sha256.New)
	
	// 原子性更新密钥和盐值
	sm.key = newKey
	sm.salt = newSalt
}

/**
 * EncryptToken 加密token数据
 * @param {string} plaintext 明文token
 * @returns {string, error} 加密后的token和错误
 */
func (sm *SecurityManager) EncryptToken(plaintext string) (string, error) {
	block, err := aes.NewCipher(sm.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

/**
 * DecryptToken 解密token数据
 * @param {string} ciphertext 加密的token
 * @returns {string, error} 解密后的token和错误
 */
func (sm *SecurityManager) DecryptToken(ciphertext string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(sm.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

/**
 * HashSensitiveData 对敏感数据进行哈希处理
 * @param {string} data 敏感数据
 * @returns {string} 哈希值
 */
func (sm *SecurityManager) HashSensitiveData(data string) string {
	hash := sha256.Sum256([]byte(data + string(sm.key)))
	return base64.URLEncoding.EncodeToString(hash[:])
}


/**
 * ValidateTokenFormat 验证token格式
 * @param {string} token token字符串
 * @returns {bool} 是否为有效格式
 */
func ValidateTokenFormat(token string) bool {
	if len(token) == 0 || len(token) > 1024 {
		return false
	}
	
	// 检查是否为有效的base64编码
	decoded, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return false
	}
	
	// 验证解码后的长度
	return len(decoded) >= 16 && len(decoded) <= 256
}

/**
 * SanitizeInput 清理输入数据
 * @param {string} input 输入数据
 * @returns {string} 清理后的数据
 */
func SanitizeInput(input string) string {
	// 移除潜在的危险字符
	sanitized := ""
	for _, char := range input {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_' || char == '.' {
			sanitized += string(char)
		}
	}
	return sanitized
}
