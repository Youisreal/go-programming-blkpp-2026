package main

import (
	"encoding/json"
	"errors"
	"html/template"
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

	log.Println("server listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", newMux(articles)))
}

func newMux(articles []Article) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		writeHTML(w, pageData{Title: "HTTP Wiki", Articles: articles, IsIndex: true})
	})

	mux.HandleFunc("GET /articles", func(w http.ResponseWriter, r *http.Request) {
		writeHTML(w, pageData{Title: "Artikel", Articles: articles, IsIndex: true})
	})

	mux.HandleFunc("GET /articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		article, ok := findArticle(articles, r.PathValue("slug"))
		if !ok {
			http.NotFound(w, r)
			return
		}

		writeHTML(w, pageData{Title: article.Title, Article: article})
	})

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	return mux
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

type pageData struct {
	Title    string
	Articles []Article
	Article  Article
	IsIndex  bool
}

func writeHTML(w http.ResponseWriter, data pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := pageTemplate.Execute(w, data); err != nil {
		log.Println("write response:", err)
	}
}

var pageTemplate = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="id">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Title}}</title>
</head>
<body>
  <main>
    <header>
      <h1>{{.Title}}</h1>
      {{if .IsIndex}}<p>Daftar artikel singkat dari data JSON.</p>{{else}}<p><a href="/articles">Kembali ke daftar artikel</a></p>{{end}}
    </header>

    {{if .IsIndex}}
      <section>
        {{range .Articles}}
          <article>
            <h2><a href="/articles/{{.Slug}}">{{.Title}}</a></h2>
            <p>{{.Summary}}</p>
            <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>
          </article>
        {{end}}
      </section>
    {{else}}
      <article>
        <p>{{.Article.Summary}}</p>
        <ul>{{range .Article.Tags}}<li>{{.}}</li>{{end}}</ul>
      </article>
    {{end}}
  </main>
</body>
</html>
`))
