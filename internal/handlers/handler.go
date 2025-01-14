package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/tudo/internal/config"
	"github.com/opchaves/tudo/internal/models"
)

var validate = validator.New()

type Container struct {
	DB  *pgxpool.Pool
	Q   *models.Queries
	JWT *jwtauth.JWTAuth
}

type Handler struct {
	DB  *pgxpool.Pool
	Q   *models.Queries
	JWT *jwtauth.JWTAuth
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{
		DB:  db,
		Q:   models.New(db),
		JWT: jwtauth.New("HS256", []byte(config.JwtSecret), nil),
	}
}

func ContainerWithDB(db *pgxpool.Pool) *Container {
	return &Container{
		DB:  db,
		Q:   models.New(db),
		JWT: jwtauth.New("HS256", []byte(config.JwtSecret), nil),
	}
}

func aLog(r *http.Request) *slog.Logger {
	return httplog.LogEntry(r.Context())
}

func getUserID(r *http.Request) int32 {
	_, claims, _ := jwtauth.FromContext(r.Context())
	return int32(claims["user_id"].(float64))
}

func NewToken(user *models.User, jwt *jwtauth.JWTAuth) (string, error) {
	_, tokenString, err := jwt.Encode(map[string]interface{}{
		"email":   user.Email,
		"user_id": user.ID,
		"exp":     time.Now().Add(config.JwtExpiry).Unix(),
	})

	return tokenString, err
}
