package testutils

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/opchaves/tudo/internal/config"
	"github.com/opchaves/tudo/internal/models"
	"github.com/opchaves/tudo/migrations"
	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
)

var (
	Pool *pgxpool.Pool
	DB   *sql.DB
)

// @see https://stackoverflow.com/a/77581618

func SetupDB() {
	var err error
	Pool, err = pgxpool.New(context.Background(), config.TestDatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	conn, err := Pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Failed to acquire connection: %v", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
	`)
	if err != nil {
		log.Fatalf("Failed to reset schema: %v", err)
	}

	goose.SetBaseFS(migrations.MigrationFiles)

	if goose.SetDialect("postgres") != nil {
		log.Fatalf("Failed to set goose dialect")
	}

	DB = stdlib.OpenDBFromPool(Pool)

	if err := goose.Up(DB, "."); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

func TeardownDB() {
	log.Println("Closing database connection pool")
	Pool.Close()
}

func CreateUser(email, password string, q *models.Queries) *models.User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	userParams := models.UsersInsertParams{
		Email:    email,
		Password: string(hashedPassword),
	}

	user, err := q.UsersInsert(context.Background(), userParams)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	return user
}

func CreateToken(user *models.User, jwt *jwtauth.JWTAuth) string {
	_, tokenString, err := jwt.Encode(map[string]interface{}{
		"email":   user.Email,
		"user_id": user.ID,
		"exp":     time.Now().Add(config.JwtExpiry).Unix(),
	})
	if err != nil {
		log.Fatalf("Failed to create token: %v", err)
	}

	return tokenString
}
