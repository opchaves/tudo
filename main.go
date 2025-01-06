package main

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/jwtauth/v5"
	"github.com/opchaves/tudo/internal/config"
	"github.com/opchaves/tudo/internal/handlers"
)

func main() {
	container := handlers.NewContainer()

	r := chi.NewRouter()

	logLevel := slog.LevelDebug
	if config.IsProduction {
		logLevel = slog.LevelInfo
	}

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

	r.Use(middleware.RequestID)
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)

	r.Route("/users", func(r chi.Router) {
		r.Post("/signup", handlers.SignUp(container))
		r.Post("/login", handlers.Login(container))
	})

	r.Route("/notes", func(r chi.Router) {
		r.Use(jwtauth.Verifier(container.JWT))
		r.Use(jwtauth.Authenticator(container.JWT))
		r.Get("/", handlers.GetNotes(container))
		r.Post("/", handlers.CreateNote(container))
		r.Get("/{id}", handlers.GetNoteByID(container))
		r.Put("/{id}", handlers.UpdateNoteByID(container))
		r.Delete("/{id}", handlers.DeleteNoteByID(container))
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
