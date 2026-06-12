package main

import (
	"net/http"
	"strings"
)

type manageData struct {
	Title      string
	Authed     bool
	CanManage  bool
	Articles   []Article
	Message    string
	LoginError string
	Next       string
}

func handleManageGet(store ArticleStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authed := uploadAuthed(r)
		_, canManage := store.(ArticleManager)
		w.Header().Set("Cache-Control", "no-store")
		render(w, "manage.html", manageData{
			Title:     "Quản lý blog · Tony Blogs",
			Authed:    authed,
			CanManage: canManage,
			Articles:  store.All(),
			Next:      "/upload/manage",
		})
	}
}

func handleManageVisible(store ArticleStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !uploadAuthed(r) {
			http.Redirect(w, r, "/upload/manage", http.StatusSeeOther)
			return
		}

		manager, canManage := store.(ArticleManager)
		if !canManage {
			w.Header().Set("Cache-Control", "no-store")
			render(w, "manage.html", manageData{
				Title:     "Quản lý blog · Tony Blogs",
				Authed:    true,
				CanManage: false,
				Articles:  store.All(),
				Message:   "Quản lý chỉ khả dụng với Firestore.",
				Next:      "/upload/manage",
			})
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		slug := strings.TrimSpace(r.FormValue("slug"))
		visible := r.FormValue("visible") == "1"
		if slug == "" {
			http.Redirect(w, r, "/upload/manage", http.StatusSeeOther)
			return
		}

		if err := manager.SetVisible(r.Context(), slug, visible); err != nil {
			w.Header().Set("Cache-Control", "no-store")
			render(w, "manage.html", manageData{
				Title:     "Quản lý blog · Tony Blogs",
				Authed:    true,
				CanManage: true,
				Articles:  store.All(),
				Message:   err.Error(),
				Next:      "/upload/manage",
			})
			return
		}

		http.Redirect(w, r, "/upload/manage", http.StatusSeeOther)
	}
}
