package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/opchaves/tudo/internal/config"
	"github.com/opchaves/tudo/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (o *SignUpRequest) Bind(r *http.Request) error {
	if o == nil {
		return errors.New("missing required fields")
	}

	err := validate.Struct(o)

	return err
}

type SignUpResponse struct {
	UserID int32 `json:"id"`
}

func (s *SignUpResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	data := &SignUpRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	newUser := models.UsersInsertParams{}
	newUser.Email = data.Email
	newUser.Password = string(hashedPassword)

	user, err := h.Q.UsersInsert(r.Context(), newUser)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	aLog(r).Info("User created", "email", user.Email)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &SignUpResponse{UserID: user.ID})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbUser, err := h.Q.UsersFindByEmail(r.Context(), user.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	_, tokenString, err := h.JWT.Encode(map[string]interface{}{
		"email":   dbUser.Email,
		"user_id": dbUser.ID,
		"exp":     time.Now().Add(config.JwtExpiry).Unix(),
	})
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    dbUser.ID,
			"name":  dbUser.Name,
			"email": dbUser.Email,
		},
		"token": tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
