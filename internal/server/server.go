package server

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/tudo/internal/config"
	"github.com/opchaves/tudo/internal/handlers"
	"github.com/opchaves/tudo/internal/models"
)

type Server struct {
	Router *chi.Mux
	DB     *pgxpool.Pool
	Q      *models.Queries
	JWT    *jwtauth.JWTAuth
}

func CreateNewServer(pool *pgxpool.Pool) *Server {
	if pool == nil {
		var err error
		pool, err = pgxpool.New(context.Background(), config.DatabaseURL)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &Server{
		Router: chi.NewRouter(),
		DB:     pool,
		Q:      models.New(pool),
		JWT:    jwtauth.New("HS256", []byte(config.JwtSecret), nil),
	}
}

func (s *Server) MountHandlers() {
	logLevel := slog.LevelDebug
	if config.IsProduction {
		logLevel = slog.LevelInfo
	}

	s.Router.Use(middleware.RequestID)
	if config.Env != "test" {
		logger := httplog.NewLogger("tudo", httplog.Options{
			JSON:            config.IsProduction,
			LogLevel:        logLevel,
			Concise:         !config.IsProduction,
			RequestHeaders:  !config.IsProduction,
			QuietDownRoutes: []string{"/ping"},
			QuietDownPeriod: 10 * time.Second,
			Tags: map[string]string{
				"env": config.Env,
			},
		})
		s.Router.Use(httplog.RequestLogger(logger))
	}
	s.Router.Use(middleware.Recoverer)

	handler := handlers.NewHandler(s.DB)

	s.Router.Route("/auth", func(r chi.Router) {
		r.Post("/signup", handler.SignUp)
		r.Post("/login", handler.Login)
	})

	s.Router.Route("/api/notes", func(r chi.Router) {
		r.Use(jwtauth.Verifier(s.JWT))
		r.Use(jwtauth.Authenticator(s.JWT))
		r.Get("/", handler.GetNotes)
		r.Post("/", handler.CreateNote)
		r.Get("/{id}", handler.GetNoteByID)
		r.Put("/{id}", handler.UpdateNoteByID)
		r.Delete("/{id}", handler.DeleteNoteByID)
	})
}
