package test

import (
	"github.com/stretchr/testify/assert"
	db2 "github.com/wando-world/wando-sso/internal/db"
	"github.com/wando-world/wando-sso/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupUserTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})

	if err != nil {
		panic("sqlite 커넥션 실패")
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic("마이그레이션 실패")
	}
	return db
}

func TestCreateUser(t *testing.T) {
	db := setupUserTestDB()
	repo := db2.NewUserRepository(db)

	email := "wando@naver.com"
	testUser := &models.User{
		// 모델에 맞는 필드값 설정
		Nickname:     "wando",
		UserID:       "wando123",
		Password:     "password123",
		Email:        &email,
		Salt:         "123",
		VerifiedCode: "123456",
		Role:         "GENERAL",
	}

	// CreateUser를 호출하고 결과를 검사
	err := repo.CreateUser(testUser)
	assert.NoError(t, err, "회원가입 성공")

	// 데이터가 실제로 삽입되었는지 검증
	var user models.User
	result := db.First(&user, "nickname = ?", "wando")
	assert.NoError(t, result.Error, "회원가입 후 조회 성공")
	assert.Equal(t, "wando@naver.com", *user.Email, "이메일 맞는지 확인")

	// 이미 있는 유저 생성시 에러발생
	err = repo.CreateUser(testUser)
	assert.Error(t, err, "이미 있는 유저이므로 에러 발생")
}

func TestFindUserForLogin(t *testing.T) {
	db := setupUserTestDB()
	repo := db2.NewUserRepository(db)

	// case 1
	email := "wando@naver.com"
	testUser := &models.User{
		// 모델에 맞는 필드값 설정
		Nickname:     "wando",
		UserID:       "wando123",
		Password:     "password123",
		Email:        &email,
		Salt:         "123",
		VerifiedCode: "123456",
		Role:         "GENERAL",
	}

	_ = repo.CreateUser(testUser)

	foundUser, err := repo.FindUserForLogin(testUser)
	assert.NoError(t, err, "조회 성공")
	assert.Equal(t, "wando@naver.com", *foundUser.Email, "이메일 확인")
	assert.Equal(t, "wando", foundUser.Nickname, "닉네임 확인")

	// case 2
	testUser = &models.User{
		// 모델에 맞는 필드값 설정
		Nickname:     "wando33",
		UserID:       "wando3333",
		Password:     "password123",
		Email:        &email,
		Salt:         "123",
		VerifiedCode: "123456",
		Role:         "GENERAL",
	}
	_ = repo.CreateUser(testUser)

	foundUser, err = repo.FindUserForLogin(testUser)
	assert.NoError(t, err, "조회 성공")
	assert.Equal(t, "wando33", foundUser.Nickname, "이메일 확인")
	assert.Equal(t, "wando3333", foundUser.UserID, "닉네임 확인")
}

func TestFindUserByID(t *testing.T) {
	db := setupUserTestDB()
	repo := db2.NewUserRepository(db)

	email := "wando@naver.com"
	testUser := &models.User{
		// 모델에 맞는 필드값 설정
		Nickname:     "wando",
		UserID:       "wando123",
		Password:     "password123",
		Email:        &email,
		Salt:         "123",
		VerifiedCode: "123456",
		Role:         "GENERAL",
	}
	_ = repo.CreateUser(testUser)

	expectedUser, _ := repo.FindUserForLogin(testUser)

	foundUser, err := repo.FindUserByID(expectedUser)
	assert.NoError(t, err, "조회 성공")
	assert.Equal(t, expectedUser.ID, foundUser.ID, "고유 id 확인")
}
