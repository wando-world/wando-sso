package utils

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
	"sync"
)

const (
	SaltLength = 16
	HashLength = 32
)

type IPasswordUtils interface {
	GenerateSalt() ([]byte, error)
	HashPassword(password string, salt []byte) string
	VerifyPassword(password, encodedHash string, salt []byte) bool
}

// PasswordUtils 구현체
type PasswordUtils struct{}

var (
	passwordInstance *PasswordUtils
	passwordOnce     sync.Once
)

func NewPasswordUtils() *PasswordUtils {
	passwordOnce.Do(func() {
		passwordInstance = &PasswordUtils{}
	})
	return passwordInstance
}

// GenerateSalt 랜덤 salt 생성
func (p *PasswordUtils) GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// HashPassword argon2 로 비밀번호 해시
func (p *PasswordUtils) HashPassword(password string, salt []byte) string {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, HashLength)
	return base64.RawStdEncoding.EncodeToString(hash)
}

// VerifyPassword 일반 텍스트 비밀번호를 해시와 비교
func (p *PasswordUtils) VerifyPassword(password, encodedHash string, salt []byte) bool {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, HashLength)
	return base64.RawStdEncoding.EncodeToString(hash) == encodedHash
}
