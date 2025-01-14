package handlers_test

import (
	"net/http"
	"testing"
)

func TestGetNotes(t *testing.T) {
	t.Run("should be empty if no notes", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/notes", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer " + authToken)

		rr := execRequest(req, aServer)

		AssertStatus(t, rr.Code, http.StatusOK)

		expected := "[]\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("should be unauthorized if no token", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/notes", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := execRequest(req, aServer)

		AssertStatus(t, rr.Code, http.StatusUnauthorized)
	})
}
