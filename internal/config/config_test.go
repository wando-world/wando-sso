package config_test

import (
	"github.com/wando-world/wando-sso/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setEnv 함수는 환경 변수를 설정하고, 설정 실패 시 테스트를 실패시킴
func setEnv(t *testing.T, key, value string) {
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("failed to set environment variable %s: %v", key, err)
	}
}

func TestNewConfig(t *testing.T) {
	// 설정된 환경 변수 초기화
	os.Clearenv()

	// 필요한 환경 변수 설정
	setEnv(t, "GO_ENV", "test")
	setEnv(t, "PORT", ":8082")
	setEnv(t, "DATABASE_URL", "test_db_url")
	setEnv(t, "ATK_SECRET", "test_atk_secret")
	setEnv(t, "RTK_SECRET", "test_rtk_secret")

	// Config 인스턴스 생성
	cfg := config.New()

	// 설정값 확인
	assert.Equal(t, ":8082", cfg.Port)
	assert.Equal(t, "test", cfg.Env)
	assert.Equal(t, "test_db_url", cfg.DbUrl)
	assert.Equal(t, "test_atk_secret", cfg.ATKSecret)
	assert.Equal(t, "test_rtk_secret", cfg.RTKSecret)
}

func TestGetConfig(t *testing.T) {
	// 설정된 환경 변수 초기화
	os.Clearenv()

	// 필요한 환경 변수 설정
	setEnv(t, "GO_ENV", "test")
	setEnv(t, "PORT", ":8082")
	setEnv(t, "DATABASE_URL", "test_db_url")
	setEnv(t, "ATK_SECRET", "test_atk_secret")
	setEnv(t, "RTK_SECRET", "test_rtk_secret")

	// GetConfig 호출하여 설정값 확인
	cfg := config.GetConfig()
	assert.Equal(t, ":8082", cfg.Port)
	assert.Equal(t, "test", cfg.Env)
	assert.Equal(t, "test_db_url", cfg.DbUrl)
	assert.Equal(t, "test_atk_secret", cfg.ATKSecret)
	assert.Equal(t, "test_rtk_secret", cfg.RTKSecret)
}
