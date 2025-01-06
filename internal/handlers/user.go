package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/opchaves/tudo/internal/config"
	"github.com/opchaves/tudo/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aLog := logCtx(r)
		var user models.UsersInsertParams

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Hash the user's password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)

		_, err = c.Q.UsersInsert(r.Context(), user)
		if err != nil {
			aLog.Warn("Failed to insert user", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		aLog.Info("User created", "email", user.Email)

		w.WriteHeader(http.StatusCreated)
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dbUser, err := c.Q.UsersFindByEmail(r.Context(), user.Email)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		_, tokenString, err := c.JWT.Encode(map[string]interface{}{
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
}
