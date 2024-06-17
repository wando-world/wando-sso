package test

import (
	"bytes"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wando-world/wando-sso/internal/api"
	models2 "github.com/wando-world/wando-sso/internal/api/models"
	"github.com/wando-world/wando-sso/internal/models"
	"github.com/wando-world/wando-sso/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockUserMapper is a mock for IUserMapper
type MockUserMapper struct {
	mock.Mock
}

func (m *MockUserMapper) CreateUserRequestToUser(req models2.CreateUserRequest) (models.User, error) {
	args := m.Called(req)
	return args.Get(0).(models.User), args.Error(1)
}

// MockUserRepository is a mock for IUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

type testEnv struct {
	echo       *echo.Echo
	mockMapper *MockUserMapper
	mockRepo   *MockUserRepository
	handler    api.IUserHandler
}

func setupTestEnv() testEnv {
	e := echo.New()
	e.Validator = utils.NewValidator()

	mockMapper := new(MockUserMapper)
	mockRepo := new(MockUserRepository)
	handler := api.NewUserHandler(mockMapper, mockRepo)

	return testEnv{
		echo:       e,
		mockMapper: mockMapper,
		mockRepo:   mockRepo,
		handler:    handler,
	}
}

func newCreateUserRequestHelper(userRequest string, env testEnv) (*httptest.ResponseRecorder, echo.Context) {
	req := httptest.NewRequest(http.MethodPost, "/sso/api/v1/user", bytes.NewReader([]byte(userRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := env.echo.NewContext(req, rec)
	return rec, c
}

func TestCreateUser(t *testing.T) {
	env := setupTestEnv()

	// Test cases
	tests := []struct {
		name        string
		userRequest string
		setupMocks  func()
		wantStatus  int
		wantError   bool
	}{
		{
			name:        "정상 Request",
			userRequest: `{"nickname":"wando","userId":"wando1234","password":"wando12345","verifiedCode":"asd123"}`,
			setupMocks: func() {
				env.mockMapper.On("CreateUserRequestToUser", mock.Anything).Return(models.User{Nickname: "wando"}, nil)
				env.mockRepo.On("CreateUser", mock.Anything).Return(nil)
			},
			wantStatus: http.StatusCreated,
			wantError:  false,
		},
		{
			name:        "정상 Request add email",
			userRequest: `{"nickname":"wando","userId":"wando1234","password":"wando12345","email":"kdw1521@naver.com","verifiedCode":"asd123"}`,
			setupMocks: func() {
				env.mockMapper.On("CreateUserRequestToUser", mock.Anything).Return(models.User{Nickname: "wando"}, nil)
				env.mockRepo.On("CreateUser", mock.Anything).Return(nil)
			},
			wantStatus: http.StatusCreated,
			wantError:  false,
		},
		{
			name:        "비정상 Request JSON",
			userRequest: `{"nickname":"wando","userId":"wando1234",}`,
			setupMocks:  func() {},
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
		},
		{
			name:        "비정상 - userId",
			userRequest: `{"nickname":"wando","userId":"wa","password":"wando12345","verifiedCode":"asd123"}`,
			setupMocks:  func() {},
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
		},
		{
			name:        "비정상 - password",
			userRequest: `{"nickname":"wando","userId":"wando1234","password":"12","verifiedCode":"asd123"}`,
			setupMocks:  func() {},
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
		},
		{
			name:        "비정상 - verifiedCode",
			userRequest: `{"nickname":"wando","userId":"wando1234","password":"wando12345","verifiedCode":"12"}`,
			setupMocks:  func() {},
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
		},
		{
			name:        "비정상 - email",
			userRequest: `{"nickname":"wando","email":"ddd","userId":"wando1234","password":"wando12345","verifiedCode":"asd123"}`,
			setupMocks:  func() {},
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			rec, c := newCreateUserRequestHelper(tt.userRequest, env)

			// Invoke the CreateUser method
			err := env.handler.CreateUser(c)

			// Check if an error was expected
			if (err != nil) != tt.wantError {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantError)
			}

			// If an error occurred, assert the status code from the HTTP error
			if err != nil {
				var httpError *echo.HTTPError
				ok := errors.As(err, &httpError)
				if !ok {
					t.Errorf("Expected echo.HTTPError, got %T", err)
				} else {
					assert.Equal(t, tt.wantStatus, httpError.Code, "status code가 다름")
				}
			} else {
				// Assert the status code for successful invocation
				assert.Equal(t, tt.wantStatus, rec.Code)
			}

			// Ensure all expectations are met
			env.mockMapper.AssertExpectations(t)
			env.mockRepo.AssertExpectations(t)
		})
	}
}
