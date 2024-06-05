package db

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

// 테스트 데이터베이스 연결 문자열 설정
const testDSN = "host=localhost user=testuser password=testpass dbname=testdb port=5555 sslmode=disable TimeZone=Asia/Seoul"

func TestInitDB(t *testing.T) {
	// 테스트 환경변수 설정
	err := os.Setenv("GO_ENV", "test")
	require.NoError(t, err, "[에러] test env 세팅 에러")

	// 데이터베이스 연결 테스트
	DB, err = gorm.Open(postgres.Open(testDSN), &gorm.Config{})
	require.NoError(t, err, "[에러] db 커넥션 에러")

	// 테스트 환경에서 데이터베이스 마이그레이션
	performMigration(os.Getenv("GO_ENV"))
}

func TestPerformMigrationInProd(t *testing.T) {
	// 로그를 버퍼로 리다이렉트
	logOutput := bytes.NewBufferString("")
	log.SetOutput(logOutput)
	defer log.SetOutput(os.Stderr) // 테스트 후 로그 출력을 다시 원래대로 설정

	// 테스트 환경변수 설정
	err := os.Setenv("GO_ENV", "prod")
	require.NoError(t, err, "환경번수 세팅 에러")

	// 마이그레이션 함수 실행
	performMigration(os.Getenv("GO_ENV"))

	// 로그 메시지 검증
	require.Contains(t, logOutput.String(), "운영 환경에서는 데이터베이스 마이그레이션을 건너뜁니다.", "prod 인경우 마이그레이션 건너뛰어야 하는데 안됨")
}
