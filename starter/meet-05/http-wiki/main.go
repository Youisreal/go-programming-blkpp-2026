package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
)

type Article struct {
	Slug    string   `json:"slug"`
	Title   string   `json:"title"`
	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`
}

func main() {
	articles, err := loadArticles("data/articles.json")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "HTTP wiki",
			"list":    "/articles",
			"detail":  "/articles/{slug}",
		})
	})

	http.HandleFunc("GET /articles", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, articles)
	})

	http.HandleFunc("GET /articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		article, ok := findArticle(articles, r.PathValue("slug"))
		if !ok {
			http.NotFound(w, r)
			return
		}

		writeJSON(w, http.StatusOK, article)
	})

	http.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	log.Println("server listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func loadArticles(path string) ([]Article, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var articles []Article
	if err := json.NewDecoder(file).Decode(&articles); err != nil {
		return nil, err
	}
	if len(articles) == 0 {
		return nil, errors.New("articles data is empty")
	}

	return articles, nil
}

func findArticle(articles []Article, slug string) (Article, bool) {
	for _, article := range articles {
		if article.Slug == slug {
			return article, true
		}
	}
	return Article{}, false
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("write response:", err)
	}
}
