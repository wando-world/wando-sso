package utils

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
)

const (
	SaltLength = 16
	HashLength = 32
)

// GenerateSalt 랜덤 salt 생성
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// HashPassword argon2 로 비밀번호 해시
func HashPassword(password string, salt []byte) string {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, HashLength)
	return base64.RawStdEncoding.EncodeToString(hash)
}

// VerifyPassword 일반 텍스트 비밀번호를 해시와 비교
func VerifyPassword(password, encodedHash string, salt []byte) bool {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, HashLength)
	return base64.RawStdEncoding.EncodeToString(hash) == encodedHash
}
