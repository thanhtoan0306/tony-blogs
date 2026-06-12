package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type articleCreateRequest struct {
	Title       string   `json:"title"`
	Slug        string   `json:"slug,omitempty"`
	Summary     string   `json:"summary"`
	BodyHTML    string   `json:"bodyHtml"`
	Thumbnail   string   `json:"thumbnail,omitempty"`
	Category    string   `json:"category,omitempty"`
	Author      string   `json:"author,omitempty"`
	PublishedAt string   `json:"publishedAt,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Visible     *bool    `json:"visible,omitempty"`
}

type apiDocsData struct {
	Title   string
	BaseURL string
}

func handleAPIDocs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=300")
		render(w, "api-docs.html", apiDocsData{
			Title:   "API · Tony Blogs",
			BaseURL: requestBaseURL(r),
		})
	}
}

func handleArticlesAPIPost(store ArticleStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if !uploadAuthed(r) {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "unauthorized — set header X-Upload-Password or login at /upload",
			})
			return
		}

		writer, canUpload := store.(ArticleWriter)
		if !canUpload {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "upload requires Firestore (default). JSON mode is read-only",
			})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		var req articleCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON body"})
			return
		}

		form := uploadForm{
			Title:       strings.TrimSpace(req.Title),
			Slug:        strings.TrimSpace(req.Slug),
			Summary:     strings.TrimSpace(req.Summary),
			BodyHTML:    strings.TrimSpace(req.BodyHTML),
			Thumbnail:   strings.TrimSpace(req.Thumbnail),
			Category:    strings.TrimSpace(req.Category),
			Author:      strings.TrimSpace(req.Author),
			PublishedAt: strings.TrimSpace(req.PublishedAt),
			Tags:        strings.Join(req.Tags, ","),
		}

		if err := validateUploadForm(form, store); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		article := formToArticle(form)
		if req.Visible != nil {
			article.Visible = req.Visible
		}

		if err := writer.Create(r.Context(), article); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":    true,
			"slug":  article.Slug,
			"url":   "/news/" + article.Slug,
			"title": article.Title,
		})
	}
}

func requestBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	host := r.Host
	if host == "" {
		host = "127.0.0.1:8093"
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}
