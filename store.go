package main

import (
	"context"
	"fmt"
	"os"
)

// ArticleStore loads and serves articles from Firestore by default (JSON via DATA_SOURCE=json).
type ArticleStore interface {
	All() []Article
	BySlug(slug string) (Article, bool)
}

func newArticleStore(ctx context.Context) (ArticleStore, string, error) {
	switch os.Getenv("DATA_SOURCE") {
	case "json":
		path := os.Getenv("ARTICLES_JSON")
		if path == "" {
			path = "mockdb/articles.json"
		}
		articles, err := loadArticlesFromJSON(path)
		if err != nil {
			return nil, "", err
		}
		return newMemoryStore(articles), path, nil
	default:
		creds := os.Getenv("FIREBASE_CREDENTIALS")
		if creds == "" {
			creds = "golang-blogs-firebase-adminsdk-fbsvc-64dacce61f.json"
		}
		projectID := os.Getenv("FIREBASE_PROJECT_ID")
		if projectID == "" {
			projectID = "golang-blogs"
		}
		collection := os.Getenv("FIRESTORE_COLLECTION")
		if collection == "" {
			collection = "articles"
		}
		store, err := newFirestoreStore(ctx, creds, projectID, collection)
		if err != nil {
			return nil, "", err
		}
		return store, fmt.Sprintf("firestore:%s/%s", projectID, collection), nil
	}
}

type memoryStore struct {
	articles []Article
	bySlug   map[string]Article
}

func newMemoryStore(articles []Article) *memoryStore {
	bySlug := make(map[string]Article, len(articles))
	for _, a := range articles {
		bySlug[a.Slug] = a
	}
	return &memoryStore{articles: articles, bySlug: bySlug}
}

func (s *memoryStore) All() []Article { return s.articles }

func (s *memoryStore) BySlug(slug string) (Article, bool) {
	a, ok := s.bySlug[slug]
	return a, ok
}
