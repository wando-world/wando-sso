package utils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wando-world/wando-sso/utils"
)

func TestGenerateATK(t *testing.T) {
	atkSecret := "testAtkSecret"
	rtkSecret := "testRtkSecret"
	jwtUtils := utils.NewJwtUtils(atkSecret, rtkSecret)

	id := uint(1)
	role := "general"
	tokenString, err := jwtUtils.GenerateATK(id, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	claims, err := jwtUtils.ParseToken(tokenString, "atk")
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, id, claims.Id)
	assert.Equal(t, role, claims.Role)
	assert.WithinDuration(t, time.Now().Add(time.Minute*15), claims.ExpiresAt.Time, time.Minute)
	assert.Equal(t, "wando", claims.Issuer)
}

func TestGenerateRTK(t *testing.T) {
	atkSecret := "testAtkSecret"
	rtkSecret := "testRtkSecret"
	jwtUtils := utils.NewJwtUtils(atkSecret, rtkSecret)

	id := uint(1)
	tokenString, err := jwtUtils.GenerateRTK(id)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	claims, err := jwtUtils.ParseToken(tokenString, "rtk")
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, id, claims.Id)
	assert.WithinDuration(t, time.Now().Add(time.Hour*24*7), claims.ExpiresAt.Time, time.Minute)
	assert.Equal(t, "wando", claims.Issuer)
}

func TestParseToken(t *testing.T) {
	atkSecret := "testAtkSecret"
	rtkSecret := "testRtkSecret"
	jwtUtils := utils.NewJwtUtils(atkSecret, rtkSecret)

	id := uint(1)
	role := "general"
	tokenString, err := jwtUtils.GenerateATK(id, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	claims, err := jwtUtils.ParseToken(tokenString, "atk")
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, id, claims.Id)
	assert.Equal(t, role, claims.Role)
}

func TestInvalidToken(t *testing.T) {
	atkSecret := "testAtkSecret"
	rtkSecret := "testRtkSecret"
	jwtUtils := utils.NewJwtUtils(atkSecret, rtkSecret)

	invalidTokenString := "invalidToken"
	claims, err := jwtUtils.ParseToken(invalidTokenString, "rtk")
	assert.Error(t, err)
	assert.Nil(t, claims)
}
