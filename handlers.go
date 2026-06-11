package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type indexData struct {
	Title    string
	Headline string
	Lede     string
	Articles []Article
	Mode     string
}

type articleData struct {
	Title   string
	Article Article
}

func handleIndex(articles []Article) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		today := todayIsoDate()
		feed := articlesForDate(articles, today)
		mode := "today"
		headline := "Crypto News Today"
		lede := formatFullDate(today) + " · " + itoa(len(feed)) + " stories"

		if len(feed) == 0 {
			feed = articles
			mode = "latest"
			headline = "Latest Crypto News"
			lede = "No stories for " + formatFullDate(today) + " · showing " + itoa(len(feed)) + " recent"
		}

		w.Header().Set("Cache-Control", "public, max-age=60")
		render(w, "index.html", indexData{
			Title:    "Crypto Today — Top Stories",
			Headline: headline,
			Lede:     lede,
			Articles: feed,
			Mode:     mode,
		})
	}
}

func handleArticle(articles []Article) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/news/")
		if slug == "" || strings.Contains(slug, "/") {
			http.NotFound(w, r)
			return
		}

		article, ok := articleBySlug(articles, slug)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			render(w, "404.html", map[string]string{"Title": "Not found"})
			return
		}

		w.Header().Set("Cache-Control", "public, max-age=300")
		render(w, "article.html", articleData{
			Title:   article.Title + " · Crypto Today",
			Article: article,
		})
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":        true,
		"service":   "golang-news",
		"mode":      "ssr",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
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
