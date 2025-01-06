package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	"github.com/opchaves/tudo/internal/db"
	"github.com/opchaves/tudo/internal/handlers"
)

var tokenAuth *jwtauth.JWTAuth

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	pool, err := db.Connect(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	tokenAuth = jwtauth.New("HS256", []byte("your_secret_key"), nil)

	r := chi.NewRouter()

	logger := httplog.NewLogger("tudo", httplog.Options{
		JSON:            false,
		LogLevel:        slog.LevelDebug,
		Concise:         !false,
		RequestHeaders:  !false,
		QuietDownRoutes: []string{"/ping"},
		QuietDownPeriod: 10 * time.Second,
		Tags: map[string]string{
			"env": string("dev"),
		},
	})

	r.Use(middleware.RequestID)
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)

	r.Route("/users", func(r chi.Router) {
		r.Post("/signup", handlers.SignUp(pool))
		r.Post("/login", handlers.Login(pool, tokenAuth))
	})

	r.Route("/notes", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))
		r.Get("/", handlers.GetNotes(pool))
		r.Post("/", handlers.CreateNote(pool))
		r.Get("/{id}", handlers.GetNoteByID(pool))
		r.Put("/{id}", handlers.UpdateNoteByID(pool))
		r.Delete("/{id}", handlers.DeleteNoteByID(pool))
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
