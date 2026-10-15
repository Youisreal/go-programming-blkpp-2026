package httpjson

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lib/pq"
)

func TestID(t *testing.T) {
	for _, id := range []string{"1", "32767"} {
		request := httptest.NewRequest("GET", "/", nil)
		request.SetPathValue("id", id)
		value, ok := ID(httptest.NewRecorder(), request)
		if !ok || value < 1 || value > 32767 {
			t.Fatalf("valid ID rejected: %s", id)
		}
	}
}

func TestDatabaseErrors(t *testing.T) {
	for _, tc := range []struct {
		err    error
		status int
		code   string
	}{
		{sql.ErrNoRows, 404, "PRODUCT_NOT_FOUND"},
		{&pq.Error{Code: "23505"}, 409, "DUPLICATE_KEY"},
		{&pq.Error{Code: "23503"}, 409, "FOREIGN_KEY_CONFLICT"},
		{&pq.Error{Code: "23514"}, 400, "INVALID_INPUT"},
		{errors.New("internal SQL secret"), 500, "INTERNAL_ERROR"},
	} {
		response := httptest.NewRecorder()
		DatabaseError(response, tc.err, "PRODUCT")
		var body struct {
			Error struct {
				Code string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if response.Code != tc.status || body.Error.Code != tc.code {
			t.Fatalf("unexpected response: %d %s", response.Code, response.Body)
		}
		if strings.Contains(response.Body.String(), "SQL secret") {
			t.Fatal("internal error leaked")
		}
	}
}
