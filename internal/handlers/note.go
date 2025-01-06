package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/tudo/internal/models"
)

func getUserID(r *http.Request) int {
	_, claims, _ := jwtauth.FromContext(r.Context())
	return int(claims["user_id"].(float64))
}

func GetNotes(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)

		rows, err := pool.Query(context.Background(), "SELECT id, user_id, title, content, created_at FROM notes WHERE user_id=$1", userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var notes []models.Note
		for rows.Next() {
			var note models.Note
			if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			notes = append(notes, note)
		}

		json.NewEncoder(w).Encode(notes)
	}
}

func CreateNote(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)

		var note models.Note
		if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		note.UserID = userID

		_, err := pool.Exec(context.Background(), "INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3)", note.UserID, note.Title, note.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func GetNoteByID(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)

		id := chi.URLParam(r, "id")
		noteID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		var note models.Note
		err = pool.QueryRow(context.Background(), "SELECT id, user_id, title, content, created_at FROM notes WHERE id=$1 AND user_id=$2", noteID, userID).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt)
		if err != nil {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(note)
	}
}

func UpdateNoteByID(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)

		id := chi.URLParam(r, "id")
		noteID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		var note models.Note
		if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = pool.Exec(context.Background(), "UPDATE notes SET title=$1, content=$2 WHERE id=$3 AND user_id=$4", note.Title, note.Content, noteID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteNoteByID(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)

		id := chi.URLParam(r, "id")
		noteID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		_, err = pool.Exec(context.Background(), "DELETE FROM notes WHERE id=$1 AND user_id=$2", noteID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
