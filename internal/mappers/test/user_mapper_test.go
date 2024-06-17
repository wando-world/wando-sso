package test

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wando-world/wando-sso/internal/api/models"
	"github.com/wando-world/wando-sso/internal/mappers"
	"testing"
)

// Mocking
type MockUtils struct {
	mock.Mock
}

func (m *MockUtils) GenerateSalt() ([]byte, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockUtils) HashPassword(password string, salt []byte) string {
	args := m.Called(password, salt)
	return args.String(0)
}

func (m *MockUtils) VerifyPassword(password, encodedHash string, salt []byte) bool {
	args := m.Called(password, encodedHash, salt)
	return args.Bool(0)
}

func TestCreateUserRequestToUser(t *testing.T) {
	mockUtils := new(MockUtils)
	mockUtils.On("GenerateSalt").Return([]byte("mockSalt"), nil)
	mockUtils.On("HashPassword", "testpassword", []byte("mockSalt")).Return("hashedPassword")

	mapper := mappers.NewUserMapper(mockUtils)

	email := "testemail@example.com"
	req := models.CreateUserRequest{
		Nickname:     "testnickname",
		UserID:       "testuserid",
		Email:        &email,
		Password:     "testpassword",
		VerifiedCode: "testverifiedcode",
	}

	user, err := mapper.CreateUserRequestToUser(req)

	assert.NoError(t, err)
	assert.Equal(t, req.Nickname, user.Nickname)
	assert.Equal(t, req.UserID, user.UserID)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.VerifiedCode, user.VerifiedCode)
	assert.Equal(t, "GENERAL", user.Role)
	assert.Equal(t, "hashedPassword", user.Password)
	assert.Equal(t, base64.RawStdEncoding.EncodeToString([]byte("mockSalt")), user.Salt)
}

func TestCreateUserRequestToUser_Error(t *testing.T) {
	mockUtils := new(MockUtils)
	mockUtils.On("GenerateSalt").Return(nil, fmt.Errorf("salt generation error"))

	mapper := mappers.NewUserMapper(mockUtils)

	email := "testemail@example.com"
	req := models.CreateUserRequest{
		Nickname:     "testnickname",
		UserID:       "testuserid",
		Email:        &email,
		Password:     "testpassword",
		VerifiedCode: "testverifiedcode",
	}

	user, err := mapper.CreateUserRequestToUser(req)

	assert.Error(t, err)
	assert.Empty(t, user)
}
