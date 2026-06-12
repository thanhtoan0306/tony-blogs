package main

import (
	"context"
	"log"
	"net/http"
	"strings"
)

func loadArticleComments(ctx context.Context, comments CommentStore, slug string) []Comment {
	if comments == nil {
		return nil
	}
	list, err := comments.ListByArticle(ctx, slug)
	if err != nil {
		log.Printf("comments %s: %v", slug, err)
		return nil
	}
	return list
}

func handleCommentPost(store ArticleStore, comments CommentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/news/")
		slug = strings.TrimSuffix(slug, "/comments")
		if slug == "" || strings.Contains(slug, "/") {
			http.NotFound(w, r)
			return
		}

		article, ok := store.BySlug(slug)
		if !ok {
			http.NotFound(w, r)
			return
		}

		if comments == nil {
			http.Error(w, "comments unavailable", http.StatusServiceUnavailable)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 8<<10)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
			return
		}

		form := commentForm{
			Nickname: r.FormValue("nickname"),
			Body:     r.FormValue("body"),
		}
		validated, err := validateCommentForm(form)
		if err != nil {
			w.Header().Set("Cache-Control", "no-store")
			render(w, "article.html", articleData{
				Title:        article.Title + " · Tony Blogs",
				Article:      article,
				Comments:     loadArticleComments(r.Context(), comments, slug),
				CommentError: err.Error(),
				CommentForm:  validated,
			})
			return
		}

		if err := comments.Create(r.Context(), newComment(slug, validated)); err != nil {
			w.Header().Set("Cache-Control", "no-store")
			render(w, "article.html", articleData{
				Title:        article.Title + " · Tony Blogs",
				Article:      article,
				Comments:     loadArticleComments(r.Context(), comments, slug),
				CommentError: err.Error(),
				CommentForm:  validated,
			})
			return
		}

		http.Redirect(w, r, "/news/"+slug+"#comments", http.StatusSeeOther)
	}
}
