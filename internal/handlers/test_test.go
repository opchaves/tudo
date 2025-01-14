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

func makePostRequest(t *testing.T, path string, data interface{}) *http.Request {
	payload, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	return req
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}
