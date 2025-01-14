package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var (
    _, b, _, _ = runtime.Caller(0)
    basepath   = filepath.Dir(b)
)

var doOnce sync.Once

type ctxKey int

const (
	CtxClaims ctxKey = iota
	CtxRefreshToken
	CtxVersion
)

var (
	Name    = getEnv("TUDO_NAME", "tudo")
	Env     = getEnv("TUDO_ENV", "development")
	Host    = getEnv("HOST", "0.0.0.0")
	Port    = getEnv("PORT", "8080")
	Origins = getEnv("ORIGINS", "")

	IsProduction = Env == "production"
	IsLocal      = Env == "development" || Env == "test"
	DatabaseURL  = getEnv("DATABASE_URL", "postgres://dev:password@localhost:5432/devdb?sslmode=disable")
	TestDatabaseURL  = getEnv("TEST_DATABASE_URL", "postgres://dev:password@localhost:5432/testdb?sslmode=disable")

	JwtSecret        = getEnv("JWT_SECRET", "superSecret")
	JwtExpiry        = toDuration("JWT_EXPIRY", "1h")
	JwtRefreshExpiry = toDuration("JWT_REFRESH_EXPIRY", "72h")
)

func toDuration(envVar string, defaultVal string) time.Duration {
	val, err := time.ParseDuration(getEnv(envVar, defaultVal))
	if err != nil {
		log.Fatalf("Invalid value for %s: %s", envVar, err)
	}
	return val
}

func toInt(envVar string, defaultVal string) int {
	val, err := strconv.Atoi(getEnv(envVar, defaultVal))
	if err != nil {
		log.Fatalf("Invalid value for %s: %s", envVar, err)
	}
	return val
}

func getEnv(name, defaultValue string) string {
	doOnce.Do(func() {
		path := filepath.Join(basepath, "../../.env")
		println(">>> .env path", path)
		readEnvFile(path)
	})

	if value := os.Getenv(name); value != "" {
		return value
	}

	return defaultValue
}

func readEnvFile(filename string) {
	env := os.Getenv("APP_ENV")
	if env != "production" {
		err := godotenv.Load(filename)
		if err != nil {
			log.Printf("No %s file found. Using default values.\n", filename)
		}
	}
}
