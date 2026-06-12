package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type indexData struct {
	Title      string
	Headline   string
	Lede       string
	Articles   []Article
	Mode       string
	DailyViews int64
}

type articleData struct {
	Title      string
	Article    Article
	DailyViews int64
}

func handleIndex(store ArticleStore, views ViewCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articles := visibleArticles(store.All())

		w.Header().Set("Cache-Control", "public, max-age=60")
		render(w, "index.html", indexData{
			Title:      "Tony Blogs — Top Stories",
			Headline:   "Latest Blogs",
			Lede:       itoa(len(articles)) + " stories",
			Articles:   articles,
			Mode:       "latest",
			DailyViews: recordPageView(r.Context(), views, "home"),
		})
	}
}

func handleArticle(store ArticleStore, views ViewCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/news/")
		if slug == "" || strings.Contains(slug, "/") {
			http.NotFound(w, r)
			return
		}

		article, ok := store.BySlug(slug)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			render(w, "404.html", map[string]string{"Title": "Not found"})
			return
		}

		w.Header().Set("Cache-Control", "public, max-age=300")
		render(w, "article.html", articleData{
			Title:      article.Title + " · Tony Blogs",
			Article:    article,
			DailyViews: recordPageView(r.Context(), views, "news:"+slug),
		})
	}
}

func handleHealth(source string, articleCount int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":             true,
			"service":        "golang-news",
			"mode":           "ssr",
			"dataSource":     source,
			"articleCount":   articleCount,
			"timestamp":      time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/static/")
	if path == "" || strings.Contains(path, "..") {
		http.NotFound(w, r)
		return
	}
	data, err := staticFS.ReadFile("static/" + path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if strings.HasSuffix(path, ".css") {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	}
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_, _ = w.Write(data)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
