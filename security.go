package wtoken

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

/**
 * SecurityManager 安全管理器
 */
type SecurityManager struct {
	key []byte // 加密密钥
}

/**
 * NewSecurityManager 创建安全管理器
 * @param {string} password 密码
 * @returns {*SecurityManager} 安全管理器实例
 */
func NewSecurityManager(password string) *SecurityManager {
	hash := sha256.Sum256([]byte(password))
	return &SecurityManager{
		key: hash[:],
	}
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
 * @returns {bool} 是否有效
 */
func ValidateTokenFormat(token string) bool {
	if len(token) == 0 {
		return false
	}
	// 检查是否为有效的base64编码
	_, err := base64.URLEncoding.DecodeString(token)
	return err == nil
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