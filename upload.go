package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ArticleWriter interface {
	Create(ctx context.Context, a Article) error
}

type ArticleManager interface {
	SetVisible(ctx context.Context, slug string, visible bool) error
}

type uploadForm struct {
	Title       string
	Slug        string
	Summary     string
	BodyHTML    string
	Thumbnail   string
	Category    string
	Author      string
	PublishedAt string
	Tags        string
}

type uploadData struct {
	Title      string
	Error      string
	LoginError string
	Authed     bool
	CanUpload  bool
	Form       uploadForm
}

func handleUploadGet(store ArticleStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authed := uploadAuthed(r)
		_, canUpload := store.(ArticleWriter)
		w.Header().Set("Cache-Control", "no-store")
		render(w, "upload.html", uploadData{
			Title:     "Upload Article · Tony Blogs",
			Authed:    authed,
			CanUpload: canUpload,
			Form: uploadForm{
				Category: "General",
				Author:   "Crypto Desk",
			},
		})
	}
}

func handleUploadLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if !checkUploadPassword(strings.TrimSpace(r.FormValue("password"))) {
			next := uploadLoginNext(r.FormValue("next"))
			w.Header().Set("Cache-Control", "no-store")
			if next == "/upload/manage" {
				render(w, "manage.html", manageData{
					Title:      "Quản lý blog · Tony Blogs",
					Authed:     false,
					LoginError: "Mật khẩu không đúng",
					Next:       next,
				})
				return
			}
			render(w, "upload.html", uploadData{
				Title:      "Upload Article · Tony Blogs",
				Authed:     false,
				CanUpload:  true,
				LoginError: "Mật khẩu không đúng",
			})
			return
		}
		setUploadAuth(w)
		http.Redirect(w, r, uploadLoginNext(r.FormValue("next")), http.StatusSeeOther)
	}
}

func handleUploadPost(store ArticleStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !uploadAuthed(r) {
			http.Redirect(w, r, "/upload", http.StatusSeeOther)
			return
		}

		writer, canUpload := store.(ArticleWriter)
		if !canUpload {
			w.Header().Set("Cache-Control", "no-store")
			render(w, "upload.html", uploadData{
				Title:     "Upload Article · Tony Blogs",
				Authed:    true,
				CanUpload: false,
				Error:     "Upload requires Firestore (default). JSON mode is read-only.",
			})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
			return
		}

		form := uploadForm{
			Title:       strings.TrimSpace(r.FormValue("title")),
			Slug:        strings.TrimSpace(r.FormValue("slug")),
			Summary:     strings.TrimSpace(r.FormValue("summary")),
			BodyHTML:    strings.TrimSpace(r.FormValue("bodyHtml")),
			Thumbnail:   strings.TrimSpace(r.FormValue("thumbnail")),
			Category:    strings.TrimSpace(r.FormValue("category")),
			Author:      strings.TrimSpace(r.FormValue("author")),
			PublishedAt: strings.TrimSpace(r.FormValue("publishedAt")),
			Tags:        strings.TrimSpace(r.FormValue("tags")),
		}

		if err := validateUploadForm(form, store); err != nil {
			w.Header().Set("Cache-Control", "no-store")
			render(w, "upload.html", uploadData{
				Title:     "Upload Article · Tony Blogs",
				Authed:    true,
				CanUpload: true,
				Error:     err.Error(),
				Form:      form,
			})
			return
		}

		article := formToArticle(form)
		if err := writer.Create(r.Context(), article); err != nil {
			w.Header().Set("Cache-Control", "no-store")
			render(w, "upload.html", uploadData{
				Title:     "Upload Article · Tony Blogs",
				Authed:    true,
				CanUpload: true,
				Error:     err.Error(),
				Form:      form,
			})
			return
		}

		http.Redirect(w, r, "/news/"+article.Slug, http.StatusSeeOther)
	}
}

func validateUploadForm(form uploadForm, store ArticleStore) error {
	if form.Title == "" {
		return fmt.Errorf("title is required")
	}
	if form.Summary == "" {
		return fmt.Errorf("summary is required")
	}
	if form.BodyHTML == "" {
		return fmt.Errorf("body HTML is required")
	}

	slug := form.Slug
	if slug == "" {
		slug = slugify(form.Title)
	} else {
		slug = slugify(slug)
	}
	if slug == "" {
		return fmt.Errorf("could not generate a valid slug")
	}
	if _, exists := store.BySlug(slug); exists {
		return fmt.Errorf("slug %q already exists", slug)
	}
	return nil
}

func formToArticle(form uploadForm) Article {
	slug := form.Slug
	if slug == "" {
		slug = slugify(form.Title)
	} else {
		slug = slugify(slug)
	}

	category := form.Category
	if category == "" {
		category = "General"
	}
	author := form.Author
	if author == "" {
		author = "Crypto Desk"
	}
	publishedAt := form.PublishedAt
	if publishedAt == "" {
		publishedAt = time.Now().UTC().Format(time.RFC3339)
	}

	return Article{
		ID:          slug,
		Slug:        slug,
		Title:       form.Title,
		Summary:     form.Summary,
		BodyHTML:    form.BodyHTML,
		Thumbnail:   form.Thumbnail,
		Category:    category,
		Author:      author,
		PublishedAt: publishedAt,
		Tags:        parseTags(form.Tags),
		Visible:     boolPtr(true),
	}
}

func uploadLoginNext(next string) string {
	if next == "/upload/manage" {
		return next
	}
	return "/upload"
}

func parseTags(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			prevDash = false
			continue
		}
		if (r == ' ' || r == '-' || r == '_') && b.Len() > 0 && !prevDash {
			b.WriteByte('-')
			prevDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
