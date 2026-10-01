package main

import "testing"

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
