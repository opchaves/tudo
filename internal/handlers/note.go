package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/opchaves/tudo/internal/models"
)

func getUserID(r *http.Request) int32 {
	_, claims, _ := jwtauth.FromContext(r.Context())
	return int32(claims["user_id"].(float64))
}

func GetNotes(c *Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)

		notes, err := c.Q.NotesGetByUser(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(notes)
	}
}

func CreateNote(c *Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)

		var note models.NotesInsertParams
		if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		note.UserID = userID

		_, err := c.Q.NotesInsert(r.Context(), note)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func GetNoteByID(c *Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)

		id := chi.URLParam(r, "id")
		noteID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		note, err := c.Q.NotesGetByID(r.Context(), models.NotesGetByIDParams{ID: int32(noteID), UserID: userID})
		if err != nil {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(note)
	}
}

func UpdateNoteByID(c *Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)

		id := chi.URLParam(r, "id")
		noteID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		var note models.NotesUpdateParams
		if err = json.NewDecoder(r.Body).Decode(&note); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		note.ID = int32(noteID)
		note.UserID = userID

		_, err = c.Q.NotesUpdate(r.Context(), note)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteNoteByID(c *Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)

		id := chi.URLParam(r, "id")
		noteID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		err = c.Q.NotesDelete(r.Context(), models.NotesDeleteParams{ID: int32(noteID), UserID: userID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
