package postgresql_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/wando-world/wando-sso/domain"
	"github.com/wando-world/wando-sso/internal/repository/postgresql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func setupUserMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
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

func TestUserRepository_CreateUser(t *testing.T) {
	gormDB, mock, teardown := setupUserMockDB(t)
	defer teardown()

	repo := postgresql.NewUserRepository(gormDB)

	user := &domain.User{
		Nickname:     "testuser",
		UserID:       "testuser",
		Password:     "password",
		Salt:         "salt",
		Email:        new(string),
		VerifiedCode: "123456",
		Role:         "GENERAL",
	}
	*user.Email = "test@example.com"

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" \("created_at","updated_at","deleted_at","nickname","user_id","password","salt","verified_code","role","email"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9,\$10\) RETURNING "id","email"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, user.Nickname, user.UserID, user.Password, user.Salt, user.VerifiedCode, user.Role, *user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, *user.Email))
	mock.ExpectCommit()

	err := repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_FindUserById(t *testing.T) {
	gormDB, mock, teardown := setupUserMockDB(t)
	defer teardown()

	repo := postgresql.NewUserRepository(gormDB)

	// 예상되는 쿼리와 결과 설정
	userID := uint(1)
	rows := sqlmock.NewRows([]string{"id", "nickname", "user_id", "password", "salt", "email", "verified_code", "role"}).
		AddRow(1, "testuser", "testuser", "password", "salt", "test@example.com", "123456", "GENERAL")

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(userID, 1).
		WillReturnRows(rows)

	// 테스트 실행
	resultUser, err := repo.FindUserById(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.Equal(t, "testuser", resultUser.Nickname)
	assert.Equal(t, "testuser", resultUser.UserID)
	assert.Equal(t, "password", resultUser.Password)
	assert.Equal(t, "salt", resultUser.Salt)
	assert.Equal(t, "test@example.com", *resultUser.Email)
	assert.Equal(t, "123456", resultUser.VerifiedCode)
	assert.Equal(t, "GENERAL", resultUser.Role)

	// 예상된 쿼리가 실행되었는지 확인
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
