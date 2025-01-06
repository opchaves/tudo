package handlers

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/tudo/internal/config"
	"github.com/opchaves/tudo/internal/models"
)

type Container struct {
	DB  *pgxpool.Pool
	Q   *models.Queries
	JWT *jwtauth.JWTAuth
}

func NewContainer() *Container {
	pool, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	return &Container{
		DB:  pool,
		Q:   models.New(pool),
		JWT: jwtauth.New("HS256", []byte(config.JwtSecret), nil),
	}
}

func logCtx(r *http.Request) *slog.Logger {
	return httplog.LogEntry(r.Context())
}
