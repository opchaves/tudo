package testutils

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/opchaves/tudo/internal/config"
	"github.com/opchaves/tudo/internal/models"
	"github.com/opchaves/tudo/internal/utils"
	"github.com/opchaves/tudo/migrations"
	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
)

var Pool *pgxpool.Pool

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

	db := stdlib.OpenDBFromPool(Pool)

	if err := goose.Up(db, "."); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

func TeardownDB() {
	log.Println("Closing database connection pool")
	Pool.Close()
}

func CreateUser(email, password string, verified bool, q *models.Queries) *models.User {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	hashedPassword := string(hash)
	uid := utils.NewIDShort()
	lastName := "Test"

	userParams := models.UsersInsertParams{
		Email:     email,
		Password:  &hashedPassword,
		Uid:       uid,
		FirstName: "Test",
		LastName:  &lastName,
		Role:      "user",
		Verified:  verified,
	}

	user, err := q.UsersInsert(context.Background(), userParams)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	return user
}
