package handlers_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/opchaves/tudo/internal/handlers"
	"github.com/opchaves/tudo/internal/testutils"
)

func TestSignUp(t *testing.T) {
	t.Run("successful signup", func(t *testing.T) {
		input := handlers.SignUpRequest{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "Test",
		}
		req := makePostRequest(t, "/auth/signup", input)

		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusCreated)
	})

	t.Run("should fail with invalid payload", func(t *testing.T) {
		input := []byte(`{"invalid": "payload"}`)
		req := makePostRequest(t, "/auth/signup", input)
		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("should fail with invalid email", func(t *testing.T) {
		input := handlers.SignUpRequest{Email: "invalid-email", Password: "password123"}
		req := makePostRequest(t, "/auth/signup", input)
		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("should fail with short password", func(t *testing.T) {
		input := handlers.SignUpRequest{Email: authUser.Email, Password: "password"}
		req := makePostRequest(t, "/auth/signup", input)
		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("should fail with duplicate email", func(t *testing.T) {
		input := handlers.SignUpRequest{
			Email:     "duplicate@example.com",
			Password:  "password123",
			FirstName: "Duplicate",
		}
		req := makePostRequest(t, "/auth/signup", input)

		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusCreated)

		req = makePostRequest(t, "/auth/signup", input)
		rr = execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})
}

func TestLogin(t *testing.T) {
	t.Run("should login if verified", func(t *testing.T) {
		input := handlers.LoginRequest{Email: authUser.Email, Password: userPassword}
		req := makePostRequest(t, "/auth/login", input)

		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusOK)

		body := parseBody(t, rr.Body.Bytes())

		if _, ok := body["token"]; !ok {
			t.Fatalf("Response does not contain token")
		}
	})

	t.Run("should fail if not verified yet", func(t *testing.T) {
		password := "password123"
		user := testutils.CreateUser("unverified@example.com", password, false, aServer.Q)
		input := handlers.LoginRequest{Email: user.Email, Password: password}

		req := makePostRequest(t, "/auth/login", input)
		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusBadRequest)

		body := parseBody(t, rr.Body.Bytes())

		if msg, ok := body["error"].(string); ok {
			if !strings.Contains(msg, "User not verified yet") {
				t.Fatalf("Response does not contain message")
			}
		} else {
			t.Fatalf("Response does not contain message")
		}
	})

	t.Run("should fail with invalid password", func(t *testing.T) {
		input := handlers.LoginRequest{Email: authUser.Email, Password: "_invalid_"}
		req := makePostRequest(t, "/auth/login", input)

		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("should fail if email not found", func(t *testing.T) {
		input := handlers.LoginRequest{Email: "notfound123@test.com", Password: "password123"}
		req := makePostRequest(t, "/auth/login", input)

		rr := execRequest(req, aServer)

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})
}
