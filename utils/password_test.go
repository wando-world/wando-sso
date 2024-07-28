package utils_test

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"github.com/wando-world/wando-sso/utils"
	"testing"
)

const (
	SaltLength = 16
	HashLength = 32
)

func TestGenerateSalt(t *testing.T) {
	passwordUtils := utils.NewPasswordUtils()

	salt, err := passwordUtils.GenerateSalt()
	assert.NoError(t, err)
	assert.NotNil(t, salt)
	assert.Equal(t, SaltLength, len(salt))
}

func TestHashPassword(t *testing.T) {
	passwordUtils := utils.NewPasswordUtils()

	salt, err := passwordUtils.GenerateSalt()
	assert.NoError(t, err)

	password := "testpassword"
	hashedPassword := passwordUtils.HashPassword(password, salt)

	// Argon2 해시 길이는 HashLength여야 함
	decodedHash, err := base64.RawStdEncoding.DecodeString(hashedPassword)
	assert.NoError(t, err)
	assert.Equal(t, HashLength, len(decodedHash))
}

func TestVerifyPassword(t *testing.T) {
	passwordUtils := utils.NewPasswordUtils()

	salt, err := passwordUtils.GenerateSalt()
	assert.NoError(t, err)

	password := "testpassword"
	hashedPassword := passwordUtils.HashPassword(password, salt)

	// VerifyPassword 함수가 올바르게 비밀번호를 검증하는지 확인
	isValid := passwordUtils.VerifyPassword(password, hashedPassword, salt)
	assert.True(t, isValid)

	// 잘못된 비밀번호를 검증하는지 확인
	isValid = passwordUtils.VerifyPassword("wrongpassword", hashedPassword, salt)
	assert.False(t, isValid)
}
