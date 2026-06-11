package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strings"
	"time"
)

type Article struct {
	ID          string   `json:"id"`
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Summary     string   `json:"summary"`
	Body        string   `json:"body"`
	BodyHTML    string   `json:"bodyHtml"`
	Category    string   `json:"category"`
	Author      string   `json:"author"`
	PublishedAt string   `json:"publishedAt"`
	Tags        []string `json:"tags"`
}

type articleStore struct {
	Articles []Article `json:"articles"`
}

func loadArticles(path string) ([]Article, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read articles: %w", err)
	}

	var store articleStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("parse articles: %w", err)
	}

	articles := store.Articles
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].PublishedAt > articles[j].PublishedAt
	})
	return articles, nil
}

func (a Article) BodyContent() template.HTML {
	if a.BodyHTML != "" {
		return template.HTML(a.BodyHTML)
	}
	var b strings.Builder
	for _, p := range a.Paragraphs() {
		b.WriteString("<p>")
		b.WriteString(template.HTMLEscapeString(p))
		b.WriteString("</p>")
	}
	return template.HTML(b.String())
}

func (a Article) Paragraphs() []string {
	parts := strings.Split(a.Body, "\n\n")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func formatPublishedAt(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso
	}
	return t.Format("Jan 2, 2006 · 3:04 PM MST")
}

func formatFullDate(isoDate string) string {
	t, err := time.Parse("2006-01-02", isoDate)
	if err != nil {
		return isoDate
	}
	return t.Format("Monday, January 2, 2006")
}

func todayIsoDate() string {
	return time.Now().UTC().Format("2006-01-02")
}

func articlesForDate(articles []Article, date string) []Article {
	var out []Article
	for _, a := range articles {
		if strings.HasPrefix(a.PublishedAt, date) {
			out = append(out, a)
		}
	}
	return out
}

func articleBySlug(articles []Article, slug string) (Article, bool) {
	for _, a := range articles {
		if a.Slug == slug {
			return a, true
		}
	}
	return Article{}, false
}
