package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/wando-world/wando-sso/internal/utils"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	pwUtils := &utils.PasswordUtils{}

	salt, err := pwUtils.GenerateSalt()
	assert.NoError(t, err)
	assert.NotNil(t, salt)
	assert.Equal(t, utils.SaltLength, len(salt))
}

func TestHashPassword(t *testing.T) {
	pwUtils := &utils.PasswordUtils{}

	password := "testpassword"
	salt, err := pwUtils.GenerateSalt()
	assert.NoError(t, err)

	hashedPassword := pwUtils.HashPassword(password, salt)
	assert.NotEmpty(t, hashedPassword)
}

func TestVerifyPassword(t *testing.T) {
	pwUtils := &utils.PasswordUtils{}

	password := "testpassword"
	salt, err := pwUtils.GenerateSalt()
	assert.NoError(t, err)

	hashedPassword := pwUtils.HashPassword(password, salt)
	assert.True(t, pwUtils.VerifyPassword(password, hashedPassword, salt))

	assert.False(t, pwUtils.VerifyPassword("wrongpassword", hashedPassword, salt))
}
