package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	if err := initTemplates(); err != nil {
		log.Fatal(err)
	}

	articlesPath := os.Getenv("ARTICLES_JSON")
	if articlesPath == "" {
		articlesPath = "mockdb/articles.json"
	}

	articles, err := loadArticles(articlesPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded %d articles from %s", len(articles), articlesPath)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8093"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /static/", handleStatic)
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /news/{slug}", handleArticle(articles))
	mux.HandleFunc("GET /", handleIndex(articles))

	addr := "127.0.0.1:" + port
	log.Printf("golang-news SSR: http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
