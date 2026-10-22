package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientCORS(t *testing.T) {
	router := setupEngine(nil, []byte("test-secret"))
	for _, origin := range []string{"http://localhost:8081", "http://127.0.0.1:8081", "https://example.com"} {
		for _, method := range []string{http.MethodOptions, http.MethodGet} {
			request := httptest.NewRequest(method, "/health", nil)
			request.Header.Set("Origin", origin)
			request.Header.Set("Access-Control-Request-Method", "POST")
			request.Header.Set("Access-Control-Request-Headers", "authorization,content-type")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			want := origin
			if origin == "https://example.com" {
				want = ""
			}
			if got := response.Header().Get("Access-Control-Allow-Origin"); got != want {
				t.Fatalf("%s %s: origin = %q, want %q", method, origin, got, want)
			}
			if method == http.MethodOptions && response.Code != http.StatusNoContent {
				t.Fatalf("preflight status = %d", response.Code)
			}
			if want != "" && response.Header().Get("Access-Control-Allow-Headers") != "Authorization, Content-Type" {
				t.Fatal("header autentikasi dan JSON harus diizinkan")
			}
		}
	}
}
