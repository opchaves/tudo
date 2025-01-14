package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/opchaves/tudo/internal/handlers"
)

func TestSignUp(t *testing.T) {
	t.Run("successful signup", func(t *testing.T) {
		input := handlers.SignUpRequest{Email: "test@example.com", Password: "pass123*"}
		req := makePostRequest(t, "/auth/signup", input)

		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusCreated)
	})

	t.Run("invalid payload", func(t *testing.T) {
		input := []byte(`{"invalid": "payload"}`)
		req := makePostRequest(t, "/auth/signup", input)
		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("invalid email", func(t *testing.T) {
		input := handlers.SignUpRequest{Email: "invalid-email", Password: "password123"}
		req := makePostRequest(t, "/auth/signup", input)
		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("duplicate email", func(t *testing.T) {
		input := handlers.SignUpRequest{Email: "duplicate@example.com", Password: "password123"}
		req := makePostRequest(t, "/auth/signup", input)

		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusCreated)

		req = makePostRequest(t, "/auth/signup", input)
		rr = execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})
}

func TestLogin(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		input := handlers.LoginRequest{Email: authUser.Email, Password: userPassword}
		req := makePostRequest(t, "/auth/login", input)

		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusOK)

		var body map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &body)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if _, ok := body["token"]; !ok {
			t.Fatalf("Response does not contain token")
		}
	})

	t.Run("should fail with invalid password", func(t *testing.T) {
		input := handlers.LoginRequest{Email: authUser.Email, Password: "_invalid_"}
		req := makePostRequest(t, "/auth/login", input)

		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})
}
