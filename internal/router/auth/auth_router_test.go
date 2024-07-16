package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/wando-world/wando-sso/internal/api"
	apiModels "github.com/wando-world/wando-sso/internal/api/models"
	"github.com/wando-world/wando-sso/internal/models"
	"github.com/wando-world/wando-sso/internal/utils"
	"testing"
)

type MockAuthMapper struct {
	mock.Mock
}

func (m *MockAuthMapper) LoginRequestToUser(req apiModels.LoginRequest) models.User {
	args := m.Called(req)
	return args.Get(0).(models.User)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindUserForLogin(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindUserByID(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

type MockPasswordUtils struct {
	mock.Mock
}

func (m *MockPasswordUtils) GenerateSalt() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockPasswordUtils) HashPassword(password string, salt []byte) string {
	args := m.Called(password, salt)
	return args.String(0)
}

func (m *MockPasswordUtils) VerifyPassword(password, encodedHash string, salt []byte) bool {
	args := m.Called(password, encodedHash, salt)
	return args.Bool(0)
}

type MockJwtUtils struct {
	mock.Mock
}

func (m *MockJwtUtils) GenerateATK(id uint, role string) (string, error) {
	args := m.Called(id, role)
	return args.String(0), args.Error(1)
}

func (m *MockJwtUtils) GenerateRTK(id uint) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func (m *MockJwtUtils) ParseToken(tokenString string) (*utils.Claims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*utils.Claims), args.Error(1)
}

func (m *MockJwtUtils) GetAtkSecret() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func (m *MockJwtUtils) GetRtkSecret() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

type testEnv struct {
	echo              *echo.Echo
	mockMapper        *MockAuthMapper
	mockRepo          *MockUserRepository
	handler           api.IAuthHandler
	mockPasswordUtils *MockPasswordUtils
	mockJwtUtils      *MockJwtUtils
}

func setupTestEnv() testEnv {
	e := echo.New()
	e.Validator = utils.NewValidator()

	mockMapper := new(MockAuthMapper)
	mockRepo := new(MockUserRepository)
	mockPasswordUtils := new(MockPasswordUtils)
	mockJwtUtils := new(MockJwtUtils)
	handler := api.NewAuthHandler(mockPasswordUtils, mockJwtUtils, mockMapper, mockRepo)

	return testEnv{
		echo:              e,
		mockMapper:        mockMapper,
		mockRepo:          mockRepo,
		mockPasswordUtils: mockPasswordUtils,
		mockJwtUtils:      mockJwtUtils,
		handler:           handler,
	}
}

func TestLogin(t *testing.T) {
	env := setupTestEnv()

	tests := []struct {
		name         string
		loginRequest string
		setupMocks   func()
		wantStatus   int
		wantError    bool
	}{
		{
			name:         "정상 Request",
			loginRequest: `{"userId":"wando","password":"1234asd","verifiedCode":"asd123"}`,
			setupMocks: func() {
				// TODO: mocking 해야함
			},
		},
	}
}
