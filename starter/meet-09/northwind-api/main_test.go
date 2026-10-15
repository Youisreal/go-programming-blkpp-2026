package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"syscall"
	"testing"

	"github.com/lib/pq"
)

func TestDatabaseStartupError(t *testing.T) {
	for _, tc := range []struct {
		err  error
		want string
	}{
		{fmt.Errorf("ping: %w", &net.OpError{Op: "dial", Net: "tcp", Err: syscall.ENETUNREACH}), "network is unreachable"},
		{&pq.Error{Code: "28P01", Message: "sensitive-password"}, "invalid_password"},
		{errors.New("parse postgresql://user:sensitive-password@host:badport"), "format DATABASE_URL"},
	} {
		message := databaseStartupError(tc.err).Error()
		if !strings.Contains(message, tc.want) || strings.Contains(message, "sensitive-password") {
			t.Fatalf("unexpected startup message: %s", message)
		}
	}
}

func TestRoutesWithoutDatabase(t *testing.T) {
	mux := newMux(nil)
	cases := []struct {
		name, method, path, body string
		status                   int
		code                     string
	}{
		{"health", "GET", "/health", "", 200, ""},
		{"health method", "POST", "/health", "", 405, "METHOD_NOT_ALLOWED"},
		{"invalid ID", "GET", "/api/products/nope", "", 400, "INVALID_ID"},
		{"zero ID", "GET", "/api/categories/0", "", 400, "INVALID_ID"},
		{"negative ID", "GET", "/api/orders/-1", "", 400, "INVALID_ID"},
		{"overflow ID", "DELETE", "/api/products/32768", "", 400, "INVALID_ID"},
		{"malformed JSON", "POST", "/api/products", "{", 400, "INVALID_JSON"},
		{"unknown field", "POST", "/api/products", `{"unknown":1}`, 400, "INVALID_JSON"},
		{"trailing JSON", "POST", "/api/categories", `{} {}`, 400, "INVALID_JSON"},
		{"negative price", "POST", "/api/products", `{"product_id":30000,"product_name":"Workshop","unit_price":-1}`, 400, "INVALID_INPUT"},
		{"negative stock", "POST", "/api/products", `{"product_id":30000,"product_name":"Workshop","units_in_stock":-1}`, 400, "INVALID_INPUT"},
		{"missing category name", "POST", "/api/categories", `{"category_id":30000}`, 400, "INVALID_INPUT"},
		{"mismatched ID", "PUT", "/api/categories/30000", `{"category_id":30001,"category_name":"Workshop"}`, 400, "INVALID_INPUT"},
		{"empty order", "POST", "/api/orders", `{"order_id":30000,"customer_id":"ALFKI","items":[]}`, 400, "INVALID_INPUT"},
		{"null body", "POST", "/api/products", `null`, 400, "INVALID_INPUT"},
		{"unsupported collection method", "PATCH", "/api/products", `{}`, 405, "METHOD_NOT_ALLOWED"},
		{"unsupported item method", "DELETE", "/api/orders/1", ``, 405, "METHOD_NOT_ALLOWED"},
		{"missing endpoint", "GET", "/missing", ``, 404, "NOT_FOUND"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body)))
			if response.Code != tc.status {
				t.Fatalf("status=%d want %d: %s", response.Code, tc.status, response.Body)
			}
			if response.Header().Get("Content-Type") != "application/json" {
				t.Fatal("expected JSON content type")
			}
			var body struct {
				Data  map[string]string `json:"data"`
				Error *struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if tc.code != "" && (body.Error == nil || body.Error.Code != tc.code || body.Error.Message == "") {
				t.Fatalf("invalid error envelope: %s", response.Body)
			}
			if tc.code == "" && body.Data["status"] != "ok" {
				t.Fatalf("invalid health: %s", response.Body)
			}
			if tc.status == http.StatusMethodNotAllowed && response.Header().Get("Allow") == "" {
				t.Fatal("missing Allow header")
			}
		})
	}
}
