package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/opchaves/tudo/internal/models"
)

type NoteResponse struct {
	*models.Note
}

func (n *NoteResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewNoteListResponse(notes []*models.Note) []render.Renderer {
	list := make([]render.Renderer, len(notes))
	for i, note := range notes {
		list[i] = &NoteResponse{note}
	}
	return list
}

type NoteRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (o *NoteRequest) Bind(r *http.Request) error {
	if o == nil {
		return errors.New("missing required fields")
	}

	err := validate.Struct(o)

	return err
}

func GetParamID(w http.ResponseWriter, r *http.Request) (int32, error) {
	id := chi.URLParam(r, "id")
	noteID, err := strconv.Atoi(id)
	if err != nil {
		render.Render(w, r, ErrText("invalid note ID"))
		return 0, err
	}

	return int32(noteID), nil
}

func (h *Handler) GetNotes(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	notes, err := h.Q.NotesGetByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.RenderList(w, r, NewNoteListResponse(notes))
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	input := &NoteRequest{}

	if err := render.Bind(r, input); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	newNote := models.NotesInsertParams{}
	newNote.Title = input.Title
	newNote.Content = input.Content
	newNote.UserID = userID

	note, err := h.Q.NotesInsert(r.Context(), newNote)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &NoteResponse{note})
}

func (h *Handler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	noteID, err := GetParamID(w, r)
	if err != nil {
		return
	}

	params := &models.NotesGetByIDParams{ID: int32(noteID), UserID: userID}
	note, err := h.Q.NotesGetByID(r.Context(), *params)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	render.Render(w, r, &NoteResponse{note})
}

func (h *Handler) UpdateNoteByID(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	noteID, err := GetParamID(w, r)
	if err != nil {
		return
	}

	input := &NoteRequest{}
	if err = render.Bind(r, input); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	noteParams := models.NotesUpdateParams{
		ID:      noteID,
		UserID:  userID,
		Title:   input.Title,
		Content: input.Content,
	}

	note, err := h.Q.NotesUpdate(r.Context(), noteParams)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &NoteResponse{note})
}

func (h *Handler) DeleteNoteByID(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	noteID, err := GetParamID(w, r)
	if err != nil {
		return
	}

	err = h.Q.NotesDelete(r.Context(), models.NotesDeleteParams{ID: noteID, UserID: userID})
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}
