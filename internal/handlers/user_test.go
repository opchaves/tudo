package handlers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/opchaves/tudo/internal/handlers"
	"github.com/opchaves/tudo/internal/models"
	"github.com/opchaves/tudo/internal/server"
	"github.com/opchaves/tudo/internal/testutils"
)

var (
	authUser  *models.User
	authToken string
	aServer   *server.Server
)

var (
	userPassword = "password123"
	userEmail    = "auth@test.com"
)

func TestMain(m *testing.M) {
	var err error

	testutils.SetupDB()

	aServer = server.CreateNewServer(testutils.Pool)

	authUser = testutils.CreateUser(userEmail, userPassword, aServer.Q)

	authToken, err = handlers.NewToken(authUser, aServer.JWT)
	if err != nil {
		log.Fatalf("Failed to create auth token: %v", err)
	}

	aServer.MountHandlers()

	code := m.Run()

	testutils.TeardownDB()
	os.Exit(code)
}

func execRequest(req *http.Request, s *server.Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func MakePostRequest(t *testing.T, path string, data interface{}) *http.Request {
	payload, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	return req
}

func AssertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}

func TestSignUp(t *testing.T) {
	t.Run("successful signup", func(t *testing.T) {
		input := handlers.SignUpRequest{Email: "test@example.com", Password: "pass123*"}
		req := MakePostRequest(t, "/auth/signup", input)

		rr := execRequest(req, aServer)

		AssertStatus(t, rr.Code, http.StatusCreated)
	})

	t.Run("invalid payload", func(t *testing.T) {
		input := []byte(`{"invalid": "payload"}`)
		req := MakePostRequest(t, "/auth/signup", input)
		rr := execRequest(req, aServer)

		AssertStatus(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("invalid email", func(t *testing.T) {
		input := handlers.SignUpRequest{Email: "invalid-email", Password: "password123"}
		req := MakePostRequest(t, "/auth/signup", input)
		rr := execRequest(req, aServer)

		AssertStatus(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("duplicate email", func(t *testing.T) {
		input := handlers.SignUpRequest{Email: "duplicate@example.com", Password: "password123"}
		req := MakePostRequest(t, "/auth/signup", input)

		rr := execRequest(req, aServer)

		AssertStatus(t, rr.Code, http.StatusCreated)

		req = MakePostRequest(t, "/auth/signup", input)
		rr = execRequest(req, aServer)

		AssertStatus(t, rr.Code, http.StatusBadRequest)
	})
}

func TestLogin(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		input := handlers.LoginRequest{Email: authUser.Email, Password: userPassword}
		req := MakePostRequest(t, "/auth/login", input)

		rr := execRequest(req, aServer)

		AssertStatus(t, rr.Code, http.StatusOK)

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
		req := MakePostRequest(t, "/auth/login", input)

		rr := execRequest(req, aServer)

		AssertStatus(t, rr.Code, http.StatusBadRequest)
	})
}
