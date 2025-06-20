package configs

import (
	"fmt"
	"os"
	"sicantik-idaman/internal/domain"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(("failed load env file"))
	}
}

func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		return defaultValue
	}

	return value
}

func SetHelper(hlp *domain.Helper) {
	hlp.Port = GetEnv("APP_PORT", "1010")
	hlp.DbDsn = GetEnv("DSN_DB", "http://localhost")
	hlp.JwtSecret = GetEnv("JWT_SECRET", "dundermifflin")
}
