package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/opchaves/tudo/internal/models"
	"github.com/opchaves/tudo/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
}

func (o *SignUpRequest) Bind(r *http.Request) error {
	if o == nil {
		return errors.New("missing required fields")
	}

	return validate.Struct(o)
}

type SignUpResponse struct {
	UserID string `json:"id"`
}

func (s *SignUpResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	data := &SignUpRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	hashedPassword := string(hash)

	newUser := models.UsersInsertParams{}
	newUser.Email = data.Email
	newUser.Password = &hashedPassword
	newUser.Uid = utils.NewIDShort()
	newUser.FirstName = data.FirstName
	newUser.LastName = &data.LastName
	newUser.Role = "user"

	user, err := h.Q.UsersInsert(r.Context(), newUser)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	aLog(r).Info("User created", "email", user.Email)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &SignUpResponse{UserID: user.Uid})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	*models.User
	ID       bool   `json:"id,omitempty"`
	Password bool   `json:"password,omitempty"`
	Token    string `json:"token"`
}

func (o *LoginRequest) Bind(r *http.Request) error {
	if o == nil {
		return errors.New("missing required fields")
	}

	return validate.Struct(o)
}

func (s *LoginResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	input := &LoginRequest{}
	if err := render.Bind(r, input); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user, err := h.Q.UsersFindByEmail(r.Context(), input.Email)
	if err != nil {
		render.Render(w, r, ErrText("Email already taken"))
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(input.Password)); err != nil {
		render.Render(w, r, ErrText("Invalid email or password"))
		return
	}

	token, err := NewToken(user, h.JWT)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &LoginResponse{User: user, Token: token})
}
