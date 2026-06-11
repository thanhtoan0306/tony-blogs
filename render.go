package main

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

var tmpl *template.Template

func initTemplates() error {
	funcs := template.FuncMap{
		"formatPublishedAt": formatPublishedAt,
		"formatFullDate":    formatFullDate,
	}

	var err error
	tmpl, err = template.New("").Funcs(funcs).ParseFS(templateFS, "templates/*.html")
	return err
}

func render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
