package postgresql

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func setupAuthMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}

	return gormDB, mock, func() {
		err := db.Close()
		if err != nil {
			return
		}
	}
}

func TestAuthRepository_FindUserForLogin(t *testing.T) {
	// 데이터베이스 설정
	gormDB, mock, teardown := setupAuthMockDB(t)
	defer teardown()

	repo := NewAuthRepository(gormDB)

	// 예상되는 쿼리와 결과 설정
	rows := sqlmock.NewRows([]string{"id", "user_id", "verified_code", "nickname", "email"}).
		AddRow(1, "testuser", "123456", "Test User", "test@example.com")

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE \("users"\."user_id" = \$1 AND "users"\."verified_code" = \$2\) AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$3`).
		WithArgs("testuser", "123456", 1).
		WillReturnRows(rows)

	// 테스트 실행
	resultUser, err := repo.FindUserForLogin(context.Background(), "testuser", "123456")
	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.Equal(t, "testuser", resultUser.UserID)
	assert.Equal(t, "123456", resultUser.VerifiedCode)
	assert.Equal(t, "Test User", resultUser.Nickname)
	assert.Equal(t, "test@example.com", *resultUser.Email)

	// 예상된 쿼리가 실행되었는지 확인
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

}

func TestAuthRepository_FindUserById(t *testing.T) {
	// 데이터베이스 설정
	gormDB, mock, teardown := setupAuthMockDB(t)
	defer teardown()

	repo := NewAuthRepository(gormDB)

	// 예상되는 쿼리와 결과 설정
	rows := sqlmock.NewRows([]string{"id", "user_id", "verified_code", "nickname", "email"}).
		AddRow(1, "testuser", "123456", "Test User", "test@example.com")

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	// 테스트 실행
	resultUser, err := repo.FindUserById(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.Equal(t, uint(1), resultUser.ID)
	assert.Equal(t, "testuser", resultUser.UserID)
	assert.Equal(t, "123456", resultUser.VerifiedCode)
	assert.Equal(t, "Test User", resultUser.Nickname)
	assert.Equal(t, "test@example.com", *resultUser.Email)

	// 예상된 쿼리가 실행되었는지 확인
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
