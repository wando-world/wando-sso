package db

import (
	"github.com/wando-world/wando-sso/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		log.Fatalf("[에러] db 커넥션 에러 %v", err)
	}

	log.Println("DB 커넥션 성공")

	// ENUM 검사 및 없으면 생성
	if err := ensureRoleTypeExists(); err != nil {
		log.Fatalf("[에러] role_type ENUM 생성 에러: %v", err)
	}

	// 환경 설정에 따른 마이그레이션 실행
	performMigration(os.Getenv("GO_ENV"))
}

// ensureRoleTypeExists 'role_type' ENUM 있는지 체크, 없다면 생성
func ensureRoleTypeExists() error {
	var exists bool
	err := DB.Raw("SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_type')").Scan(&exists).Error
	if err != nil {
		return err
	}

	// ENUM 타입이 이미 존재하면 함수 종료
	if exists {
		log.Println("'role_type' ENUM 이미 존재함.")
		return nil
	}

	// ENUM 타입 생성
	err = DB.Exec("CREATE TYPE role_type AS ENUM ('ADMIN', 'GENERAL')").Error
	if err != nil {
		return err
	}

	log.Println("'role_type' ENUM 생성 완료.")
	return nil
}

// performMigration 환경 설정에 따라 데이터베이스 마이그레이션을 수행
func performMigration(environment string) {
	if environment != "prod" {
		migrateDatabase()
	} else {
		log.Println("운영 환경에서는 데이터베이스 마이그레이션을 건너뜁니다.")
	}
}

// migrateDatabase 모델에 대한 데이터베이스 마이그레이션을 수행
func migrateDatabase() {
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("[에러] 마이그레이션 실패: %v", err)
	}
	log.Println("데이터베이스 마이그레이션 성공적으로 완료됨")
}
