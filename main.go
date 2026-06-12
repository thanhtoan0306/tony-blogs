package main

import (
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := initTemplates(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	prices := newPriceBoard()
	prices.Start(ctx)

	store, source, err := newArticleStore(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if closer, ok := store.(*firestoreStore); ok {
		defer closer.Close()
	}
	log.Printf("loaded %d articles from %s", len(store.All()), source)

	views, viewSource, err := newViewCounter(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if closer, ok := views.(*firestoreViewCounter); ok {
		defer closer.Close()
	}
	log.Printf("page views: %s", viewSource)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8093"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /static/", handleStatic)
	mux.HandleFunc("GET /health", handleHealth(source, len(store.All())))
	mux.HandleFunc("GET /news/{slug}", handleArticle(store, views))
	mux.HandleFunc("GET /upload", handleUploadGet(store))
	mux.HandleFunc("POST /upload/login", handleUploadLogin())
	mux.HandleFunc("POST /upload", handleUploadPost(store))
	mux.HandleFunc("GET /upload/manage", handleManageGet(store))
	mux.HandleFunc("POST /upload/manage/visible", handleManageVisible(store))
	mux.HandleFunc("GET /prices", handlePricesPage(prices, views))
	mux.HandleFunc("GET /api/prices", handlePricesAPI(prices))
	mux.HandleFunc("GET /api", handleAPIDocs())
	mux.HandleFunc("POST /api/articles", handleArticlesAPIPost(store))
	mux.HandleFunc("GET /", handleIndex(store, views))

	addr := listenAddr(port)
	log.Printf("golang-news SSR: http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func listenAddr(port string) string {
	if os.Getenv("VERCEL") != "" {
		return "0.0.0.0:" + port
	}
	return "127.0.0.1:" + port
}
