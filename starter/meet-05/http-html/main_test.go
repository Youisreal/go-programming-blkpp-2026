package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoadAndFindArticles(t *testing.T) {
	articles, err := loadArticles("data/articles.json")
	if err != nil {
		t.Fatal(err)
	}

	article, ok := findArticle(articles, "json")
	if !ok {
		t.Fatal("json article not found")
	}
	if article.Title != "JSON" {
		t.Fatalf("title = %q, want JSON", article.Title)
	}
}

func TestArticlesRenderHTML(t *testing.T) {
	articles, err := loadArticles("data/articles.json")
	if err != nil {
		t.Fatal(err)
	}

	server := newMux(articles)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/articles/json", nil)

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if contentType := response.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Fatalf("content-type = %q, want text/html; charset=utf-8", contentType)
	}
	if body := response.Body.String(); !strings.Contains(body, "<h1>JSON</h1>") || !strings.Contains(body, "Format data teks") {
		t.Fatalf("response body missing article HTML:\n%s", body)
	}
}
