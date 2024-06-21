package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port      string
	Env       string
	DbUrl     string
	ATKSecret string
	RTKSecret string
}

func New() *Config {
	// GO_ENV 로 해당 하는 .env 파일 로드
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "dev" // default
	}

	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("[에러] %s file 불러오기 실패", envFile)
	}
	return &Config{
		Port:      getEnv("PORT", ":8081"), // default port
		Env:       env,
		DbUrl:     getEnv("DATABASE_URL", ""),
		ATKSecret: getEnv("ATK_SECRET", ""),
		RTKSecret: getEnv("RTK_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
